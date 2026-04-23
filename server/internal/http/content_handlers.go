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

type createPostReq struct {
	ForumID uint64 `json:"forumId"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func createPostHandler(deps Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		ai, ok := RequireAuth(c, deps)
		if !ok {
			return
		}

		var req createPostReq
		if err := c.ShouldBindJSON(&req); err != nil {
			RespondFail(c, http.StatusBadRequest, 30001, "参数非法")
			return
		}

		title := strings.TrimSpace(req.Title)
		content := strings.TrimSpace(req.Content)
		if req.ForumID == 0 || len([]rune(title)) == 0 || len([]rune(title)) > 200 || len([]rune(content)) == 0 || len([]rune(content)) > 5000 {
			RespondFail(c, http.StatusBadRequest, 30001, "参数非法")
			return
		}

		var forum model.Forum
		if err := deps.DB.WithContext(c.Request.Context()).First(&forum, req.ForumID).Error; err != nil {
			RespondFail(c, http.StatusBadRequest, 30002, "贴吧不存在")
			return
		}

		p := model.Post{
			ForumID:      req.ForumID,
			AuthorID:     ai.UserID,
			Title:        title,
			Content:      content,
			Status:       "pending",
			LikeCount:    0,
			CommentCount: 0,
			ViewCount:    0,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		if err := deps.DB.WithContext(c.Request.Context()).Create(&p).Error; err != nil {
			RespondFail(c, http.StatusInternalServerError, 50000, "internal error")
			return
		}

		RespondOK(c, map[string]any{
			"postId": p.ID,
			"status": p.Status,
		})
	}
}

type updatePostReq struct {
	Title   *string `json:"title"`
	Content *string `json:"content"`
}

func updatePostHandler(deps Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		ai, ok := RequireAuth(c, deps)
		if !ok {
			return
		}

		postID, err := strconv.ParseUint(c.Param("postId"), 10, 64)
		if err != nil {
			RespondFail(c, http.StatusBadRequest, 30001, "参数非法")
			return
		}

		var req updatePostReq
		if err := c.ShouldBindJSON(&req); err != nil {
			RespondFail(c, http.StatusBadRequest, 30001, "参数非法")
			return
		}

		var p model.Post
		if err := deps.DB.WithContext(c.Request.Context()).First(&p, postID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				RespondFail(c, http.StatusNotFound, 30003, "帖子不存在")
				return
			}
			RespondFail(c, http.StatusInternalServerError, 50000, "internal error")
			return
		}

		if p.AuthorID != ai.UserID {
			RespondFail(c, http.StatusForbidden, 30004, "无权限")
			return
		}

		updates := map[string]any{}
		if req.Title != nil {
			title := strings.TrimSpace(*req.Title)
			if len([]rune(title)) == 0 || len([]rune(title)) > 200 {
				RespondFail(c, http.StatusBadRequest, 30001, "参数非法")
				return
			}
			updates["title"] = title
		}
		if req.Content != nil {
			content := strings.TrimSpace(*req.Content)
			if len([]rune(content)) == 0 || len([]rune(content)) > 5000 {
				RespondFail(c, http.StatusBadRequest, 30001, "参数非法")
				return
			}
			updates["content"] = content
		}
		if len(updates) == 0 {
			RespondFail(c, http.StatusBadRequest, 30001, "参数非法")
			return
		}
		updates["updated_at"] = time.Now()

		if err := deps.DB.WithContext(c.Request.Context()).Model(&model.Post{}).Where("id = ?", postID).Updates(updates).Error; err != nil {
			RespondFail(c, http.StatusInternalServerError, 50000, "internal error")
			return
		}
		RespondOK(c, map[string]any{})
	}
}

func deletePostHandler(deps Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		ai, ok := RequireAuth(c, deps)
		if !ok {
			return
		}

		postID, err := strconv.ParseUint(c.Param("postId"), 10, 64)
		if err != nil {
			RespondFail(c, http.StatusBadRequest, 30003, "帖子不存在")
			return
		}

		var p model.Post
		if err := deps.DB.WithContext(c.Request.Context()).First(&p, postID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				RespondFail(c, http.StatusNotFound, 30003, "帖子不存在")
				return
			}
			RespondFail(c, http.StatusInternalServerError, 50000, "internal error")
			return
		}

		if ai.Role != "admin" && ai.UserID != p.AuthorID {
			RespondFail(c, http.StatusForbidden, 30004, "无权限")
			return
		}

		if err := deps.DB.WithContext(c.Request.Context()).Delete(&p).Error; err != nil {
			RespondFail(c, http.StatusInternalServerError, 50000, "internal error")
			return
		}
		RespondOK(c, map[string]any{})
	}
}

func listForumPostsHandler(deps Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		forumID, err := strconv.ParseUint(c.Param("forumId"), 10, 64)
		if err != nil {
			RespondFail(c, http.StatusBadRequest, 30002, "贴吧不存在")
			return
		}
		var forum model.Forum
		if err := deps.DB.WithContext(c.Request.Context()).First(&forum, forumID).Error; err != nil {
			RespondFail(c, http.StatusBadRequest, 30002, "贴吧不存在")
			return
		}

		pageNum := parsePositiveIntWithDefault(c.Query("pageNum"), 1)
		pageSize := parsePositiveIntWithDefault(c.Query("pageSize"), 10)
		if pageNum <= 0 || pageSize <= 0 || pageSize > 100 {
			RespondFail(c, http.StatusBadRequest, 30001, "参数非法")
			return
		}

		ai, hasAuth := ParseAuth(c, deps)

		q := deps.DB.WithContext(c.Request.Context()).Model(&model.Post{}).Where("forum_id = ?", forumID)
		if !hasAuth {
			q = q.Where("status = ?", "visible")
		} else {
			q = q.Where(deps.DB.Where("status = ?", "visible").Or("author_id = ?", ai.UserID))
		}

		var total int64
		if err := q.Count(&total).Error; err != nil {
			RespondFail(c, http.StatusInternalServerError, 50000, "internal error")
			return
		}

		var rows []model.Post
		if err := q.Order("created_at desc").Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&rows).Error; err != nil {
			RespondFail(c, http.StatusInternalServerError, 50000, "internal error")
			return
		}

		list := make([]map[string]any, 0, len(rows))
		for _, p := range rows {
			list = append(list, map[string]any{
				"postId":       p.ID,
				"forumId":      p.ForumID,
				"title":        p.Title,
				"authorId":     p.AuthorID,
				"status":       p.Status,
				"likeCount":    p.LikeCount,
				"commentCount": p.CommentCount,
				"viewCount":    p.ViewCount,
				"createdAt":    p.CreatedAt.Format(time.RFC3339),
			})
		}

		RespondOK(c, map[string]any{
			"list":  list,
			"total": total,
		})
	}
}

func listVisiblePostsHandler(deps Deps) gin.HandlerFunc {
	type row struct {
		ID           uint64    `gorm:"column:id"`
		ForumID      uint64    `gorm:"column:forum_id"`
		ForumName    string    `gorm:"column:forum_name"`
		Title        string    `gorm:"column:title"`
		Content      string    `gorm:"column:content"`
		AuthorID     uint64    `gorm:"column:author_id"`
		Status       string    `gorm:"column:status"`
		LikeCount    int64     `gorm:"column:like_count"`
		CommentCount int64     `gorm:"column:comment_count"`
		ViewCount    int64     `gorm:"column:view_count"`
		CreatedAt    time.Time `gorm:"column:created_at"`
	}

	return func(c *gin.Context) {
		pageNum := parsePositiveIntWithDefault(c.Query("pageNum"), 1)
		pageSize := parsePositiveIntWithDefault(c.Query("pageSize"), 10)
		if pageNum <= 0 || pageSize <= 0 || pageSize > 100 {
			RespondFail(c, http.StatusBadRequest, 30001, "参数非法")
			return
		}

		base := deps.DB.WithContext(c.Request.Context()).
			Table("posts").
			Joins("JOIN forums ON forums.id = posts.forum_id AND forums.deleted_at IS NULL").
			Where("posts.deleted_at IS NULL").
			Where("posts.status = ?", "visible")

		var total int64
		if err := base.Count(&total).Error; err != nil {
			RespondFail(c, http.StatusInternalServerError, 50000, "internal error")
			return
		}

		var rows []row
		if err := base.Select(strings.Join([]string{
			"posts.id",
			"posts.forum_id",
			"forums.name AS forum_name",
			"posts.title",
			"posts.content",
			"posts.author_id",
			"posts.status",
			"posts.like_count",
			"posts.comment_count",
			"posts.view_count",
			"posts.created_at",
		}, ", ")).
			Order("posts.created_at desc").
			Limit(pageSize).
			Offset((pageNum - 1) * pageSize).
			Scan(&rows).Error; err != nil {
			RespondFail(c, http.StatusInternalServerError, 50000, "internal error")
			return
		}

		list := make([]map[string]any, 0, len(rows))
		for _, p := range rows {
			list = append(list, map[string]any{
				"postId":       p.ID,
				"forumId":      p.ForumID,
				"forumName":    p.ForumName,
				"title":        p.Title,
				"content":      p.Content,
				"authorId":     p.AuthorID,
				"status":       p.Status,
				"likeCount":    p.LikeCount,
				"commentCount": p.CommentCount,
				"viewCount":    p.ViewCount,
				"createdAt":    p.CreatedAt.Format(time.RFC3339),
			})
		}

		RespondOK(c, map[string]any{"list": list, "total": total})
	}
}

func getPostHandler(deps Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		postID, err := strconv.ParseUint(c.Param("postId"), 10, 64)
		if err != nil {
			RespondFail(c, http.StatusNotFound, 30003, "帖子不存在")
			return
		}

		var p model.Post
		if err := deps.DB.WithContext(c.Request.Context()).First(&p, postID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				RespondFail(c, http.StatusNotFound, 30003, "帖子不存在")
				return
			}
			RespondFail(c, http.StatusInternalServerError, 50000, "internal error")
			return
		}

		ai, hasAuth := ParseAuth(c, deps)
		if p.Status != "visible" {
			if !hasAuth || (ai.Role != "admin" && ai.UserID != p.AuthorID) {
				RespondFail(c, http.StatusForbidden, 30005, "不可见")
				return
			}
		}

		_ = deps.DB.WithContext(c.Request.Context()).
			Model(&model.Post{}).
			Where("id = ?", p.ID).
			UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error

		if err := deps.DB.WithContext(c.Request.Context()).First(&p, postID).Error; err != nil {
			RespondFail(c, http.StatusInternalServerError, 50000, "internal error")
			return
		}

		RespondOK(c, map[string]any{
			"postId":       p.ID,
			"forumId":      p.ForumID,
			"title":        p.Title,
			"content":      p.Content,
			"authorId":     p.AuthorID,
			"status":       p.Status,
			"likeCount":    p.LikeCount,
			"commentCount": p.CommentCount,
			"viewCount":    p.ViewCount,
			"createdAt":    p.CreatedAt.Format(time.RFC3339),
		})
	}
}

type createCommentReq struct {
	Content         string `json:"content"`
	ParentCommentID uint64 `json:"parentCommentId"`
}

func createCommentHandler(deps Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		ai, ok := RequireAuth(c, deps)
		if !ok {
			return
		}

		postID, err := strconv.ParseUint(c.Param("postId"), 10, 64)
		if err != nil {
			RespondFail(c, http.StatusBadRequest, 30003, "帖子不存在")
			return
		}

		var req createCommentReq
		if err := c.ShouldBindJSON(&req); err != nil {
			RespondFail(c, http.StatusBadRequest, 30011, "参数非法")
			return
		}
		content := strings.TrimSpace(req.Content)
		if len([]rune(content)) == 0 || len([]rune(content)) > 1000 {
			RespondFail(c, http.StatusBadRequest, 30011, "参数非法")
			return
		}

		var p model.Post
		if err := deps.DB.WithContext(c.Request.Context()).First(&p, postID).Error; err != nil {
			RespondFail(c, http.StatusBadRequest, 30003, "帖子不存在")
			return
		}
		if p.Status != "visible" && ai.Role != "admin" && ai.UserID != p.AuthorID {
			RespondFail(c, http.StatusForbidden, 30005, "不可见")
			return
		}

		var parentAuthorID uint64
		if req.ParentCommentID != 0 {
			var parent model.Comment
			if err := deps.DB.WithContext(c.Request.Context()).First(&parent, req.ParentCommentID).Error; err != nil || parent.PostID != postID {
				RespondFail(c, http.StatusBadRequest, 30012, "父评论不存在")
				return
			}
			parentAuthorID = parent.AuthorID
		}

		tx := deps.DB.WithContext(c.Request.Context()).Begin()
		cm := model.Comment{
			PostID:          postID,
			AuthorID:        ai.UserID,
			ParentCommentID: req.ParentCommentID,
			Content:         content,
			LikeCount:       0,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}
		if err := tx.Create(&cm).Error; err != nil {
			_ = tx.Rollback().Error
			RespondFail(c, http.StatusInternalServerError, 50000, "internal error")
			return
		}
		if err := tx.Model(&model.Post{}).Where("id = ?", postID).UpdateColumn("comment_count", gorm.Expr("comment_count + 1")).Error; err != nil {
			_ = tx.Rollback().Error
			RespondFail(c, http.StatusInternalServerError, 50000, "internal error")
			return
		}
		if err := tx.Commit().Error; err != nil {
			RespondFail(c, http.StatusInternalServerError, 50000, "internal error")
			return
		}

		notifyComment(deps, c, p.AuthorID, postID, cm.ID, ai.UserID)
		notifyReplyComment(deps, c, parentAuthorID, postID, cm.ID, ai.UserID)

		RespondOK(c, map[string]any{"commentId": cm.ID})
	}
}

func deleteCommentHandler(deps Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		ai, ok := RequireAuth(c, deps)
		if !ok {
			return
		}

		commentID, err := strconv.ParseUint(c.Param("commentId"), 10, 64)
		if err != nil {
			RespondFail(c, http.StatusNotFound, 30013, "评论不存在")
			return
		}

		var cm model.Comment
		if err := deps.DB.WithContext(c.Request.Context()).First(&cm, commentID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				RespondFail(c, http.StatusNotFound, 30013, "评论不存在")
				return
			}
			RespondFail(c, http.StatusInternalServerError, 50000, "internal error")
			return
		}

		if ai.Role != "admin" && ai.UserID != cm.AuthorID {
			RespondFail(c, http.StatusForbidden, 30004, "无权限")
			return
		}

		if err := deps.DB.WithContext(c.Request.Context()).Delete(&cm).Error; err != nil {
			RespondFail(c, http.StatusInternalServerError, 50000, "internal error")
			return
		}

		RespondOK(c, map[string]any{})
	}
}

func listCommentsHandler(deps Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		postID, err := strconv.ParseUint(c.Param("postId"), 10, 64)
		if err != nil {
			RespondFail(c, http.StatusBadRequest, 30003, "帖子不存在")
			return
		}

		var p model.Post
		if err := deps.DB.WithContext(c.Request.Context()).First(&p, postID).Error; err != nil {
			RespondFail(c, http.StatusBadRequest, 30003, "帖子不存在")
			return
		}

		if p.Status != "visible" {
			ai, hasAuth := ParseAuth(c, deps)
			if !hasAuth || (ai.Role != "admin" && ai.UserID != p.AuthorID) {
				RespondFail(c, http.StatusForbidden, 30005, "不可见")
				return
			}
		}

		pageNum := parsePositiveIntWithDefault(c.Query("pageNum"), 1)
		pageSize := parsePositiveIntWithDefault(c.Query("pageSize"), 10)
		if pageNum <= 0 || pageSize <= 0 || pageSize > 200 {
			RespondFail(c, http.StatusBadRequest, 30001, "参数非法")
			return
		}

		q := deps.DB.WithContext(c.Request.Context()).Model(&model.Comment{}).Where("post_id = ?", postID)
		var total int64
		if err := q.Count(&total).Error; err != nil {
			RespondFail(c, http.StatusInternalServerError, 50000, "internal error")
			return
		}

		var rows []model.Comment
		if err := q.Order("created_at asc").Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&rows).Error; err != nil {
			RespondFail(c, http.StatusInternalServerError, 50000, "internal error")
			return
		}

		list := make([]map[string]any, 0, len(rows))
		for _, cm := range rows {
			list = append(list, map[string]any{
				"commentId":       cm.ID,
				"postId":          cm.PostID,
				"authorId":        cm.AuthorID,
				"parentCommentId": cm.ParentCommentID,
				"content":         cm.Content,
				"likeCount":       cm.LikeCount,
				"createdAt":       cm.CreatedAt.Format(time.RFC3339),
			})
		}

		RespondOK(c, map[string]any{"list": list, "total": total})
	}
}

type likeReq struct {
	TargetType string `json:"targetType"`
	TargetID   uint64 `json:"targetId"`
}

func likeHandler(deps Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		ai, ok := RequireAuth(c, deps)
		if !ok {
			return
		}

		var req likeReq
		if err := c.ShouldBindJSON(&req); err != nil {
			RespondFail(c, http.StatusBadRequest, 30021, "参数非法")
			return
		}
		req.TargetType = strings.ToLower(strings.TrimSpace(req.TargetType))
		if (req.TargetType != "post" && req.TargetType != "comment") || req.TargetID == 0 {
			RespondFail(c, http.StatusBadRequest, 30021, "参数非法")
			return
		}

		tx := deps.DB.WithContext(c.Request.Context()).Begin()
		l := model.Like{
			UserID:     ai.UserID,
			TargetType: req.TargetType,
			TargetID:   req.TargetID,
			CreatedAt:  time.Now(),
		}
		if err := tx.Create(&l).Error; err != nil {
			_ = tx.Rollback().Error
			if isMySQLDup2(err) {
				RespondOK(c, map[string]any{})
				return
			}
			RespondFail(c, http.StatusInternalServerError, 50000, "internal error")
			return
		}

		var toUserID uint64
		switch req.TargetType {
		case "post":
			var p model.Post
			if err := tx.First(&p, req.TargetID).Error; err != nil {
				_ = tx.Rollback().Error
				RespondFail(c, http.StatusBadRequest, 30022, "目标不存在")
				return
			}
			toUserID = p.AuthorID
			if err := tx.Model(&model.Post{}).Where("id = ?", req.TargetID).UpdateColumn("like_count", gorm.Expr("like_count + 1")).Error; err != nil {
				_ = tx.Rollback().Error
				RespondFail(c, http.StatusInternalServerError, 50000, "internal error")
				return
			}
		case "comment":
			var cm model.Comment
			if err := tx.First(&cm, req.TargetID).Error; err != nil {
				_ = tx.Rollback().Error
				RespondFail(c, http.StatusBadRequest, 30022, "目标不存在")
				return
			}
			toUserID = cm.AuthorID
			if err := tx.Model(&model.Comment{}).Where("id = ?", req.TargetID).UpdateColumn("like_count", gorm.Expr("like_count + 1")).Error; err != nil {
				_ = tx.Rollback().Error
				RespondFail(c, http.StatusInternalServerError, 50000, "internal error")
				return
			}
		}

		if err := tx.Commit().Error; err != nil {
			RespondFail(c, http.StatusInternalServerError, 50000, "internal error")
			return
		}
		notifyLike(deps, c, toUserID, req.TargetType, req.TargetID, ai.UserID)
		RespondOK(c, map[string]any{})
	}
}

func feedHandler(deps Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		ai, ok := RequireAuth(c, deps)
		if !ok {
			return
		}

		pageNum := parsePositiveIntWithDefault(c.Query("pageNum"), 1)
		pageSize := parsePositiveIntWithDefault(c.Query("pageSize"), 10)
		if pageNum <= 0 || pageSize <= 0 || pageSize > 100 {
			RespondFail(c, http.StatusBadRequest, 30001, "参数非法")
			return
		}

		var forumIDs []uint64
		if err := deps.DB.WithContext(c.Request.Context()).
			Model(&model.Follow{}).
			Where("user_id = ? AND target_type = ? AND status = ?", ai.UserID, "forum", "active").
			Pluck("target_id", &forumIDs).Error; err != nil {
			RespondFail(c, http.StatusInternalServerError, 50000, "internal error")
			return
		}

		if len(forumIDs) == 0 {
			if err := deps.DB.WithContext(c.Request.Context()).
				Model(&model.ForumFollow{}).
				Where("user_id = ?", ai.UserID).
				Pluck("forum_id", &forumIDs).Error; err != nil {
				RespondFail(c, http.StatusInternalServerError, 50000, "internal error")
				return
			}
		}

		if len(forumIDs) == 0 {
			RespondOK(c, map[string]any{"list": []any{}, "total": 0})
			return
		}

		q := deps.DB.WithContext(c.Request.Context()).
			Model(&model.Post{}).
			Where("status = ?", "visible").
			Where("forum_id IN ?", forumIDs)

		var total int64
		if err := q.Count(&total).Error; err != nil {
			RespondFail(c, http.StatusInternalServerError, 50000, "internal error")
			return
		}

		var rows []model.Post
		if err := q.Order("created_at desc").
			Limit(pageSize).
			Offset((pageNum - 1) * pageSize).
			Find(&rows).Error; err != nil {
			RespondFail(c, http.StatusInternalServerError, 50000, "internal error")
			return
		}

		list := make([]map[string]any, 0, len(rows))
		for _, p := range rows {
			list = append(list, map[string]any{
				"postId":       p.ID,
				"forumId":      p.ForumID,
				"title":        p.Title,
				"authorId":     p.AuthorID,
				"status":       p.Status,
				"likeCount":    p.LikeCount,
				"commentCount": p.CommentCount,
				"viewCount":    p.ViewCount,
				"createdAt":    p.CreatedAt.Format(time.RFC3339),
			})
		}

		RespondOK(c, map[string]any{"list": list, "total": total})
	}
}

func isMySQLDup2(err error) bool {
	var myErr *mysql.MySQLError
	if errors.As(err, &myErr) {
		return myErr.Number == 1062
	}
	return strings.Contains(strings.ToLower(err.Error()), "duplicate")
}
