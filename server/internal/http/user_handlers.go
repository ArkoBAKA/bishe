package http

import (
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"server/internal/model"
)

type registerUserResp struct {
	UserID    uint64 `json:"userId"`
	AvatarURL string `json:"avatarUrl"`
}

func userRegisterHandler(deps Deps) gin.HandlerFunc {
	maxReqBytes := int64(deps.Config.Upload.MaxRequestMB) * 1024 * 1024
	if maxReqBytes <= 0 {
		maxReqBytes = 100 * 1024 * 1024
	}
	maxAvatarBytes := int64(deps.Config.Upload.MaxFileMB) * 1024 * 1024
	if maxAvatarBytes <= 0 {
		maxAvatarBytes = 50 * 1024 * 1024
	}

	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxReqBytes)
		if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
			RespondFail(c, http.StatusBadRequest, 10001, "参数非法")
			return
		}

		account := strings.TrimSpace(c.PostForm("account"))
		password := c.PostForm("password")
		nickname := strings.TrimSpace(c.PostForm("nickname"))

		if !validAccount(account) {
			RespondFail(c, http.StatusBadRequest, 10001, "参数非法")
			return
		}
		if !validPassword(password) {
			RespondFail(c, http.StatusBadRequest, 10003, "密码强度不足")
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			RespondFail(c, http.StatusInternalServerError, 50000, "internal error")
			return
		}

		u := model.User{
			Username:     account,
			Account:      account,
			PasswordHash: string(hash),
			Nickname:     nickname,
			AvatarURL:    "",
			Bio:          "",
			Role:         "user",
			Status:       "normal",
		}

		tx := deps.DB.WithContext(c.Request.Context()).Begin()
		if err := tx.Create(&u).Error; err != nil {
			_ = tx.Rollback().Error
			if isDupErr(err) {
				RespondFail(c, http.StatusBadRequest, 10002, "账号已存在")
				return
			}
			RespondFail(c, http.StatusInternalServerError, 50000, "internal error")
			return
		}

		if u.Nickname == "" {
			u.Nickname = "用户" + strconv.FormatUint(u.ID, 10)
			if err := tx.Model(&model.User{}).Where("id = ?", u.ID).Update("nickname", u.Nickname).Error; err != nil {
				_ = tx.Rollback().Error
				RespondFail(c, http.StatusInternalServerError, 50000, "internal error")
				return
			}
		}

		avatarHeader, err := c.FormFile("avatarFile")
		if err == nil && avatarHeader != nil {
			if avatarHeader.Size <= 0 || avatarHeader.Size > maxAvatarBytes {
				_ = tx.Rollback().Error
				RespondFail(c, http.StatusBadRequest, 10004, "头像上传失败")
				return
			}
			ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(avatarHeader.Filename), "."))
			if !isAllowedAvatarExt(ext) {
				_ = tx.Rollback().Error
				RespondFail(c, http.StatusBadRequest, 10004, "头像上传失败")
				return
			}

			relPath := filepath.Join("public", "avatar", strconv.FormatUint(u.ID, 10), newUUIDv4()+"."+ext)
			fullPath := filepath.Join(uploadBaseDir(deps), relPath)
			if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
				_ = tx.Rollback().Error
				RespondFail(c, http.StatusInternalServerError, 10004, "头像上传失败")
				return
			}
			if err := saveMultipartToFile(avatarHeader, fullPath, maxAvatarBytes); err != nil {
				_ = tx.Rollback().Error
				RespondFail(c, http.StatusInternalServerError, 10004, "头像上传失败")
				return
			}

			avatarURL := "/api/v1/static/public/" + filepath.ToSlash(strings.TrimPrefix(relPath, "public"+string(filepath.Separator)))
			if err := tx.Model(&model.User{}).Where("id = ?", u.ID).Update("avatar_url", avatarURL).Error; err != nil {
				_ = os.Remove(fullPath)
				_ = tx.Rollback().Error
				RespondFail(c, http.StatusInternalServerError, 10004, "头像上传失败")
				return
			}
			u.AvatarURL = avatarURL
		}

		if err := tx.Commit().Error; err != nil {
			RespondFail(c, http.StatusInternalServerError, 50000, "internal error")
			return
		}

		RespondOK(c, registerUserResp{UserID: u.ID, AvatarURL: u.AvatarURL})
	}
}

type loginReq struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func userLoginHandler(deps Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req loginReq
		if err := c.ShouldBindJSON(&req); err != nil {
			RespondFail(c, http.StatusBadRequest, 10011, "账号或密码错误")
			return
		}
		account := strings.TrimSpace(req.Account)
		if account == "" || req.Password == "" {
			RespondFail(c, http.StatusBadRequest, 10011, "账号或密码错误")
			return
		}

		var u model.User
		if err := deps.DB.WithContext(c.Request.Context()).Where("account = ?", account).First(&u).Error; err != nil {
			RespondFail(c, http.StatusBadRequest, 10011, "账号或密码错误")
			return
		}

		if u.Status == "banned" {
			if u.BanUntil == nil || u.BanUntil.After(time.Now()) {
				RespondFail(c, http.StatusBadRequest, 10012, "账号已封禁")
				return
			}
		}

		if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password)); err != nil {
			RespondFail(c, http.StatusBadRequest, 10011, "账号或密码错误")
			return
		}

		token, expiresIn := IssueJWT(deps.Config.Auth.JWTSecret, u.ID, u.Role, deps.Config.Auth.TokenTTLSeconds)
		RespondOK(c, map[string]any{
			"token":     token,
			"tokenType": "Bearer",
			"expiresIn": expiresIn,
			"user": map[string]any{
				"userId":    u.ID,
				"account":   u.Account,
				"nickname":  u.Nickname,
				"avatarUrl": u.AvatarURL,
				"bio":       u.Bio,
				"role":      u.Role,
				"status":    u.Status,
			},
		})
	}
}

func userLogoutHandler(deps Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, ok := RequireAuth(c, deps); !ok {
			return
		}
		RespondOK(c, map[string]any{})
	}
}

func userMeHandler(deps Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		ai, ok := RequireAuth(c, deps)
		if !ok {
			return
		}

		var u model.User
		if err := deps.DB.WithContext(c.Request.Context()).First(&u, ai.UserID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				RespondFail(c, http.StatusBadRequest, 10021, "用户不存在")
				return
			}
			RespondFail(c, http.StatusInternalServerError, 50000, "internal error")
			return
		}

		RespondOK(c, map[string]any{
			"userId":    u.ID,
			"account":   u.Account,
			"nickname":  u.Nickname,
			"avatarUrl": u.AvatarURL,
			"bio":       u.Bio,
			"role":      u.Role,
			"status":    u.Status,
			"createdAt": u.CreatedAt.Format(time.RFC3339),
		})
	}
}

func userPublicHandler(deps Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("userId"), 10, 64)
		if err != nil {
			RespondFail(c, http.StatusBadRequest, 10051, "用户不存在")
			return
		}

		var u model.User
		if err := deps.DB.WithContext(c.Request.Context()).First(&u, id).Error; err != nil {
			RespondFail(c, http.StatusBadRequest, 10051, "用户不存在")
			return
		}

		RespondOK(c, map[string]any{
			"userId":    u.ID,
			"nickname":  u.Nickname,
			"avatarUrl": u.AvatarURL,
			"bio":       u.Bio,
			"role":      u.Role,
		})
	}
}

func userUpdateProfileHandler(deps Deps) gin.HandlerFunc {
	maxReqBytes := int64(deps.Config.Upload.MaxRequestMB) * 1024 * 1024
	if maxReqBytes <= 0 {
		maxReqBytes = 100 * 1024 * 1024
	}
	maxAvatarBytes := int64(deps.Config.Upload.MaxFileMB) * 1024 * 1024
	if maxAvatarBytes <= 0 {
		maxAvatarBytes = 50 * 1024 * 1024
	}

	return func(c *gin.Context) {
		ai, ok := RequireAuth(c, deps)
		if !ok {
			return
		}

		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxReqBytes)
		if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
			RespondFail(c, http.StatusBadRequest, 10031, "参数非法")
			return
		}

		nickname := strings.TrimSpace(c.PostForm("nickname"))
		bio := strings.TrimSpace(c.PostForm("bio"))

		avatarHeader, avatarErr := c.FormFile("avatarFile")
		hasAvatar := avatarErr == nil && avatarHeader != nil

		updates := map[string]any{}
		if nickname != "" {
			if len([]rune(nickname)) > 64 {
				RespondFail(c, http.StatusBadRequest, 10031, "参数非法")
				return
			}
			updates["nickname"] = nickname
		}
		if bio != "" {
			if len([]rune(bio)) > 255 {
				RespondFail(c, http.StatusBadRequest, 10031, "参数非法")
				return
			}
			updates["bio"] = bio
		}

		tx := deps.DB.WithContext(c.Request.Context()).Begin()
		var avatarFullPath string
		if hasAvatar {
			if avatarHeader.Size <= 0 || avatarHeader.Size > maxAvatarBytes {
				_ = tx.Rollback().Error
				RespondFail(c, http.StatusBadRequest, 10032, "头像上传失败")
				return
			}
			ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(avatarHeader.Filename), "."))
			if !isAllowedAvatarExt(ext) {
				_ = tx.Rollback().Error
				RespondFail(c, http.StatusBadRequest, 10032, "头像上传失败")
				return
			}
			relPath := filepath.Join("public", "avatar", strconv.FormatUint(ai.UserID, 10), newUUIDv4()+"."+ext)
			avatarFullPath = filepath.Join(uploadBaseDir(deps), relPath)
			if err := os.MkdirAll(filepath.Dir(avatarFullPath), 0o755); err != nil {
				_ = tx.Rollback().Error
				RespondFail(c, http.StatusInternalServerError, 10032, "头像上传失败")
				return
			}
			if err := saveMultipartToFile(avatarHeader, avatarFullPath, maxAvatarBytes); err != nil {
				_ = tx.Rollback().Error
				RespondFail(c, http.StatusInternalServerError, 10032, "头像上传失败")
				return
			}
			avatarURL := "/api/v1/static/public/" + filepath.ToSlash(strings.TrimPrefix(relPath, "public"+string(filepath.Separator)))
			updates["avatar_url"] = avatarURL
		}

		if len(updates) == 0 {
			if avatarFullPath != "" {
				_ = os.Remove(avatarFullPath)
			}
			_ = tx.Rollback().Error
			RespondFail(c, http.StatusBadRequest, 10031, "参数非法")
			return
		}

		if err := tx.Model(&model.User{}).Where("id = ?", ai.UserID).Updates(updates).Error; err != nil {
			if avatarFullPath != "" {
				_ = os.Remove(avatarFullPath)
			}
			_ = tx.Rollback().Error
			RespondFail(c, http.StatusInternalServerError, 50000, "internal error")
			return
		}
		if err := tx.Commit().Error; err != nil {
			RespondFail(c, http.StatusInternalServerError, 50000, "internal error")
			return
		}
		RespondOK(c, map[string]any{})
	}
}

type updatePasswordReq struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required"`
}

func userUpdatePasswordHandler(deps Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		ai, ok := RequireAuth(c, deps)
		if !ok {
			return
		}

		var req updatePasswordReq
		if err := c.ShouldBindJSON(&req); err != nil {
			RespondFail(c, http.StatusBadRequest, 10042, "新密码不合法")
			return
		}
		if req.OldPassword == "" {
			RespondFail(c, http.StatusBadRequest, 10041, "旧密码错误")
			return
		}
		if req.NewPassword == "" || !validPassword(req.NewPassword) {
			RespondFail(c, http.StatusBadRequest, 10042, "新密码不合法")
			return
		}
		if req.OldPassword == req.NewPassword {
			RespondFail(c, http.StatusBadRequest, 10042, "新密码不合法")
			return
		}

		var u model.User
		if err := deps.DB.WithContext(c.Request.Context()).First(&u, ai.UserID).Error; err != nil {
			RespondFail(c, http.StatusBadRequest, 10041, "旧密码错误")
			return
		}
		if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.OldPassword)); err != nil {
			RespondFail(c, http.StatusBadRequest, 10041, "旧密码错误")
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			RespondFail(c, http.StatusInternalServerError, 50000, "internal error")
			return
		}
		if err := deps.DB.WithContext(c.Request.Context()).Model(&model.User{}).Where("id = ?", ai.UserID).Update("password_hash", string(hash)).Error; err != nil {
			RespondFail(c, http.StatusInternalServerError, 50000, "internal error")
			return
		}
		RespondOK(c, map[string]any{})
	}
}

func validAccount(s string) bool {
	if s == "" {
		return false
	}
	if len([]rune(s)) < 3 || len([]rune(s)) > 64 {
		return false
	}
	return true
}

func validPassword(s string) bool {
	if len(s) < 8 || len(s) > 72 {
		return false
	}
	hasLetter := false
	hasDigit := false
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
			hasLetter = true
			continue
		}
		if r >= '0' && r <= '9' {
			hasDigit = true
			continue
		}
	}
	if !hasLetter || !hasDigit {
		return false
	}
	return true
}

func isDupErr(err error) bool {
	var myErr *mysql.MySQLError
	if errors.As(err, &myErr) {
		return myErr.Number == 1062
	}
	return strings.Contains(strings.ToLower(err.Error()), "duplicate")
}

func isAllowedAvatarExt(ext string) bool {
	switch strings.ToLower(strings.TrimPrefix(ext, ".")) {
	case "jpg", "jpeg", "png", "gif", "webp":
		return true
	default:
		return false
	}
}

func uploadBaseDir(deps Deps) string {
	if deps.Config.Upload.BaseDir != "" {
		return deps.Config.Upload.BaseDir
	}
	return "data/uploads"
}

func saveMultipartToFile(fh *multipart.FileHeader, dest string, maxBytes int64) error {
	src, err := fh.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	tmpDir := filepath.Join(filepath.Dir(dest), ".tmp")
	if err := os.MkdirAll(tmpDir, 0o755); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(tmpDir, "up-*")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()

	limited := io.LimitReader(src, maxBytes+1)
	written, copyErr := io.Copy(tmp, limited)
	closeErr := tmp.Close()
	if copyErr != nil || closeErr != nil {
		_ = os.Remove(tmpName)
		return errors.New("write failed")
	}
	if written > maxBytes {
		_ = os.Remove(tmpName)
		return errors.New("too large")
	}
	if err := os.Rename(tmpName, dest); err != nil {
		_ = os.Remove(tmpName)
		return err
	}
	return nil
}
