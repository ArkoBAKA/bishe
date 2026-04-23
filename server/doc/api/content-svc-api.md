# content-svc API

> 负责帖子、评论（楼中楼）、点赞（仅新增）、Feed（方案 A）。

## 接口清单

- `POST /api/v1/posts`
- `GET /api/v1/posts`
- `PUT /api/v1/posts/{postId}`
- `DELETE /api/v1/posts/{postId}`
- `GET /api/v1/forums/{forumId}/posts`
- `GET /api/v1/posts/{postId}`
- `POST /api/v1/posts/{postId}/comments`
- `DELETE /api/v1/comments/{commentId}`
- `GET /api/v1/posts/{postId}/comments`
- `POST /api/v1/likes`
- `GET /api/v1/feed`

---

## 发布帖子

- **Method**: `POST`
- **Path**: `/api/v1/posts`
- **Auth**: 需要
- **角色**: `user`
- **说明**: 发布帖子，初始可见状态为 `待审核`

### 请求参数

#### Body
| 字段 | 类型 | 必填 | 默认值 | 说明 |
|---|---|---|---|---|
| forumId | number | 是 | - | 贴吧 ID |
| title | string | 是 | - | 标题（1~200） |
| content | string | 是 | - | 正文（1~5000） |

### 响应

#### 成功示例
```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "postId": 10001,
    "status": "pending"
  }
}
```

#### data 字段结构
| 字段 | 类型 | 说明 |
|---|---|---|
| postId | number | 帖子 ID |
| status | string | 初始状态：`pending`（待审核） |

### 错误码
| code | 含义 | 触发条件 |
|---|---|---|
| 40101 | 未登录或 Token 无效 | 鉴权失败 |
| 30001 | 参数非法 | forumId/title/content 不合法 |
| 30002 | 贴吧不存在 | forumId 无对应记录 |

### 业务规则（实现必须遵守）
- 帖子状态采用单一可见状态流转：`pending`（待审核）/ `visible`（可见）/ `rejected`（审核拒绝）/ `hidden`（已隐藏）。
- 热榜输入字段 `likeCount`、`commentCount`、`viewCount` 必填，创建时初始化为 0。

### 实现备注（给后端 Agent）
- 所属服务：gateway + content-svc
- 涉及实体：posts
- 事务策略：单表写入
- 幂等约束：无

---

## 点赞（仅新增）

- **Method**: `POST`
- **Path**: `/api/v1/likes`
- **Auth**: 需要
- **角色**: `user`
- **说明**: 点赞帖子或评论（仅新增，不取消）

### 请求参数

#### Body
| 字段 | 类型 | 必填 | 默认值 | 说明 |
|---|---|---|---|---|
| targetType | string | 是 | - | `post` / `comment` |
| targetId | number | 是 | - | 目标 ID |

### 响应

#### 成功示例
```json
{
  "code": 0,
  "message": "ok",
  "data": {}
}
```

### 错误码
| code | 含义 | 触发条件 |
|---|---|---|
| 40101 | 未登录或 Token 无效 | 鉴权失败 |
| 30021 | 参数非法 | targetType/targetId 不合法 |
| 30022 | 目标不存在 | postId/commentId 不存在或不可见 |

### 业务规则（实现必须遵守）
- 幂等：同一用户对同一目标只能成功一次。
- 不提供取消点赞接口。

---

## 编辑帖子

- **Method**: `PUT`
- **Path**: `/api/v1/posts/{postId}`
- **Auth**: 需要
- **角色**: `user`
- **说明**: 编辑帖子（仅作者）

### 请求参数

#### Path
| 字段 | 类型 | 必填 | 默认值 | 说明 |
|---|---|---|---|---|
| postId | number | 是 | - | 帖子 ID |

#### Body
| 字段 | 类型 | 必填 | 默认值 | 说明 |
|---|---|---|---|---|
| title | string | 否 | - | 标题（1~200） |
| content | string | 否 | - | 正文（1~5000） |

### 响应

#### 成功示例
```json
{
  "code": 0,
  "message": "ok",
  "data": {}
}
```

### 错误码
| code | 含义 | 触发条件 |
|---|---|---|
| 40101 | 未登录或 Token 无效 | 鉴权失败 |
| 30001 | 参数非法 | 字段不合法或无任何变更 |
| 30003 | 帖子不存在 | postId 无对应记录 |
| 30004 | 无权限 | 非作者编辑 |

### 业务规则（实现必须遵守）
- 至少一个字段有变更才执行更新。
- 仅作者可编辑。

### 实现备注（给后端 Agent）
- 所属服务：gateway + content-svc
- 涉及实体：posts
- 事务策略：单表更新
- 幂等约束：重复提交同样内容不报错

---

## 删除帖子

- **Method**: `DELETE`
- **Path**: `/api/v1/posts/{postId}`
- **Auth**: 需要
- **角色**: `user` / `admin`
- **说明**: 删除帖子（软删；作者或管理员）

### 请求参数

#### Path
| 字段 | 类型 | 必填 | 默认值 | 说明 |
|---|---|---|---|---|
| postId | number | 是 | - | 帖子 ID |

### 响应

#### 成功示例
```json
{
  "code": 0,
  "message": "ok",
  "data": {}
}
```

### 错误码
| code | 含义 | 触发条件 |
|---|---|---|
| 40101 | 未登录或 Token 无效 | 鉴权失败 |
| 30003 | 帖子不存在 | postId 无对应记录 |
| 30004 | 无权限 | 非作者且非管理员 |

### 业务规则（实现必须遵守）
- 删除为软删（保留审计）。

---

## 可见帖子列表（全站）

- **Method**: `GET`
- **Path**: `/api/v1/posts`
- **Auth**: 不需要
- **角色**: `guest`
- **说明**: 分页查询全站审核通过（`visible`）的帖子列表，按创建时间倒序（最新在前）；每条帖子返回帖子详情字段 + 所属贴吧信息（`forumId` + `forumName`）

### 请求参数

#### Query
| 字段 | 类型 | 必填 | 默认值 | 说明 |
|---|---|---|---|---|
| pageNum | number | 否 | 1 | 页码 |
| pageSize | number | 否 | 10 | 每页条数 |

### 响应

#### 成功示例
```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "list": [
      {
        "postId": 10001,
        "forumId": 10001,
        "forumName": "Golang",
        "title": "Hello",
        "content": "正文",
        "authorId": 10001,
        "status": "visible",
        "likeCount": 0,
        "commentCount": 0,
        "viewCount": 0,
        "createdAt": "2026-04-21T10:30:00+08:00"
      }
    ],
    "total": 1
  }
}
```

#### data 字段结构
| 字段 | 类型 | 说明 |
|---|---|---|
| list | array | 帖子列表 |
| total | number | 总数 |

### 错误码
| code | 含义 | 触发条件 |
|---|---|---|
| 30001 | 参数非法 | pageNum/pageSize 非法 |

### 业务规则（实现必须遵守）
- 仅返回 `visible`（审核通过）帖子。
- 排序：`createdAt` 倒序（最新在前）。
- `forumName` 需由 `forumId` 反查得到。

---

## 贴吧帖子列表

- **Method**: `GET`
- **Path**: `/api/v1/forums/{forumId}/posts`
- **Auth**: 不需要
- **角色**: `guest`
- **说明**: 分页查询某贴吧下帖子列表（默认仅返回可见帖子）

### 请求参数

#### Path
| 字段 | 类型 | 必填 | 默认值 | 说明 |
|---|---|---|---|---|
| forumId | number | 是 | - | 贴吧 ID |

#### Query
| 字段 | 类型 | 必填 | 默认值 | 说明 |
|---|---|---|---|---|
| pageNum | number | 否 | 1 | 页码 |
| pageSize | number | 否 | 10 | 每页条数 |

### 响应

#### 成功示例
```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "list": [
      {
        "postId": 10001,
        "forumId": 10001,
        "title": "Hello",
        "authorId": 10001,
        "status": "visible",
        "likeCount": 0,
        "commentCount": 0,
        "viewCount": 0,
        "createdAt": "2026-04-21T10:30:00+08:00"
      }
    ],
    "total": 1
  }
}
```

#### data 字段结构
| 字段 | 类型 | 说明 |
|---|---|---|
| list | array | 帖子列表 |
| total | number | 总数 |

### 错误码
| code | 含义 | 触发条件 |
|---|---|---|
| 30002 | 贴吧不存在 | forumId 无对应记录 |
| 30001 | 参数非法 | pageNum/pageSize 非法 |

### 业务规则（实现必须遵守）
- 默认仅返回 `visible`。
- 若携带有效 token 且为作者，则可看到自己的非可见帖子。

---

## 帖子详情

- **Method**: `GET`
- **Path**: `/api/v1/posts/{postId}`
- **Auth**: 不需要
- **角色**: `guest`
- **说明**: 获取帖子详情（会增加 viewCount）

### 请求参数

#### Path
| 字段 | 类型 | 必填 | 默认值 | 说明 |
|---|---|---|---|---|
| postId | number | 是 | - | 帖子 ID |

### 响应

#### 成功示例
```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "postId": 10001,
    "forumId": 10001,
    "title": "Hello",
    "content": "正文",
    "authorId": 10001,
    "status": "visible",
    "likeCount": 0,
    "commentCount": 0,
    "viewCount": 1,
    "createdAt": "2026-04-21T10:30:00+08:00"
  }
}
```

#### data 字段结构
| 字段 | 类型 | 说明 |
|---|---|---|
| postId | number | 帖子 ID |
| forumId | number | 贴吧 ID |
| title | string | 标题 |
| content | string | 正文 |
| authorId | number | 作者 ID |
| status | string | 状态 |
| likeCount | number | 点赞数 |
| commentCount | number | 评论数 |
| viewCount | number | 浏览数 |
| createdAt | string | 创建时间（RFC3339） |

### 错误码
| code | 含义 | 触发条件 |
|---|---|---|
| 30003 | 帖子不存在 | postId 无对应记录 |
| 30005 | 不可见 | 非可见状态且无权限访问 |

---

## 发评论（楼中楼）

- **Method**: `POST`
- **Path**: `/api/v1/posts/{postId}/comments`
- **Auth**: 需要
- **角色**: `user`
- **说明**: 对帖子发表评论；支持楼中楼（parentCommentId）

### 请求参数

#### Path
| 字段 | 类型 | 必填 | 默认值 | 说明 |
|---|---|---|---|---|
| postId | number | 是 | - | 帖子 ID |

#### Body
| 字段 | 类型 | 必填 | 默认值 | 说明 |
|---|---|---|---|---|
| content | string | 是 | - | 评论内容（1~1000） |
| parentCommentId | number | 否 | - | 父评论 ID（楼中楼） |

### 响应

#### 成功示例
```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "commentId": 20001
  }
}
```

#### data 字段结构
| 字段 | 类型 | 说明 |
|---|---|---|
| commentId | number | 评论 ID |

### 错误码
| code | 含义 | 触发条件 |
|---|---|---|
| 40101 | 未登录或 Token 无效 | 鉴权失败 |
| 30003 | 帖子不存在 | postId 无对应记录 |
| 30011 | 参数非法 | content 为空/过长，parentCommentId 不合法 |
| 30012 | 父评论不存在 | parentCommentId 不存在或不属于该 postId |

### 业务规则（实现必须遵守）
- 创建评论后，帖子 `commentCount` +1。

---

## 删除评论

- **Method**: `DELETE`
- **Path**: `/api/v1/comments/{commentId}`
- **Auth**: 需要
- **角色**: `user` / `admin`
- **说明**: 删除评论（软删；作者或管理员）

### 请求参数

#### Path
| 字段 | 类型 | 必填 | 默认值 | 说明 |
|---|---|---|---|---|
| commentId | number | 是 | - | 评论 ID |

### 响应

#### 成功示例
```json
{
  "code": 0,
  "message": "ok",
  "data": {}
}
```

### 错误码
| code | 含义 | 触发条件 |
|---|---|---|
| 40101 | 未登录或 Token 无效 | 鉴权失败 |
| 30013 | 评论不存在 | commentId 无对应记录 |
| 30004 | 无权限 | 非作者且非管理员 |

---

## 评论列表

- **Method**: `GET`
- **Path**: `/api/v1/posts/{postId}/comments`
- **Auth**: 不需要
- **角色**: `guest`
- **说明**: 分页查询评论列表（按创建时间正序）

### 请求参数

#### Path
| 字段 | 类型 | 必填 | 默认值 | 说明 |
|---|---|---|---|---|
| postId | number | 是 | - | 帖子 ID |

#### Query
| 字段 | 类型 | 必填 | 默认值 | 说明 |
|---|---|---|---|---|
| pageNum | number | 否 | 1 | 页码 |
| pageSize | number | 否 | 10 | 每页条数 |

### 响应

#### 成功示例
```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "list": [
      {
        "commentId": 20001,
        "postId": 10001,
        "authorId": 10001,
        "parentCommentId": 0,
        "content": "评论内容",
        "likeCount": 0,
        "createdAt": "2026-04-21T10:30:00+08:00"
      }
    ],
    "total": 1
  }
}
```

### 错误码
| code | 含义 | 触发条件 |
|---|---|---|
| 30003 | 帖子不存在 | postId 无对应记录 |
| 30001 | 参数非法 | pageNum/pageSize 非法 |

---

## Feed（方案 A）

- **Method**: `GET`
- **Path**: `/api/v1/feed`
- **Auth**: 需要
- **角色**: `user`
- **说明**: 查询“我关注的贴吧”下的帖子时间线（按发布时间倒序分页）

### 请求参数

#### Query
| 字段 | 类型 | 必填 | 默认值 | 说明 |
|---|---|---|---|---|
| pageNum | number | 否 | 1 | 页码 |
| pageSize | number | 否 | 10 | 每页条数 |

### 响应

#### 成功示例
```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "list": [],
    "total": 0
  }
}
```

### 业务规则（实现必须遵守）
- 仅返回 `visible` 帖子。
- 关注关系以 `forum_follows(userId, forumId)` 为准（后续可由 relation-svc 负责维护）。

---

## 建表语句（MySQL 8）

```sql
CREATE TABLE IF NOT EXISTS `posts` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `forum_id` BIGINT UNSIGNED NOT NULL,
  `author_id` BIGINT UNSIGNED NOT NULL,
  `title` VARCHAR(200) NOT NULL,
  `content` VARCHAR(5000) NOT NULL,
  `status` VARCHAR(16) NOT NULL DEFAULT 'pending',
  `like_count` BIGINT NOT NULL DEFAULT 0,
  `comment_count` BIGINT NOT NULL DEFAULT 0,
  `view_count` BIGINT NOT NULL DEFAULT 0,
  `created_at` DATETIME NOT NULL,
  `updated_at` DATETIME NOT NULL,
  `deleted_at` DATETIME NULL,
  PRIMARY KEY (`id`),
  KEY `idx_posts_forum_id_created_at` (`forum_id`, `created_at`),
  KEY `idx_posts_author_id_created_at` (`author_id`, `created_at`),
  KEY `idx_posts_status_created_at` (`status`, `created_at`),
  KEY `idx_posts_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `comments` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `post_id` BIGINT UNSIGNED NOT NULL,
  `author_id` BIGINT UNSIGNED NOT NULL,
  `parent_comment_id` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `content` VARCHAR(1000) NOT NULL,
  `like_count` BIGINT NOT NULL DEFAULT 0,
  `created_at` DATETIME NOT NULL,
  `updated_at` DATETIME NOT NULL,
  `deleted_at` DATETIME NULL,
  PRIMARY KEY (`id`),
  KEY `idx_comments_post_id_created_at` (`post_id`, `created_at`),
  KEY `idx_comments_author_id_created_at` (`author_id`, `created_at`),
  KEY `idx_comments_parent_comment_id` (`parent_comment_id`),
  KEY `idx_comments_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `likes` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT UNSIGNED NOT NULL,
  `target_type` VARCHAR(16) NOT NULL,
  `target_id` BIGINT UNSIGNED NOT NULL,
  `created_at` DATETIME NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_likes_user_target` (`user_id`, `target_type`, `target_id`),
  KEY `idx_likes_target` (`target_type`, `target_id`),
  KEY `idx_likes_user_id_created_at` (`user_id`, `created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `forum_follows` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT UNSIGNED NOT NULL,
  `forum_id` BIGINT UNSIGNED NOT NULL,
  `created_at` DATETIME NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_forum_follows_user_forum` (`user_id`, `forum_id`),
  KEY `idx_forum_follows_forum_id_created_at` (`forum_id`, `created_at`),
  KEY `idx_forum_follows_user_id_created_at` (`user_id`, `created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
```

其余接口按 `api/README.md` 的固定模板补充详细字段、示例与错误码。
