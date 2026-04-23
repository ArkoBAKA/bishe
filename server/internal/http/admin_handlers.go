package http

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"server/internal/model"
)

func requireAdmin(c *gin.Context, deps Deps) (AuthInfo, bool) {
	ai, ok := RequireAuth(c, deps)
	if !ok {
		return AuthInfo{}, false
	}
	if ai.Role != "admin" {
		RespondFail(c, http.StatusForbidden, 90001, "需要管理员权限")
		return AuthInfo{}, false
	}
	return ai, true
}

func adminListPendingPostsHandler(deps Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, ok := requireAdmin(c, deps); !ok {
			return
		}

		pageNum := parsePositiveIntWithDefault(c.Query("pageNum"), 1)
		pageSize := parsePositiveIntWithDefault(c.Query("pageSize"), 10)
		if pageNum <= 0 || pageSize <= 0 || pageSize > 200 {
			RespondFail(c, http.StatusBadRequest, 20001, "参数非法")
			return
		}

		q := deps.DB.WithContext(c.Request.Context()).
			Model(&model.Post{}).
			Where("status = ?", "pending")

		var total int64
		if err := q.Count(&total).Error; err != nil {
			RespondFail(c, http.StatusInternalServerError, 50000, "internal error")
			return
		}

		var rows []model.Post
		if err := q.Order("created_at asc").
			Limit(pageSize).
			Offset((pageNum - 1) * pageSize).
			Find(&rows).Error; err != nil {
			RespondFail(c, http.StatusInternalServerError, 50000, "internal error")
			return
		}

		list := make([]map[string]any, 0, len(rows))
		for _, p := range rows {
			list = append(list, map[string]any{
				"postId":    p.ID,
				"forumId":   p.ForumID,
				"title":     p.Title,
				"authorId":  p.AuthorID,
				"createdAt": p.CreatedAt.Format(time.RFC3339),
			})
		}

		RespondOK(c, map[string]any{"list": list, "total": total})
	}
}

type reviewPostReq struct {
	Action       string `json:"action"`
	ReviewRemark string `json:"reviewRemark"`
}

func adminReviewPostHandler(deps Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		ai, ok := requireAdmin(c, deps)
		if !ok {
			return
		}

		postID, err := strconv.ParseUint(c.Param("postId"), 10, 64)
		if err != nil {
			RespondFail(c, http.StatusBadRequest, 90003, "帖子不存在")
			return
		}

		var req reviewPostReq
		if err := c.ShouldBindJSON(&req); err != nil {
			RespondFail(c, http.StatusBadRequest, 90002, "参数非法")
			return
		}
		req.Action = strings.ToLower(strings.TrimSpace(req.Action))
		req.ReviewRemark = strings.TrimSpace(req.ReviewRemark)
		if req.Action != "approve" && req.Action != "reject" && req.Action != "hide" {
			RespondFail(c, http.StatusBadRequest, 90002, "参数非法")
			return
		}
		if len([]rune(req.ReviewRemark)) > 255 {
			RespondFail(c, http.StatusBadRequest, 90002, "参数非法")
			return
		}

		var p model.Post
		if err := deps.DB.WithContext(c.Request.Context()).First(&p, postID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				RespondFail(c, http.StatusNotFound, 90003, "帖子不存在")
				return
			}
			RespondFail(c, http.StatusInternalServerError, 50000, "internal error")
			return
		}

		newStatus := ""
		switch req.Action {
		case "approve":
			newStatus = "visible"
		case "reject":
			newStatus = "rejected"
		case "hide":
			newStatus = "hidden"
		}
		now := time.Now()

		if err := deps.DB.WithContext(c.Request.Context()).
			Model(&model.Post{}).
			Where("id = ?", postID).
			Updates(map[string]any{
				"status":        newStatus,
				"reviewed_by":   ai.UserID,
				"reviewed_at":   &now,
				"review_remark": req.ReviewRemark,
				"updated_at":    now,
			}).Error; err != nil {
			RespondFail(c, http.StatusInternalServerError, 50000, "internal error")
			return
		}

		RespondOK(c, map[string]any{"status": newStatus})
	}
}

type createReportReq struct {
	TargetType string `json:"targetType"`
	TargetID   uint64 `json:"targetId"`
	Reason     string `json:"reason"`
	Detail     string `json:"detail"`
}

func createReportHandler(deps Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		ai, ok := RequireAuth(c, deps)
		if !ok {
			return
		}

		var req createReportReq
		if err := c.ShouldBindJSON(&req); err != nil {
			RespondFail(c, http.StatusBadRequest, 91001, "参数非法")
			return
		}
		req.TargetType = strings.ToLower(strings.TrimSpace(req.TargetType))
		req.Reason = strings.TrimSpace(req.Reason)
		req.Detail = strings.TrimSpace(req.Detail)

		if (req.TargetType != "post" && req.TargetType != "comment" && req.TargetType != "user") || req.TargetID == 0 {
			RespondFail(c, http.StatusBadRequest, 91001, "参数非法")
			return
		}
		if len([]rune(req.Reason)) == 0 || len([]rune(req.Reason)) > 64 {
			RespondFail(c, http.StatusBadRequest, 91001, "参数非法")
			return
		}
		if len([]rune(req.Detail)) > 1000 {
			RespondFail(c, http.StatusBadRequest, 91001, "参数非法")
			return
		}

		r := model.Report{
			ReporterID:    ai.UserID,
			TargetType:    req.TargetType,
			TargetID:      req.TargetID,
			Reason:        req.Reason,
			Detail:        req.Detail,
			Status:        "pending",
			ProcessedBy:   0,
			ProcessedAt:   nil,
			ProcessAction: "",
			ProcessRemark: "",
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}
		if err := deps.DB.WithContext(c.Request.Context()).Create(&r).Error; err != nil {
			RespondFail(c, http.StatusInternalServerError, 50000, "internal error")
			return
		}

		RespondOK(c, map[string]any{"reportId": r.ID})
	}
}

func adminListReportsHandler(deps Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, ok := requireAdmin(c, deps); !ok {
			return
		}

		pageNum := parsePositiveIntWithDefault(c.Query("pageNum"), 1)
		pageSize := parsePositiveIntWithDefault(c.Query("pageSize"), 10)
		if pageNum <= 0 || pageSize <= 0 || pageSize > 200 {
			RespondFail(c, http.StatusBadRequest, 20001, "参数非法")
			return
		}
		status := strings.ToLower(strings.TrimSpace(c.Query("status")))
		if status != "" && status != "pending" && status != "processed" {
			RespondFail(c, http.StatusBadRequest, 20001, "参数非法")
			return
		}

		q := deps.DB.WithContext(c.Request.Context()).Model(&model.Report{})
		if status != "" {
			q = q.Where("status = ?", status)
		}

		var total int64
		if err := q.Count(&total).Error; err != nil {
			RespondFail(c, http.StatusInternalServerError, 50000, "internal error")
			return
		}

		var rows []model.Report
		if err := q.Order("created_at desc").
			Limit(pageSize).
			Offset((pageNum - 1) * pageSize).
			Find(&rows).Error; err != nil {
			RespondFail(c, http.StatusInternalServerError, 50000, "internal error")
			return
		}

		list := make([]map[string]any, 0, len(rows))
		for _, r := range rows {
			item := map[string]any{
				"reportId":   r.ID,
				"reporterId": r.ReporterID,
				"targetType": r.TargetType,
				"targetId":   r.TargetID,
				"reason":     r.Reason,
				"detail":     r.Detail,
				"status":     r.Status,
				"createdAt":  r.CreatedAt.Format(time.RFC3339),
			}
			list = append(list, item)
		}

		RespondOK(c, map[string]any{"list": list, "total": total})
	}
}

type processReportReq struct {
	Action          string `json:"action"`
	ProcessRemark   string `json:"processRemark"`
	BanUntil        string `json:"banUntil"`
	DurationSeconds int64  `json:"durationSeconds"`
}

func adminProcessReportHandler(deps Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		ai, ok := requireAdmin(c, deps)
		if !ok {
			return
		}

		reportID, err := strconv.ParseUint(c.Param("reportId"), 10, 64)
		if err != nil {
			RespondFail(c, http.StatusNotFound, 92002, "举报单不存在")
			return
		}

		var req processReportReq
		if err := c.ShouldBindJSON(&req); err != nil {
			RespondFail(c, http.StatusBadRequest, 92001, "参数非法")
			return
		}
		req.Action = strings.TrimSpace(req.Action)
		req.ProcessRemark = strings.TrimSpace(req.ProcessRemark)
		if len([]rune(req.ProcessRemark)) > 255 {
			RespondFail(c, http.StatusBadRequest, 92001, "参数非法")
			return
		}

		var rpt model.Report
		if err := deps.DB.WithContext(c.Request.Context()).First(&rpt, reportID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				RespondFail(c, http.StatusNotFound, 92002, "举报单不存在")
				return
			}
			RespondFail(c, http.StatusInternalServerError, 50000, "internal error")
			return
		}

		if rpt.Status == "processed" {
			RespondFail(c, http.StatusConflict, 92003, "状态冲突")
			return
		}

		action := strings.ToLower(req.Action)
		switch action {
		case "close", "deletepost", "deletecomment", "hidepost", "banuser", "deletePost", "deleteComment", "hidePost", "banUser":
		default:
			RespondFail(c, http.StatusBadRequest, 92001, "参数非法")
			return
		}
		switch action {
		case "deletePost":
			action = "deletepost"
		case "deleteComment":
			action = "deletecomment"
		case "hidePost":
			action = "hidepost"
		case "banUser":
			action = "banuser"
		}

		if action == "deletepost" || action == "hidepost" {
			if rpt.TargetType != "post" {
				RespondFail(c, http.StatusBadRequest, 92004, "目标不匹配")
				return
			}
		}
		if action == "deletecomment" {
			if rpt.TargetType != "comment" {
				RespondFail(c, http.StatusBadRequest, 92004, "目标不匹配")
				return
			}
		}
		if action == "banuser" {
			if rpt.TargetType != "user" {
				RespondFail(c, http.StatusBadRequest, 92004, "目标不匹配")
				return
			}
		}

		tx := deps.DB.WithContext(c.Request.Context()).Begin()
		now := time.Now()

		switch action {
		case "close":
		case "deletepost":
			if err := tx.Delete(&model.Post{}, rpt.TargetID).Error; err != nil {
				_ = tx.Rollback().Error
				RespondFail(c, http.StatusInternalServerError, 50000, "internal error")
				return
			}
		case "hidepost":
			if err := tx.Model(&model.Post{}).Where("id = ?", rpt.TargetID).Update("status", "hidden").Error; err != nil {
				_ = tx.Rollback().Error
				RespondFail(c, http.StatusInternalServerError, 50000, "internal error")
				return
			}
		case "deletecomment":
			if err := tx.Delete(&model.Comment{}, rpt.TargetID).Error; err != nil {
				_ = tx.Rollback().Error
				RespondFail(c, http.StatusInternalServerError, 50000, "internal error")
				return
			}
		case "banuser":
			banUntil := (*time.Time)(nil)
			if strings.TrimSpace(req.BanUntil) != "" {
				t, err := time.Parse(time.RFC3339, strings.TrimSpace(req.BanUntil))
				if err != nil {
					_ = tx.Rollback().Error
					RespondFail(c, http.StatusBadRequest, 92001, "参数非法")
					return
				}
				banUntil = &t
			} else if req.DurationSeconds > 0 {
				t := now.Add(time.Duration(req.DurationSeconds) * time.Second)
				banUntil = &t
			} else {
				t := now.Add(7 * 24 * time.Hour)
				banUntil = &t
			}
			if err := tx.Model(&model.User{}).Where("id = ?", rpt.TargetID).Updates(map[string]any{
				"status":     "banned",
				"ban_until":  banUntil,
				"updated_at": now,
			}).Error; err != nil {
				_ = tx.Rollback().Error
				RespondFail(c, http.StatusInternalServerError, 50000, "internal error")
				return
			}
		}

		if err := tx.Model(&model.Report{}).Where("id = ?", rpt.ID).Updates(map[string]any{
			"status":         "processed",
			"processed_by":   ai.UserID,
			"processed_at":   &now,
			"process_action": action,
			"process_remark": req.ProcessRemark,
			"updated_at":     now,
		}).Error; err != nil {
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

func adminDeletePostHandler(deps Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, ok := requireAdmin(c, deps); !ok {
			return
		}
		postID, err := strconv.ParseUint(c.Param("postId"), 10, 64)
		if err != nil {
			RespondFail(c, http.StatusNotFound, 90003, "帖子不存在")
			return
		}
		res := deps.DB.WithContext(c.Request.Context()).Delete(&model.Post{}, postID)
		if res.Error != nil {
			RespondFail(c, http.StatusInternalServerError, 50000, "internal error")
			return
		}
		if res.RowsAffected == 0 {
			RespondFail(c, http.StatusNotFound, 90003, "帖子不存在")
			return
		}
		RespondOK(c, map[string]any{})
	}
}

func adminDeleteCommentHandler(deps Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, ok := requireAdmin(c, deps); !ok {
			return
		}
		commentID, err := strconv.ParseUint(c.Param("commentId"), 10, 64)
		if err != nil {
			RespondFail(c, http.StatusNotFound, 90013, "评论不存在")
			return
		}
		res := deps.DB.WithContext(c.Request.Context()).Delete(&model.Comment{}, commentID)
		if res.Error != nil {
			RespondFail(c, http.StatusInternalServerError, 50000, "internal error")
			return
		}
		if res.RowsAffected == 0 {
			RespondFail(c, http.StatusNotFound, 90013, "评论不存在")
			return
		}
		RespondOK(c, map[string]any{})
	}
}

type banUserReq struct {
	BanUntil        string `json:"banUntil"`
	DurationSeconds int64  `json:"durationSeconds"`
	Remark          string `json:"remark"`
}

func adminBanUserHandler(deps Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, ok := requireAdmin(c, deps); !ok {
			return
		}

		userID, err := strconv.ParseUint(c.Param("userId"), 10, 64)
		if err != nil {
			RespondFail(c, http.StatusNotFound, 93002, "用户不存在")
			return
		}

		var req banUserReq
		if err := c.ShouldBindJSON(&req); err != nil {
			RespondFail(c, http.StatusBadRequest, 93001, "参数非法")
			return
		}
		req.Remark = strings.TrimSpace(req.Remark)
		if len([]rune(req.Remark)) > 255 {
			RespondFail(c, http.StatusBadRequest, 93001, "参数非法")
			return
		}

		now := time.Now()
		banUntil := (*time.Time)(nil)
		if strings.TrimSpace(req.BanUntil) != "" {
			t, err := time.Parse(time.RFC3339, strings.TrimSpace(req.BanUntil))
			if err != nil {
				RespondFail(c, http.StatusBadRequest, 93001, "参数非法")
				return
			}
			banUntil = &t
		} else if req.DurationSeconds > 0 {
			t := now.Add(time.Duration(req.DurationSeconds) * time.Second)
			banUntil = &t
		} else {
			t := now.Add(7 * 24 * time.Hour)
			banUntil = &t
		}

		res := deps.DB.WithContext(c.Request.Context()).
			Model(&model.User{}).
			Where("id = ?", userID).
			Updates(map[string]any{
				"status":     "banned",
				"ban_until":  banUntil,
				"updated_at": now,
			})
		if res.Error != nil {
			RespondFail(c, http.StatusInternalServerError, 50000, "internal error")
			return
		}
		if res.RowsAffected == 0 {
			RespondFail(c, http.StatusNotFound, 93002, "用户不存在")
			return
		}

		banUntilStr := ""
		if banUntil != nil {
			banUntilStr = banUntil.Format(time.RFC3339)
		}
		RespondOK(c, map[string]any{
			"status":   "banned",
			"banUntil": banUntilStr,
		})
	}
}
