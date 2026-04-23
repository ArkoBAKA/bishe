package http

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"

	"server/internal/model"
)

type createForumReq struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	CoverURL    string `json:"coverUrl"`
}

func createForumHandler(deps Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		ai, ok := RequireAuth(c, deps)
		if !ok {
			return
		}

		var req createForumReq
		if err := c.ShouldBindJSON(&req); err != nil {
			RespondFail(c, http.StatusBadRequest, 20001, "参数非法")
			return
		}

		name := strings.TrimSpace(req.Name)
		desc := strings.TrimSpace(req.Description)
		cover := strings.TrimSpace(req.CoverURL)

		if !validForumName(name) || len([]rune(desc)) > 255 || len(cover) > 512 {
			RespondFail(c, http.StatusBadRequest, 20001, "参数非法")
			return
		}

		f := model.Forum{
			Name:        name,
			Description: desc,
			CoverURL:    cover,
			OwnerID:     ai.UserID,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		if err := deps.DB.WithContext(c.Request.Context()).Create(&f).Error; err != nil {
			if isMySQLDup(err) {
				RespondFail(c, http.StatusBadRequest, 20002, "贴吧已存在")
				return
			}
			RespondFail(c, http.StatusInternalServerError, 50000, "internal error")
			return
		}

		RespondOK(c, map[string]any{"forumId": f.ID})
	}
}

type updateForumReq struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	CoverURL    *string `json:"coverUrl"`
}

func updateForumHandler(deps Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		ai, ok := RequireAuth(c, deps)
		if !ok {
			return
		}

		forumID, err := strconv.ParseUint(c.Param("forumId"), 10, 64)
		if err != nil {
			RespondFail(c, http.StatusBadRequest, 20001, "参数非法")
			return
		}

		var req updateForumReq
		if err := c.ShouldBindJSON(&req); err != nil {
			RespondFail(c, http.StatusBadRequest, 20001, "参数非法")
			return
		}

		var f model.Forum
		if err := deps.DB.WithContext(c.Request.Context()).First(&f, forumID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				RespondFail(c, http.StatusNotFound, 20003, "贴吧不存在")
				return
			}
			RespondFail(c, http.StatusInternalServerError, 50000, "internal error")
			return
		}

		if f.OwnerID != ai.UserID {
			RespondFail(c, http.StatusForbidden, 20004, "无权限")
			return
		}

		updates := map[string]any{}
		if req.Name != nil {
			name := strings.TrimSpace(*req.Name)
			if !validForumName(name) {
				RespondFail(c, http.StatusBadRequest, 20001, "参数非法")
				return
			}
			updates["name"] = name
		}
		if req.Description != nil {
			desc := strings.TrimSpace(*req.Description)
			if len([]rune(desc)) > 255 {
				RespondFail(c, http.StatusBadRequest, 20001, "参数非法")
				return
			}
			updates["description"] = desc
		}
		if req.CoverURL != nil {
			cover := strings.TrimSpace(*req.CoverURL)
			if len(cover) > 512 {
				RespondFail(c, http.StatusBadRequest, 20001, "参数非法")
				return
			}
			updates["cover_url"] = cover
		}

		if len(updates) == 0 {
			RespondFail(c, http.StatusBadRequest, 20001, "参数非法")
			return
		}
		updates["updated_at"] = time.Now()

		if err := deps.DB.WithContext(c.Request.Context()).Model(&model.Forum{}).Where("id = ?", forumID).Updates(updates).Error; err != nil {
			if isMySQLDup(err) {
				RespondFail(c, http.StatusBadRequest, 20002, "贴吧已存在")
				return
			}
			RespondFail(c, http.StatusInternalServerError, 50000, "internal error")
			return
		}

		RespondOK(c, map[string]any{})
	}
}

func listForumsHandler(deps Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		pageNum := parsePositiveIntWithDefault(c.Query("pageNum"), 1)
		pageSize := parsePositiveIntWithDefault(c.Query("pageSize"), 10)
		if pageNum <= 0 || pageSize <= 0 || pageSize > 100 {
			RespondFail(c, http.StatusBadRequest, 20001, "参数非法")
			return
		}

		keyword := strings.TrimSpace(c.Query("keyword"))

		q := deps.DB.WithContext(c.Request.Context()).Model(&model.Forum{})
		if keyword != "" {
			q = q.Where("name LIKE ?", "%"+keyword+"%")
		}

		var total int64
		if err := q.Count(&total).Error; err != nil {
			RespondFail(c, http.StatusInternalServerError, 50000, "internal error")
			return
		}

		var rows []model.Forum
		if err := q.Order("created_at desc").
			Limit(pageSize).
			Offset((pageNum - 1) * pageSize).
			Find(&rows).Error; err != nil {
			RespondFail(c, http.StatusInternalServerError, 50000, "internal error")
			return
		}

		list := make([]map[string]any, 0, len(rows))
		for _, f := range rows {
			list = append(list, map[string]any{
				"forumId":     f.ID,
				"name":        f.Name,
				"description": f.Description,
				"coverUrl":    f.CoverURL,
				"ownerId":     f.OwnerID,
				"createdAt":   f.CreatedAt.Format(time.RFC3339),
			})
		}

		RespondOK(c, map[string]any{
			"list":  list,
			"total": total,
		})
	}
}

func getForumHandler(deps Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		forumID, err := strconv.ParseUint(c.Param("forumId"), 10, 64)
		if err != nil {
			RespondFail(c, http.StatusBadRequest, 20003, "贴吧不存在")
			return
		}

		var f model.Forum
		if err := deps.DB.WithContext(c.Request.Context()).First(&f, forumID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				RespondFail(c, http.StatusNotFound, 20003, "贴吧不存在")
				return
			}
			RespondFail(c, http.StatusInternalServerError, 50000, "internal error")
			return
		}

		RespondOK(c, map[string]any{
			"forumId":     f.ID,
			"name":        f.Name,
			"description": f.Description,
			"coverUrl":    f.CoverURL,
			"ownerId":     f.OwnerID,
			"createdAt":   f.CreatedAt.Format(time.RFC3339),
		})
	}
}

func validForumName(name string) bool {
	if name == "" {
		return false
	}
	if len([]rune(name)) > 64 {
		return false
	}
	return true
}

func parsePositiveIntWithDefault(s string, def int) int {
	s = strings.TrimSpace(s)
	if s == "" {
		return def
	}
	n, err := strconv.Atoi(s)
	if err != nil || n <= 0 {
		return -1
	}
	return n
}

func isMySQLDup(err error) bool {
	var myErr *mysql.MySQLError
	if errors.As(err, &myErr) {
		return myErr.Number == 1062
	}
	return strings.Contains(strings.ToLower(err.Error()), "duplicate")
}
