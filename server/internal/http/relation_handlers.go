package http

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"server/internal/model"
)

type followReq struct {
	TargetType string `json:"targetType"`
	TargetID   uint64 `json:"targetId"`
}

func createFollowHandler(deps Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		ai, ok := RequireAuth(c, deps)
		if !ok {
			return
		}

		var req followReq
		if err := c.ShouldBindJSON(&req); err != nil {
			RespondFail(c, http.StatusBadRequest, 70001, "参数非法")
			return
		}

		req.TargetType = strings.ToLower(strings.TrimSpace(req.TargetType))
		if (req.TargetType != "forum" && req.TargetType != "user") || req.TargetID == 0 {
			RespondFail(c, http.StatusBadRequest, 70001, "参数非法")
			return
		}
		if req.TargetType == "user" && req.TargetID == ai.UserID {
			RespondFail(c, http.StatusBadRequest, 70001, "参数非法")
			return
		}

		if err := assertFollowTargetExists(deps, c, req.TargetType, req.TargetID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				RespondFail(c, http.StatusNotFound, 70002, "目标不存在")
				return
			}
			RespondFail(c, http.StatusInternalServerError, 50000, "internal error")
			return
		}

		now := time.Now()
		var f model.Follow
		err := deps.DB.WithContext(c.Request.Context()).
			Where("user_id = ? AND target_type = ? AND target_id = ?", ai.UserID, req.TargetType, req.TargetID).
			First(&f).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			RespondFail(c, http.StatusInternalServerError, 50000, "internal error")
			return
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			f = model.Follow{
				UserID:     ai.UserID,
				TargetType: req.TargetType,
				TargetID:   req.TargetID,
				Status:     "active",
				CanceledAt: nil,
				CreatedAt:  now,
				UpdatedAt:  now,
			}
			if err := deps.DB.WithContext(c.Request.Context()).Create(&f).Error; err != nil {
				if isMySQLDup(err) || isMySQLDup2(err) {
					RespondOK(c, map[string]any{})
					return
				}
				RespondFail(c, http.StatusInternalServerError, 50000, "internal error")
				return
			}
			RespondOK(c, map[string]any{})
			return
		}

		if f.Status == "active" {
			RespondOK(c, map[string]any{})
			return
		}

		if err := deps.DB.WithContext(c.Request.Context()).
			Model(&model.Follow{}).
			Where("id = ?", f.ID).
			Updates(map[string]any{
				"status":      "active",
				"canceled_at": nil,
				"updated_at":  now,
			}).Error; err != nil {
			RespondFail(c, http.StatusInternalServerError, 50000, "internal error")
			return
		}

		RespondOK(c, map[string]any{})
	}
}

func cancelFollowHandler(deps Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		ai, ok := RequireAuth(c, deps)
		if !ok {
			return
		}

		var req followReq
		if err := c.ShouldBindJSON(&req); err != nil {
			RespondFail(c, http.StatusBadRequest, 70001, "参数非法")
			return
		}

		req.TargetType = strings.ToLower(strings.TrimSpace(req.TargetType))
		if (req.TargetType != "forum" && req.TargetType != "user") || req.TargetID == 0 {
			RespondFail(c, http.StatusBadRequest, 70001, "参数非法")
			return
		}
		if req.TargetType == "user" && req.TargetID == ai.UserID {
			RespondFail(c, http.StatusBadRequest, 70001, "参数非法")
			return
		}

		now := time.Now()
		res := deps.DB.WithContext(c.Request.Context()).
			Model(&model.Follow{}).
			Where("user_id = ? AND target_type = ? AND target_id = ? AND status = ?", ai.UserID, req.TargetType, req.TargetID, "active").
			Updates(map[string]any{
				"status":      "canceled",
				"canceled_at": &now,
				"updated_at":  now,
			})
		if res.Error != nil {
			RespondFail(c, http.StatusInternalServerError, 50000, "internal error")
			return
		}

		RespondOK(c, map[string]any{})
	}
}

func listMyFollowsHandler(deps Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		ai, ok := RequireAuth(c, deps)
		if !ok {
			return
		}

		pageNum := parsePositiveIntWithDefault(c.Query("pageNum"), 1)
		pageSize := parsePositiveIntWithDefault(c.Query("pageSize"), 10)
		if pageNum <= 0 || pageSize <= 0 || pageSize > 200 {
			RespondFail(c, http.StatusBadRequest, 70003, "参数非法")
			return
		}

		targetType := strings.ToLower(strings.TrimSpace(c.Query("targetType")))
		if targetType != "" && targetType != "forum" && targetType != "user" {
			RespondFail(c, http.StatusBadRequest, 70003, "参数非法")
			return
		}

		q := deps.DB.WithContext(c.Request.Context()).
			Model(&model.Follow{}).
			Where("user_id = ? AND status = ?", ai.UserID, "active")
		if targetType != "" {
			q = q.Where("target_type = ?", targetType)
		}

		var total int64
		if err := q.Count(&total).Error; err != nil {
			RespondFail(c, http.StatusInternalServerError, 50000, "internal error")
			return
		}

		var rows []model.Follow
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
				"followId":   f.ID,
				"targetType": f.TargetType,
				"targetId":   f.TargetID,
				"createdAt":  f.CreatedAt.Format(time.RFC3339),
			})
		}

		RespondOK(c, map[string]any{"list": list, "total": total})
	}
}

func assertFollowTargetExists(deps Deps, c *gin.Context, targetType string, targetID uint64) error {
	switch targetType {
	case "forum":
		return deps.DB.WithContext(c.Request.Context()).Select("id").First(&model.Forum{}, targetID).Error
	case "user":
		return deps.DB.WithContext(c.Request.Context()).Select("id").First(&model.User{}, targetID).Error
	default:
		return gorm.ErrRecordNotFound
	}
}
