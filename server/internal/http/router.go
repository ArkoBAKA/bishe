package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"server/internal/config"
)

type Deps struct {
	DB     *gorm.DB
	Redis  *redis.Client
	Config config.Config
}

func NewRouter(deps Deps) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	r.GET("/health", WithDoc(APIDoc{
		Name:        "健康检查",
		Method:      "GET",
		Path:        "/health",
		Auth:        "不需要",
		Role:        "guest",
		Description: "检查 MySQL/Redis 连通性",
		Data: []FieldDoc{
			{Name: "ok", Type: "bool", Desc: "mysql 与 redis 均可用"},
			{Name: "mysql", Type: "bool", Desc: "MySQL 是否可用"},
			{Name: "redis", Type: "bool", Desc: "Redis 是否可用"},
			{Name: "mysqlError", Type: "string", Desc: "MySQL 错误（可选）"},
			{Name: "redisError", Type: "string", Desc: "Redis 错误（可选）"},
		},
	}, healthHandler(deps)))

	r.POST("/users", WithDoc(APIDoc{
		Name:        "创建用户（示例）",
		Method:      "POST",
		Path:        "/users",
		Auth:        "不需要",
		Role:        "guest",
		Description: "示例接口：创建用户记录（username 唯一）",
		Body: []FieldDoc{
			{Name: "username", Type: "string", Required: true, Desc: "用户名"},
		},
	}, createUserHandler(deps)))
	r.GET("/users/:id", WithDoc(APIDoc{
		Name:        "查询用户（示例）",
		Method:      "GET",
		Path:        "/users/:id",
		Auth:        "不需要",
		Role:        "guest",
		Description: "示例接口：按 id 查询用户",
		Query:       nil,
	}, getUserHandler(deps)))

	v1 := r.Group("/api/v1")
	{
		v1.GET("/docs", WithDoc(APIDoc{
			Name:        "接口文档（自动生成）",
			Method:      "GET",
			Path:        "/api/v1/docs",
			Auth:        "不需要",
			Role:        "guest",
			Description: "返回自动生成的接口文档。format=md/json",
			Query: []FieldDoc{
				{Name: "format", Type: "string", Required: false, Default: "md", Desc: "md 或 json"},
			},
		}, docsHandler()))

		v1.GET("/static/:bucket/*filepath", WithDoc(APIDoc{
			Name:        "静态文件访问",
			Method:      "GET",
			Path:        "/api/v1/static/:bucket/*filepath",
			Auth:        "public 不需要；private 需要",
			Role:        "guest/user",
			Description: "访问上传后的文件。public 可匿名访问，private 需要 Bearer token",
		}, staticHandler(deps)))

		v1.POST("/upload", WithDoc(APIDoc{
			Name:        "单文件上传",
			Method:      "POST",
			Path:        "/api/v1/upload",
			Auth:        "需要",
			Role:        "user",
			Description: "multipart/form-data 上传单文件，返回可访问 URL",
			Body: []FieldDoc{
				{Name: "file", Type: "file", Required: true, Desc: "文件本体"},
				{Name: "bucket", Type: "string", Required: false, Default: "public", Desc: "public/private"},
				{Name: "scene", Type: "string", Required: false, Default: "", Desc: "业务场景标记"},
			},
			Data: []FieldDoc{
				{Name: "fileId", Type: "number", Desc: "文件元数据 ID"},
				{Name: "url", Type: "string", Desc: "静态访问 URL"},
			},
		}, uploadHandler(deps)))

		v1.GET("/upload/:fileId", WithDoc(APIDoc{
			Name:        "查询文件元数据",
			Method:      "GET",
			Path:        "/api/v1/upload/:fileId",
			Auth:        "需要",
			Role:        "user",
			Description: "按 fileId 查询上传文件元数据",
		}, getUploadMetaHandler(deps)))

		v1.DELETE("/upload/:fileId", WithDoc(APIDoc{
			Name:        "删除文件（软删）",
			Method:      "DELETE",
			Path:        "/api/v1/upload/:fileId",
			Auth:        "需要",
			Role:        "user/admin",
			Description: "软删除文件元数据（仅本人或管理员）",
		}, deleteUploadHandler(deps)))

		v1.POST("/users/register", WithDoc(APIDoc{
			Name:        "用户注册",
			Method:      "POST",
			Path:        "/api/v1/users/register",
			Auth:        "不需要",
			Role:        "guest",
			Description: "账号密码注册，支持头像文件上传",
			Body: []FieldDoc{
				{Name: "account", Type: "string", Required: true, Desc: "登录账号（唯一）"},
				{Name: "password", Type: "string", Required: true, Desc: "明文密码（服务端哈希）"},
				{Name: "nickname", Type: "string", Required: false, Desc: "昵称，不传则由服务端生成默认昵称"},
				{Name: "avatarFile", Type: "file", Required: false, Desc: "头像文件，multipart/form-data 上传"},
			},
			Data: []FieldDoc{
				{Name: "userId", Type: "number", Desc: "新注册用户 ID"},
				{Name: "avatarUrl", Type: "string", Desc: "头像 URL（未上传则为空字符串）"},
			},
		}, userRegisterHandler(deps)))

		v1.POST("/users/login", WithDoc(APIDoc{
			Name:        "用户登录",
			Method:      "POST",
			Path:        "/api/v1/users/login",
			Auth:        "不需要",
			Role:        "guest",
			Description: "账号密码登录并签发 JWT",
			Body: []FieldDoc{
				{Name: "account", Type: "string", Required: true, Desc: "登录账号"},
				{Name: "password", Type: "string", Required: true, Desc: "登录密码"},
			},
			Data: []FieldDoc{
				{Name: "token", Type: "string", Desc: "JWT"},
				{Name: "tokenType", Type: "string", Desc: "固定 Bearer"},
				{Name: "expiresIn", Type: "number", Desc: "秒"},
				{Name: "user", Type: "object", Desc: "当前用户摘要"},
			},
		}, userLoginHandler(deps)))

		v1.POST("/users/logout", WithDoc(APIDoc{
			Name:        "用户登出",
			Method:      "POST",
			Path:        "/api/v1/users/logout",
			Auth:        "需要",
			Role:        "user/admin",
			Description: "登出当前会话（本项目不做 token 黑名单）",
		}, userLogoutHandler(deps)))

		v1.GET("/users/me", WithDoc(APIDoc{
			Name:        "当前用户",
			Method:      "GET",
			Path:        "/api/v1/users/me",
			Auth:        "需要",
			Role:        "user/admin",
			Description: "查询当前登录用户信息",
		}, userMeHandler(deps)))

		v1.GET("/users/:userId", WithDoc(APIDoc{
			Name:        "用户公开信息",
			Method:      "GET",
			Path:        "/api/v1/users/{userId}",
			Auth:        "不需要",
			Role:        "all",
			Description: "获取指定用户公开信息",
		}, userPublicHandler(deps)))

		v1.PUT("/users/me/profile", WithDoc(APIDoc{
			Name:        "修改资料",
			Method:      "PUT",
			Path:        "/api/v1/users/me/profile",
			Auth:        "需要",
			Role:        "user/admin",
			Description: "修改当前用户资料，支持头像文件上传",
			Body: []FieldDoc{
				{Name: "nickname", Type: "string", Required: false, Desc: "昵称"},
				{Name: "avatarFile", Type: "file", Required: false, Desc: "新头像文件，multipart/form-data 上传"},
				{Name: "bio", Type: "string", Required: false, Desc: "简介"},
			},
		}, userUpdateProfileHandler(deps)))

		v1.PUT("/users/me/password", WithDoc(APIDoc{
			Name:        "修改密码",
			Method:      "PUT",
			Path:        "/api/v1/users/me/password",
			Auth:        "需要",
			Role:        "user/admin",
			Description: "修改当前用户密码",
			Body: []FieldDoc{
				{Name: "oldPassword", Type: "string", Required: true, Desc: "旧密码"},
				{Name: "newPassword", Type: "string", Required: true, Desc: "新密码"},
			},
		}, userUpdatePasswordHandler(deps)))

		v1.POST("/forums", WithDoc(APIDoc{
			Name:        "创建贴吧",
			Method:      "POST",
			Path:        "/api/v1/forums",
			Auth:        "需要",
			Role:        "user",
			Description: "创建贴吧",
			Body: []FieldDoc{
				{Name: "name", Type: "string", Required: true, Desc: "贴吧名称（全局唯一）"},
				{Name: "description", Type: "string", Required: false, Desc: "贴吧简介"},
				{Name: "coverUrl", Type: "string", Required: false, Desc: "封面图地址"},
			},
			Data: []FieldDoc{
				{Name: "forumId", Type: "number", Desc: "新创建贴吧 ID"},
			},
		}, createForumHandler(deps)))

		v1.PUT("/forums/:forumId", WithDoc(APIDoc{
			Name:        "编辑贴吧",
			Method:      "PUT",
			Path:        "/api/v1/forums/{forumId}",
			Auth:        "需要",
			Role:        "user",
			Description: "编辑指定贴吧（仅创建者/owner）",
		}, updateForumHandler(deps)))

		v1.GET("/forums", WithDoc(APIDoc{
			Name:        "贴吧列表",
			Method:      "GET",
			Path:        "/api/v1/forums",
			Auth:        "不需要",
			Role:        "guest",
			Description: "分页查询贴吧列表（按创建时间倒序）",
			Query: []FieldDoc{
				{Name: "pageNum", Type: "number", Required: false, Default: "1", Desc: "页码，从 1 开始"},
				{Name: "pageSize", Type: "number", Required: false, Default: "10", Desc: "每页条数"},
				{Name: "keyword", Type: "string", Required: false, Default: "", Desc: "名称关键字（模糊匹配）"},
			},
			Data: []FieldDoc{
				{Name: "list", Type: "array", Desc: "贴吧列表"},
				{Name: "total", Type: "number", Desc: "总数"},
			},
		}, listForumsHandler(deps)))

		v1.GET("/forums/:forumId", WithDoc(APIDoc{
			Name:        "贴吧详情",
			Method:      "GET",
			Path:        "/api/v1/forums/{forumId}",
			Auth:        "不需要",
			Role:        "guest",
			Description: "获取贴吧详情",
		}, getForumHandler(deps)))

		v1.POST("/posts", WithDoc(APIDoc{
			Name:        "发布帖子",
			Method:      "POST",
			Path:        "/api/v1/posts",
			Auth:        "需要",
			Role:        "user",
			Description: "发布帖子，初始可见状态为 pending（待审核）",
			Body: []FieldDoc{
				{Name: "forumId", Type: "number", Required: true, Desc: "贴吧 ID"},
				{Name: "title", Type: "string", Required: true, Desc: "标题（1~200）"},
				{Name: "content", Type: "string", Required: true, Desc: "正文（1~5000）"},
			},
			Data: []FieldDoc{
				{Name: "postId", Type: "number", Desc: "帖子 ID"},
				{Name: "status", Type: "string", Desc: "初始状态 pending"},
			},
		}, createPostHandler(deps)))

		v1.PUT("/posts/:postId", WithDoc(APIDoc{
			Name:        "编辑帖子",
			Method:      "PUT",
			Path:        "/api/v1/posts/{postId}",
			Auth:        "需要",
			Role:        "user",
			Description: "编辑帖子（仅作者）",
		}, updatePostHandler(deps)))

		v1.DELETE("/posts/:postId", WithDoc(APIDoc{
			Name:        "删除帖子",
			Method:      "DELETE",
			Path:        "/api/v1/posts/{postId}",
			Auth:        "需要",
			Role:        "user/admin",
			Description: "删除帖子（软删；作者或管理员）",
		}, deletePostHandler(deps)))

		v1.GET("/posts", WithDoc(APIDoc{
			Name:        "可见帖子列表（全站）",
			Method:      "GET",
			Path:        "/api/v1/posts",
			Auth:        "不需要",
			Role:        "guest",
			Description: "分页查询全站审核通过（visible）的帖子列表，按创建时间倒序（最新在前）",
			Query: []FieldDoc{
				{Name: "pageNum", Type: "number", Required: false, Default: "1", Desc: "页码"},
				{Name: "pageSize", Type: "number", Required: false, Default: "10", Desc: "每页条数"},
			},
			Data: []FieldDoc{
				{Name: "list", Type: "array", Desc: "帖子列表"},
				{Name: "total", Type: "number", Desc: "总数"},
			},
		}, listVisiblePostsHandler(deps)))

		v1.GET("/forums/:forumId/posts", WithDoc(APIDoc{
			Name:        "贴吧帖子列表",
			Method:      "GET",
			Path:        "/api/v1/forums/{forumId}/posts",
			Auth:        "不需要",
			Role:        "guest",
			Description: "分页查询某贴吧下帖子列表（默认仅返回可见帖子）",
			Query: []FieldDoc{
				{Name: "pageNum", Type: "number", Required: false, Default: "1", Desc: "页码"},
				{Name: "pageSize", Type: "number", Required: false, Default: "10", Desc: "每页条数"},
			},
			Data: []FieldDoc{
				{Name: "list", Type: "array", Desc: "帖子列表"},
				{Name: "total", Type: "number", Desc: "总数"},
			},
		}, listForumPostsHandler(deps)))

		v1.GET("/posts/:postId", WithDoc(APIDoc{
			Name:        "帖子详情",
			Method:      "GET",
			Path:        "/api/v1/posts/{postId}",
			Auth:        "不需要",
			Role:        "guest",
			Description: "获取帖子详情（会增加 viewCount）",
		}, getPostHandler(deps)))

		v1.POST("/posts/:postId/comments", WithDoc(APIDoc{
			Name:        "发评论（楼中楼）",
			Method:      "POST",
			Path:        "/api/v1/posts/{postId}/comments",
			Auth:        "需要",
			Role:        "user",
			Description: "对帖子发表评论；支持楼中楼（parentCommentId）",
			Body: []FieldDoc{
				{Name: "content", Type: "string", Required: true, Desc: "评论内容（1~1000）"},
				{Name: "parentCommentId", Type: "number", Required: false, Desc: "父评论 ID（楼中楼）"},
			},
			Data: []FieldDoc{
				{Name: "commentId", Type: "number", Desc: "评论 ID"},
			},
		}, createCommentHandler(deps)))

		v1.DELETE("/comments/:commentId", WithDoc(APIDoc{
			Name:        "删除评论",
			Method:      "DELETE",
			Path:        "/api/v1/comments/{commentId}",
			Auth:        "需要",
			Role:        "user/admin",
			Description: "删除评论（软删；作者或管理员）",
		}, deleteCommentHandler(deps)))

		v1.GET("/posts/:postId/comments", WithDoc(APIDoc{
			Name:        "评论列表",
			Method:      "GET",
			Path:        "/api/v1/posts/{postId}/comments",
			Auth:        "不需要",
			Role:        "guest",
			Description: "分页查询评论列表（按创建时间正序）",
			Query: []FieldDoc{
				{Name: "pageNum", Type: "number", Required: false, Default: "1", Desc: "页码"},
				{Name: "pageSize", Type: "number", Required: false, Default: "10", Desc: "每页条数"},
			},
			Data: []FieldDoc{
				{Name: "list", Type: "array", Desc: "评论列表"},
				{Name: "total", Type: "number", Desc: "总数"},
			},
		}, listCommentsHandler(deps)))

		v1.POST("/likes", WithDoc(APIDoc{
			Name:        "点赞（仅新增）",
			Method:      "POST",
			Path:        "/api/v1/likes",
			Auth:        "需要",
			Role:        "user",
			Description: "点赞帖子或评论（仅新增，不取消）",
			Body: []FieldDoc{
				{Name: "targetType", Type: "string", Required: true, Desc: "post/comment"},
				{Name: "targetId", Type: "number", Required: true, Desc: "目标 ID"},
			},
		}, likeHandler(deps)))

		v1.GET("/feed", WithDoc(APIDoc{
			Name:        "Feed（方案 A）",
			Method:      "GET",
			Path:        "/api/v1/feed",
			Auth:        "需要",
			Role:        "user",
			Description: "查询我关注的贴吧下的帖子时间线（按发布时间倒序分页）",
			Query: []FieldDoc{
				{Name: "pageNum", Type: "number", Required: false, Default: "1", Desc: "页码"},
				{Name: "pageSize", Type: "number", Required: false, Default: "10", Desc: "每页条数"},
			},
			Data: []FieldDoc{
				{Name: "list", Type: "array", Desc: "帖子列表"},
				{Name: "total", Type: "number", Desc: "总数"},
			},
		}, feedHandler(deps)))

		v1.POST("/follows", WithDoc(APIDoc{
			Name:        "创建关注关系",
			Method:      "POST",
			Path:        "/api/v1/follows",
			Auth:        "需要",
			Role:        "user",
			Description: "关注贴吧或用户（幂等；取消后可恢复）",
			Body: []FieldDoc{
				{Name: "targetType", Type: "string", Required: true, Desc: "forum/user"},
				{Name: "targetId", Type: "number", Required: true, Desc: "目标 ID"},
			},
		}, createFollowHandler(deps)))

		v1.DELETE("/follows", WithDoc(APIDoc{
			Name:        "取消关注",
			Method:      "DELETE",
			Path:        "/api/v1/follows",
			Auth:        "需要",
			Role:        "user",
			Description: "取消关注（幂等；使用状态位不物理删除）",
			Body: []FieldDoc{
				{Name: "targetType", Type: "string", Required: true, Desc: "forum/user"},
				{Name: "targetId", Type: "number", Required: true, Desc: "目标 ID"},
			},
		}, cancelFollowHandler(deps)))

		v1.GET("/follows/me", WithDoc(APIDoc{
			Name:        "我的关注列表",
			Method:      "GET",
			Path:        "/api/v1/follows/me",
			Auth:        "需要",
			Role:        "user",
			Description: "获取我的关注列表（默认 active，createdAt desc）",
			Query: []FieldDoc{
				{Name: "pageNum", Type: "number", Required: false, Default: "1", Desc: "页码"},
				{Name: "pageSize", Type: "number", Required: false, Default: "10", Desc: "每页数量"},
				{Name: "targetType", Type: "string", Required: false, Default: "", Desc: "forum/user"},
			},
			Data: []FieldDoc{
				{Name: "list", Type: "array", Desc: "关注列表"},
				{Name: "total", Type: "number", Desc: "总数"},
			},
		}, listMyFollowsHandler(deps)))

		v1.GET("/notifications", WithDoc(APIDoc{
			Name:        "通知列表",
			Method:      "GET",
			Path:        "/api/v1/notifications",
			Auth:        "需要",
			Role:        "user/admin",
			Description: "获取我的通知列表（按创建时间倒序）",
			Query: []FieldDoc{
				{Name: "pageNum", Type: "number", Required: false, Default: "1", Desc: "页码"},
				{Name: "pageSize", Type: "number", Required: false, Default: "10", Desc: "每页数量"},
				{Name: "isRead", Type: "boolean", Required: false, Default: "", Desc: "是否已读筛选"},
			},
			Data: []FieldDoc{
				{Name: "list", Type: "array", Desc: "通知列表"},
				{Name: "total", Type: "number", Desc: "总数"},
			},
		}, listNotificationsHandler(deps)))

		v1.PUT("/notifications/:notificationId/read", WithDoc(APIDoc{
			Name:        "标记通知已读",
			Method:      "PUT",
			Path:        "/api/v1/notifications/{notificationId}/read",
			Auth:        "需要",
			Role:        "user/admin",
			Description: "将指定通知标记为已读（幂等）",
		}, markNotificationReadHandler(deps)))

		v1.PUT("/notifications/read-all", WithDoc(APIDoc{
			Name:        "全部通知已读",
			Method:      "PUT",
			Path:        "/api/v1/notifications/read-all",
			Auth:        "需要",
			Role:        "user/admin",
			Description: "将我所有未读通知标记为已读（幂等）",
			Data: []FieldDoc{
				{Name: "affected", Type: "number", Desc: "被更新为已读的数量"},
			},
		}, markAllNotificationsReadHandler(deps)))

		v1.POST("/reports", WithDoc(APIDoc{
			Name:        "提交举报",
			Method:      "POST",
			Path:        "/api/v1/reports",
			Auth:        "需要",
			Role:        "user",
			Description: "举报帖子/评论/用户，生成待处理审计单",
			Body: []FieldDoc{
				{Name: "targetType", Type: "string", Required: true, Desc: "post/comment/user"},
				{Name: "targetId", Type: "number", Required: true, Desc: "目标 ID"},
				{Name: "reason", Type: "string", Required: true, Desc: "简短原因（1~64）"},
				{Name: "detail", Type: "string", Required: false, Desc: "详细描述（<=1000）"},
			},
			Data: []FieldDoc{
				{Name: "reportId", Type: "number", Desc: "举报单 ID"},
			},
		}, createReportHandler(deps)))

		v1.GET("/admin/posts/pending", WithDoc(APIDoc{
			Name:        "待审核帖子列表",
			Method:      "GET",
			Path:        "/api/v1/admin/posts/pending",
			Auth:        "需要",
			Role:        "admin",
			Description: "分页查询待审核帖子列表（status=pending）",
			Query: []FieldDoc{
				{Name: "pageNum", Type: "number", Required: false, Default: "1", Desc: "页码"},
				{Name: "pageSize", Type: "number", Required: false, Default: "10", Desc: "每页条数"},
			},
			Data: []FieldDoc{
				{Name: "list", Type: "array", Desc: "帖子列表"},
				{Name: "total", Type: "number", Desc: "总数"},
			},
		}, adminListPendingPostsHandler(deps)))

		v1.PUT("/admin/posts/:postId/review", WithDoc(APIDoc{
			Name:        "审核帖子并更新状态",
			Method:      "PUT",
			Path:        "/api/v1/admin/posts/{postId}/review",
			Auth:        "需要",
			Role:        "admin",
			Description: "审核帖子并更新可见状态",
			Body: []FieldDoc{
				{Name: "action", Type: "string", Required: true, Desc: "approve/reject/hide"},
				{Name: "reviewRemark", Type: "string", Required: false, Desc: "处理备注"},
			},
			Data: []FieldDoc{
				{Name: "status", Type: "string", Desc: "更新后的帖子状态"},
			},
		}, adminReviewPostHandler(deps)))

		v1.GET("/admin/reports", WithDoc(APIDoc{
			Name:        "举报列表（管理员）",
			Method:      "GET",
			Path:        "/api/v1/admin/reports",
			Auth:        "需要",
			Role:        "admin",
			Description: "分页查询举报审计单列表",
			Query: []FieldDoc{
				{Name: "pageNum", Type: "number", Required: false, Default: "1", Desc: "页码"},
				{Name: "pageSize", Type: "number", Required: false, Default: "10", Desc: "每页条数"},
				{Name: "status", Type: "string", Required: false, Default: "", Desc: "pending/processed"},
			},
			Data: []FieldDoc{
				{Name: "list", Type: "array", Desc: "举报单列表"},
				{Name: "total", Type: "number", Desc: "总数"},
			},
		}, adminListReportsHandler(deps)))

		v1.PUT("/admin/reports/:reportId/process", WithDoc(APIDoc{
			Name:        "处理举报并执行动作",
			Method:      "PUT",
			Path:        "/api/v1/admin/reports/{reportId}/process",
			Auth:        "需要",
			Role:        "admin",
			Description: "处理举报审计单并执行动作",
			Body: []FieldDoc{
				{Name: "action", Type: "string", Required: true, Desc: "close/deletePost/deleteComment/hidePost/banUser"},
				{Name: "processRemark", Type: "string", Required: false, Desc: "处理备注"},
				{Name: "banUntil", Type: "string", Required: false, Desc: "RFC3339（banUser 可用）"},
				{Name: "durationSeconds", Type: "number", Required: false, Desc: "秒（banUser 可用）"},
			},
		}, adminProcessReportHandler(deps)))

		v1.DELETE("/admin/posts/:postId", WithDoc(APIDoc{
			Name:        "管理员删除帖子",
			Method:      "DELETE",
			Path:        "/api/v1/admin/posts/{postId}",
			Auth:        "需要",
			Role:        "admin",
			Description: "删除帖子（软删）",
		}, adminDeletePostHandler(deps)))

		v1.DELETE("/admin/comments/:commentId", WithDoc(APIDoc{
			Name:        "管理员删除评论",
			Method:      "DELETE",
			Path:        "/api/v1/admin/comments/{commentId}",
			Auth:        "需要",
			Role:        "admin",
			Description: "删除评论（软删）",
		}, adminDeleteCommentHandler(deps)))

		v1.PUT("/admin/users/:userId/ban", WithDoc(APIDoc{
			Name:        "封禁用户",
			Method:      "PUT",
			Path:        "/api/v1/admin/users/{userId}/ban",
			Auth:        "需要",
			Role:        "admin",
			Description: "封禁用户（设置 status=banned 与 banUntil）",
			Body: []FieldDoc{
				{Name: "banUntil", Type: "string", Required: false, Desc: "封禁截止时间（RFC3339）"},
				{Name: "durationSeconds", Type: "number", Required: false, Desc: "封禁时长（秒）"},
				{Name: "remark", Type: "string", Required: false, Desc: "备注"},
			},
			Data: []FieldDoc{
				{Name: "status", Type: "string", Desc: "banned"},
				{Name: "banUntil", Type: "string", Desc: "RFC3339"},
			},
		}, adminBanUserHandler(deps)))
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	})

	return r
}
