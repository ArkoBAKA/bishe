package http

import (
	"github.com/gin-gonic/gin"
)

func notifyLike(deps Deps, c *gin.Context, toUserID uint64, targetType string, targetID uint64, fromUserID uint64) {
	if toUserID == 0 || toUserID == fromUserID {
		return
	}
	createNotification(deps, c, toUserID, "like", "收到点赞", "你的内容被点赞了", map[string]any{
		"targetType": targetType,
		"targetId":   targetID,
		"fromUserId": fromUserID,
	})
}

func notifyComment(deps Deps, c *gin.Context, toUserID uint64, postID uint64, commentID uint64, fromUserID uint64) {
	if toUserID == 0 || toUserID == fromUserID {
		return
	}
	createNotification(deps, c, toUserID, "comment", "收到评论", "你的帖子收到了新评论", map[string]any{
		"postId":     postID,
		"commentId":  commentID,
		"fromUserId": fromUserID,
	})
}

func notifyReplyComment(deps Deps, c *gin.Context, toUserID uint64, postID uint64, commentID uint64, fromUserID uint64) {
	if toUserID == 0 || toUserID == fromUserID {
		return
	}
	createNotification(deps, c, toUserID, "comment", "评论被回复", "你的评论收到了回复", map[string]any{
		"postId":     postID,
		"commentId":  commentID,
		"fromUserId": fromUserID,
	})
}
