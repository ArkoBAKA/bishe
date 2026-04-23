package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"server/internal/model"
)

func listNotificationsHandler(deps Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		ai, ok := RequireAuth(c, deps)
		if !ok {
			return
		}

		pageNum := parsePositiveIntWithDefault(c.Query("pageNum"), 1)
		pageSize := parsePositiveIntWithDefault(c.Query("pageSize"), 10)
		if pageNum <= 0 || pageSize <= 0 || pageSize > 200 {
			RespondFail(c, http.StatusBadRequest, 80001, "参数非法")
			return
		}

		isReadParam := strings.TrimSpace(c.Query("isRead"))
		hasIsRead := false
		isRead := false
		if isReadParam != "" {
			switch strings.ToLower(isReadParam) {
			case "true", "1":
				hasIsRead = true
				isRead = true
			case "false", "0":
				hasIsRead = true
				isRead = false
			default:
				RespondFail(c, http.StatusBadRequest, 80001, "参数非法")
				return
			}
		}

		q := deps.DB.WithContext(c.Request.Context()).
			Model(&model.Notification{}).
			Where("user_id = ?", ai.UserID)
		if hasIsRead {
			q = q.Where("is_read = ?", isRead)
		}

		var total int64
		if err := q.Count(&total).Error; err != nil {
			RespondFail(c, http.StatusInternalServerError, 50000, "internal error")
			return
		}

		var rows []model.Notification
		if err := q.Order("created_at desc").
			Limit(pageSize).
			Offset((pageNum - 1) * pageSize).
			Find(&rows).Error; err != nil {
			RespondFail(c, http.StatusInternalServerError, 50000, "internal error")
			return
		}

		list := make([]map[string]any, 0, len(rows))
		for _, n := range rows {
			dataObj := map[string]any{}
			if strings.TrimSpace(n.DataJSON) != "" {
				_ = json.Unmarshal([]byte(n.DataJSON), &dataObj)
			}
			readAt := ""
			if n.ReadAt != nil {
				readAt = n.ReadAt.Format(time.RFC3339)
			}
			list = append(list, map[string]any{
				"notificationId": n.ID,
				"type":           n.Type,
				"title":          n.Title,
				"content":        n.Content,
				"data":           dataObj,
				"isRead":         n.IsRead,
				"createdAt":      n.CreatedAt.Format(time.RFC3339),
				"readAt":         readAt,
			})
		}

		RespondOK(c, map[string]any{"list": list, "total": total})
	}
}

func markNotificationReadHandler(deps Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		ai, ok := RequireAuth(c, deps)
		if !ok {
			return
		}

		id, err := strconv.ParseUint(c.Param("notificationId"), 10, 64)
		if err != nil || id == 0 {
			RespondFail(c, http.StatusNotFound, 80002, "通知不存在")
			return
		}

		var n model.Notification
		if err := deps.DB.WithContext(c.Request.Context()).
			Where("id = ? AND user_id = ?", id, ai.UserID).
			First(&n).Error; err != nil {
			RespondFail(c, http.StatusNotFound, 80002, "通知不存在")
			return
		}

		if n.IsRead {
			RespondOK(c, map[string]any{})
			return
		}

		now := time.Now()
		if err := deps.DB.WithContext(c.Request.Context()).
			Model(&model.Notification{}).
			Where("id = ? AND user_id = ? AND is_read = ?", id, ai.UserID, false).
			Updates(map[string]any{
				"is_read":    true,
				"read_at":    &now,
				"updated_at": now,
			}).Error; err != nil {
			RespondFail(c, http.StatusInternalServerError, 50000, "internal error")
			return
		}

		RespondOK(c, map[string]any{})
	}
}

func markAllNotificationsReadHandler(deps Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		ai, ok := RequireAuth(c, deps)
		if !ok {
			return
		}

		now := time.Now()
		res := deps.DB.WithContext(c.Request.Context()).
			Model(&model.Notification{}).
			Where("user_id = ? AND is_read = ?", ai.UserID, false).
			Updates(map[string]any{
				"is_read":    true,
				"read_at":    &now,
				"updated_at": now,
			})
		if res.Error != nil {
			RespondFail(c, http.StatusInternalServerError, 50000, "internal error")
			return
		}

		RespondOK(c, map[string]any{"affected": res.RowsAffected})
	}
}

func createNotification(deps Deps, ctx *gin.Context, userID uint64, typ, title, content string, data map[string]any) {
	if userID == 0 {
		return
	}
	typ = strings.TrimSpace(typ)
	title = strings.TrimSpace(title)
	content = strings.TrimSpace(content)
	if typ == "" {
		typ = "system"
	}
	if len([]rune(title)) > 64 {
		title = string([]rune(title)[:64])
	}
	if len([]rune(content)) > 255 {
		content = string([]rune(content)[:255])
	}
	dataJSON := ""
	if data != nil {
		if b, err := json.Marshal(data); err == nil {
			if len(b) <= 2000 {
				dataJSON = string(b)
			} else {
				dataJSON = string(b[:2000])
			}
		}
	}
	now := time.Now()
	_ = deps.DB.WithContext(ctx.Request.Context()).Create(&model.Notification{
		UserID:    userID,
		Type:      typ,
		Title:     title,
		Content:   content,
		DataJSON:  dataJSON,
		IsRead:    false,
		ReadAt:    nil,
		CreatedAt: now,
		UpdatedAt: now,
	}).Error
}
