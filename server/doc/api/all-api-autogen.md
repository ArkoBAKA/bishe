# API 文档（自动生成）

统一响应结构：

```json
{"code":0,"message":"ok","data":{}}
```

## 管理员删除评论

- **Method**: `DELETE`
- **Path**: `/api/v1/admin/comments/{commentId}`
- **Auth**: 需要
- **角色**: `admin`
- **说明**: 删除评论（软删）

### 响应

#### 成功示例
```json
{"code":0,"message":"ok","data":{}}
```

## 管理员删除帖子

- **Method**: `DELETE`
- **Path**: `/api/v1/admin/posts/{postId}`
- **Auth**: 需要
- **角色**: `admin`
- **说明**: 删除帖子（软删）

### 响应

#### 成功示例
```json
{"code":0,"message":"ok","data":{}}
```

## 删除评论

- **Method**: `DELETE`
- **Path**: `/api/v1/comments/{commentId}`
- **Auth**: 需要
- **角色**: `user/admin`
- **说明**: 删除评论（软删；作者或管理员）

### 响应

#### 成功示例
```json
{"code":0,"message":"ok","data":{}}
```

## 取消关注

- **Method**: `DELETE`
- **Path**: `/api/v1/follows`
- **Auth**: 需要
- **角色**: `user`
- **说明**: 取消关注（幂等；使用状态位不物理删除）

### 请求参数

#### Body
| 字段 | 类型 | 必填 | 默认值 | 说明 |
|---|---|---|---|---|
| targetType | string | 是 |  | forum/user |
| targetId | number | 是 |  | 目标 ID |

### 响应

#### 成功示例
```json
{"code":0,"message":"ok","data":{}}
```

## 删除帖子

- **Method**: `DELETE`
- **Path**: `/api/v1/posts/{postId}`
- **Auth**: 需要
- **角色**: `user/admin`
- **说明**: 删除帖子（软删；作者或管理员）

### 响应

#### 成功示例
```json
{"code":0,"message":"ok","data":{}}
```

## 删除文件（软删）

- **Method**: `DELETE`
- **Path**: `/api/v1/upload/:fileId`
- **Auth**: 需要
- **角色**: `user/admin`
- **说明**: 软删除文件元数据（仅本人或管理员）

### 响应

#### 成功示例
```json
{"code":0,"message":"ok","data":{}}
```

## 待审核帖子列表

- **Method**: `GET`
- **Path**: `/api/v1/admin/posts/pending`
- **Auth**: 需要
- **角色**: `admin`
- **说明**: 分页查询待审核帖子列表（status=pending）

### 请求参数

#### Query
| 字段 | 类型 | 必填 | 默认值 | 说明 |
|---|---|---|---|---|
| pageNum | number | 否 | 1 | 页码 |
| pageSize | number | 否 | 10 | 每页条数 |

### 响应

#### 成功示例
```json
{"code":0,"message":"ok","data":{}}
```

#### data 字段结构
| 字段 | 类型 | 说明 |
|---|---|---|
| list | array | 帖子列表 |
| total | number | 总数 |

## 举报列表（管理员）

- **Method**: `GET`
- **Path**: `/api/v1/admin/reports`
- **Auth**: 需要
- **角色**: `admin`
- **说明**: 分页查询举报审计单列表

### 请求参数

#### Query
| 字段 | 类型 | 必填 | 默认值 | 说明 |
|---|---|---|---|---|
| pageNum | number | 否 | 1 | 页码 |
| pageSize | number | 否 | 10 | 每页条数 |
| status | string | 否 |  | pending/processed |

### 响应

#### 成功示例
```json
{"code":0,"message":"ok","data":{}}
```

#### data 字段结构
| 字段 | 类型 | 说明 |
|---|---|---|
| list | array | 举报单列表 |
| total | number | 总数 |

## 接口文档（自动生成）

- **Method**: `GET`
- **Path**: `/api/v1/docs`
- **Auth**: 不需要
- **角色**: `guest`
- **说明**: 返回自动生成的接口文档。format=md/json

### 请求参数

#### Query
| 字段 | 类型 | 必填 | 默认值 | 说明 |
|---|---|---|---|---|
| format | string | 否 | md | md 或 json |

### 响应

#### 成功示例
```json
{"code":0,"message":"ok","data":{}}
```

## Feed（方案 A）

- **Method**: `GET`
- **Path**: `/api/v1/feed`
- **Auth**: 需要
- **角色**: `user`
- **说明**: 查询我关注的贴吧下的帖子时间线（按发布时间倒序分页）

### 请求参数

#### Query
| 字段 | 类型 | 必填 | 默认值 | 说明 |
|---|---|---|---|---|
| pageNum | number | 否 | 1 | 页码 |
| pageSize | number | 否 | 10 | 每页条数 |

### 响应

#### 成功示例
```json
{"code":0,"message":"ok","data":{}}
```

#### data 字段结构
| 字段 | 类型 | 说明 |
|---|---|---|
| list | array | 帖子列表 |
| total | number | 总数 |

## 我的关注列表

- **Method**: `GET`
- **Path**: `/api/v1/follows/me`
- **Auth**: 需要
- **角色**: `user`
- **说明**: 获取我的关注列表（默认 active，createdAt desc）

### 请求参数

#### Query
| 字段 | 类型 | 必填 | 默认值 | 说明 |
|---|---|---|---|---|
| pageNum | number | 否 | 1 | 页码 |
| pageSize | number | 否 | 10 | 每页数量 |
| targetType | string | 否 |  | forum/user |

### 响应

#### 成功示例
```json
{"code":0,"message":"ok","data":{}}
```

#### data 字段结构
| 字段 | 类型 | 说明 |
|---|---|---|
| list | array | 关注列表 |
| total | number | 总数 |

## 贴吧列表

- **Method**: `GET`
- **Path**: `/api/v1/forums`
- **Auth**: 不需要
- **角色**: `guest`
- **说明**: 分页查询贴吧列表（按创建时间倒序）

### 请求参数

#### Query
| 字段 | 类型 | 必填 | 默认值 | 说明 |
|---|---|---|---|---|
| pageNum | number | 否 | 1 | 页码，从 1 开始 |
| pageSize | number | 否 | 10 | 每页条数 |
| keyword | string | 否 |  | 名称关键字（模糊匹配） |

### 响应

#### 成功示例
```json
{"code":0,"message":"ok","data":{}}
```

#### data 字段结构
| 字段 | 类型 | 说明 |
|---|---|---|
| list | array | 贴吧列表 |
| total | number | 总数 |

## 贴吧详情

- **Method**: `GET`
- **Path**: `/api/v1/forums/{forumId}`
- **Auth**: 不需要
- **角色**: `guest`
- **说明**: 获取贴吧详情

### 响应

#### 成功示例
```json
{"code":0,"message":"ok","data":{}}
```

## 贴吧帖子列表

- **Method**: `GET`
- **Path**: `/api/v1/forums/{forumId}/posts`
- **Auth**: 不需要
- **角色**: `guest`
- **说明**: 分页查询某贴吧下帖子列表（默认仅返回可见帖子）

### 请求参数

#### Query
| 字段 | 类型 | 必填 | 默认值 | 说明 |
|---|---|---|---|---|
| pageNum | number | 否 | 1 | 页码 |
| pageSize | number | 否 | 10 | 每页条数 |

### 响应

#### 成功示例
```json
{"code":0,"message":"ok","data":{}}
```

#### data 字段结构
| 字段 | 类型 | 说明 |
|---|---|---|
| list | array | 帖子列表 |
| total | number | 总数 |

## 通知列表

- **Method**: `GET`
- **Path**: `/api/v1/notifications`
- **Auth**: 需要
- **角色**: `user/admin`
- **说明**: 获取我的通知列表（按创建时间倒序）

### 请求参数

#### Query
| 字段 | 类型 | 必填 | 默认值 | 说明 |
|---|---|---|---|---|
| pageNum | number | 否 | 1 | 页码 |
| pageSize | number | 否 | 10 | 每页数量 |
| isRead | boolean | 否 |  | 是否已读筛选 |

### 响应

#### 成功示例
```json
{"code":0,"message":"ok","data":{}}
```

#### data 字段结构
| 字段 | 类型 | 说明 |
|---|---|---|
| list | array | 通知列表 |
| total | number | 总数 |

## 帖子详情

- **Method**: `GET`
- **Path**: `/api/v1/posts/{postId}`
- **Auth**: 不需要
- **角色**: `guest`
- **说明**: 获取帖子详情（会增加 viewCount）

### 响应

#### 成功示例
```json
{"code":0,"message":"ok","data":{}}
```

## 评论列表

- **Method**: `GET`
- **Path**: `/api/v1/posts/{postId}/comments`
- **Auth**: 不需要
- **角色**: `guest`
- **说明**: 分页查询评论列表（按创建时间正序）

### 请求参数

#### Query
| 字段 | 类型 | 必填 | 默认值 | 说明 |
|---|---|---|---|---|
| pageNum | number | 否 | 1 | 页码 |
| pageSize | number | 否 | 10 | 每页条数 |

### 响应

#### 成功示例
```json
{"code":0,"message":"ok","data":{}}
```

#### data 字段结构
| 字段 | 类型 | 说明 |
|---|---|---|
| list | array | 评论列表 |
| total | number | 总数 |

## 静态文件访问

- **Method**: `GET`
- **Path**: `/api/v1/static/:bucket/*filepath`
- **Auth**: public 不需要；private 需要
- **角色**: `guest/user`
- **说明**: 访问上传后的文件。public 可匿名访问，private 需要 Bearer token

### 响应

#### 成功示例
```json
{"code":0,"message":"ok","data":{}}
```

## 查询文件元数据

- **Method**: `GET`
- **Path**: `/api/v1/upload/:fileId`
- **Auth**: 需要
- **角色**: `user`
- **说明**: 按 fileId 查询上传文件元数据

### 响应

#### 成功示例
```json
{"code":0,"message":"ok","data":{}}
```

## 当前用户

- **Method**: `GET`
- **Path**: `/api/v1/users/me`
- **Auth**: 需要
- **角色**: `user/admin`
- **说明**: 查询当前登录用户信息

### 响应

#### 成功示例
```json
{"code":0,"message":"ok","data":{}}
```

## 用户公开信息

- **Method**: `GET`
- **Path**: `/api/v1/users/{userId}`
- **Auth**: 不需要
- **角色**: `all`
- **说明**: 获取指定用户公开信息

### 响应

#### 成功示例
```json
{"code":0,"message":"ok","data":{}}
```

## 健康检查

- **Method**: `GET`
- **Path**: `/health`
- **Auth**: 不需要
- **角色**: `guest`
- **说明**: 检查 MySQL/Redis 连通性

### 响应

#### 成功示例
```json
{"code":0,"message":"ok","data":{}}
```

#### data 字段结构
| 字段 | 类型 | 说明 |
|---|---|---|
| ok | bool | mysql 与 redis 均可用 |
| mysql | bool | MySQL 是否可用 |
| redis | bool | Redis 是否可用 |
| mysqlError | string | MySQL 错误（可选） |
| redisError | string | Redis 错误（可选） |

## 查询用户（示例）

- **Method**: `GET`
- **Path**: `/users/:id`
- **Auth**: 不需要
- **角色**: `guest`
- **说明**: 示例接口：按 id 查询用户

### 响应

#### 成功示例
```json
{"code":0,"message":"ok","data":{}}
```

## 创建关注关系

- **Method**: `POST`
- **Path**: `/api/v1/follows`
- **Auth**: 需要
- **角色**: `user`
- **说明**: 关注贴吧或用户（幂等；取消后可恢复）

### 请求参数

#### Body
| 字段 | 类型 | 必填 | 默认值 | 说明 |
|---|---|---|---|---|
| targetType | string | 是 |  | forum/user |
| targetId | number | 是 |  | 目标 ID |

### 响应

#### 成功示例
```json
{"code":0,"message":"ok","data":{}}
```

## 创建贴吧

- **Method**: `POST`
- **Path**: `/api/v1/forums`
- **Auth**: 需要
- **角色**: `user`
- **说明**: 创建贴吧

### 请求参数

#### Body
| 字段 | 类型 | 必填 | 默认值 | 说明 |
|---|---|---|---|---|
| name | string | 是 |  | 贴吧名称（全局唯一） |
| description | string | 否 |  | 贴吧简介 |
| coverUrl | string | 否 |  | 封面图地址 |

### 响应

#### 成功示例
```json
{"code":0,"message":"ok","data":{}}
```

#### data 字段结构
| 字段 | 类型 | 说明 |
|---|---|---|
| forumId | number | 新创建贴吧 ID |

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
| targetType | string | 是 |  | post/comment |
| targetId | number | 是 |  | 目标 ID |

### 响应

#### 成功示例
```json
{"code":0,"message":"ok","data":{}}
```

## 发布帖子

- **Method**: `POST`
- **Path**: `/api/v1/posts`
- **Auth**: 需要
- **角色**: `user`
- **说明**: 发布帖子，初始可见状态为 pending（待审核）

### 请求参数

#### Body
| 字段 | 类型 | 必填 | 默认值 | 说明 |
|---|---|---|---|---|
| forumId | number | 是 |  | 贴吧 ID |
| title | string | 是 |  | 标题（1~200） |
| content | string | 是 |  | 正文（1~5000） |

### 响应

#### 成功示例
```json
{"code":0,"message":"ok","data":{}}
```

#### data 字段结构
| 字段 | 类型 | 说明 |
|---|---|---|
| postId | number | 帖子 ID |
| status | string | 初始状态 pending |

## 发评论（楼中楼）

- **Method**: `POST`
- **Path**: `/api/v1/posts/{postId}/comments`
- **Auth**: 需要
- **角色**: `user`
- **说明**: 对帖子发表评论；支持楼中楼（parentCommentId）

### 请求参数

#### Body
| 字段 | 类型 | 必填 | 默认值 | 说明 |
|---|---|---|---|---|
| content | string | 是 |  | 评论内容（1~1000） |
| parentCommentId | number | 否 |  | 父评论 ID（楼中楼） |

### 响应

#### 成功示例
```json
{"code":0,"message":"ok","data":{}}
```

#### data 字段结构
| 字段 | 类型 | 说明 |
|---|---|---|
| commentId | number | 评论 ID |

## 提交举报

- **Method**: `POST`
- **Path**: `/api/v1/reports`
- **Auth**: 需要
- **角色**: `user`
- **说明**: 举报帖子/评论/用户，生成待处理审计单

### 请求参数

#### Body
| 字段 | 类型 | 必填 | 默认值 | 说明 |
|---|---|---|---|---|
| targetType | string | 是 |  | post/comment/user |
| targetId | number | 是 |  | 目标 ID |
| reason | string | 是 |  | 简短原因（1~64） |
| detail | string | 否 |  | 详细描述（<=1000） |

### 响应

#### 成功示例
```json
{"code":0,"message":"ok","data":{}}
```

#### data 字段结构
| 字段 | 类型 | 说明 |
|---|---|---|
| reportId | number | 举报单 ID |

## 单文件上传

- **Method**: `POST`
- **Path**: `/api/v1/upload`
- **Auth**: 需要
- **角色**: `user`
- **说明**: multipart/form-data 上传单文件，返回可访问 URL

### 请求参数

#### Body
| 字段 | 类型 | 必填 | 默认值 | 说明 |
|---|---|---|---|---|
| file | file | 是 |  | 文件本体 |
| bucket | string | 否 | public | public/private |
| scene | string | 否 |  | 业务场景标记 |

### 响应

#### 成功示例
```json
{"code":0,"message":"ok","data":{}}
```

#### data 字段结构
| 字段 | 类型 | 说明 |
|---|---|---|
| fileId | number | 文件元数据 ID |
| url | string | 静态访问 URL |

## 用户登录

- **Method**: `POST`
- **Path**: `/api/v1/users/login`
- **Auth**: 不需要
- **角色**: `guest`
- **说明**: 账号密码登录并签发 JWT

### 请求参数

#### Body
| 字段 | 类型 | 必填 | 默认值 | 说明 |
|---|---|---|---|---|
| account | string | 是 |  | 登录账号 |
| password | string | 是 |  | 登录密码 |

### 响应

#### 成功示例
```json
{"code":0,"message":"ok","data":{}}
```

#### data 字段结构
| 字段 | 类型 | 说明 |
|---|---|---|
| token | string | JWT |
| tokenType | string | 固定 Bearer |
| expiresIn | number | 秒 |
| user | object | 当前用户摘要 |

## 用户登出

- **Method**: `POST`
- **Path**: `/api/v1/users/logout`
- **Auth**: 需要
- **角色**: `user/admin`
- **说明**: 登出当前会话（本项目不做 token 黑名单）

### 响应

#### 成功示例
```json
{"code":0,"message":"ok","data":{}}
```

## 用户注册

- **Method**: `POST`
- **Path**: `/api/v1/users/register`
- **Auth**: 不需要
- **角色**: `guest`
- **说明**: 账号密码注册，支持头像文件上传

### 请求参数

#### Body
| 字段 | 类型 | 必填 | 默认值 | 说明 |
|---|---|---|---|---|
| account | string | 是 |  | 登录账号（唯一） |
| password | string | 是 |  | 明文密码（服务端哈希） |
| nickname | string | 否 |  | 昵称，不传则由服务端生成默认昵称 |
| avatarFile | file | 否 |  | 头像文件，multipart/form-data 上传 |

### 响应

#### 成功示例
```json
{"code":0,"message":"ok","data":{}}
```

#### data 字段结构
| 字段 | 类型 | 说明 |
|---|---|---|
| userId | number | 新注册用户 ID |
| avatarUrl | string | 头像 URL（未上传则为空字符串） |

## 创建用户（示例）

- **Method**: `POST`
- **Path**: `/users`
- **Auth**: 不需要
- **角色**: `guest`
- **说明**: 示例接口：创建用户记录（username 唯一）

### 请求参数

#### Body
| 字段 | 类型 | 必填 | 默认值 | 说明 |
|---|---|---|---|---|
| username | string | 是 |  | 用户名 |

### 响应

#### 成功示例
```json
{"code":0,"message":"ok","data":{}}
```

## 审核帖子并更新状态

- **Method**: `PUT`
- **Path**: `/api/v1/admin/posts/{postId}/review`
- **Auth**: 需要
- **角色**: `admin`
- **说明**: 审核帖子并更新可见状态

### 请求参数

#### Body
| 字段 | 类型 | 必填 | 默认值 | 说明 |
|---|---|---|---|---|
| action | string | 是 |  | approve/reject/hide |
| reviewRemark | string | 否 |  | 处理备注 |

### 响应

#### 成功示例
```json
{"code":0,"message":"ok","data":{}}
```

#### data 字段结构
| 字段 | 类型 | 说明 |
|---|---|---|
| status | string | 更新后的帖子状态 |

## 处理举报并执行动作

- **Method**: `PUT`
- **Path**: `/api/v1/admin/reports/{reportId}/process`
- **Auth**: 需要
- **角色**: `admin`
- **说明**: 处理举报审计单并执行动作

### 请求参数

#### Body
| 字段 | 类型 | 必填 | 默认值 | 说明 |
|---|---|---|---|---|
| action | string | 是 |  | close/deletePost/deleteComment/hidePost/banUser |
| processRemark | string | 否 |  | 处理备注 |
| banUntil | string | 否 |  | RFC3339（banUser 可用） |
| durationSeconds | number | 否 |  | 秒（banUser 可用） |

### 响应

#### 成功示例
```json
{"code":0,"message":"ok","data":{}}
```

## 封禁用户

- **Method**: `PUT`
- **Path**: `/api/v1/admin/users/{userId}/ban`
- **Auth**: 需要
- **角色**: `admin`
- **说明**: 封禁用户（设置 status=banned 与 banUntil）

### 请求参数

#### Body
| 字段 | 类型 | 必填 | 默认值 | 说明 |
|---|---|---|---|---|
| banUntil | string | 否 |  | 封禁截止时间（RFC3339） |
| durationSeconds | number | 否 |  | 封禁时长（秒） |
| remark | string | 否 |  | 备注 |

### 响应

#### 成功示例
```json
{"code":0,"message":"ok","data":{}}
```

#### data 字段结构
| 字段 | 类型 | 说明 |
|---|---|---|
| status | string | banned |
| banUntil | string | RFC3339 |

## 编辑贴吧

- **Method**: `PUT`
- **Path**: `/api/v1/forums/{forumId}`
- **Auth**: 需要
- **角色**: `user`
- **说明**: 编辑指定贴吧（仅创建者/owner）

### 响应

#### 成功示例
```json
{"code":0,"message":"ok","data":{}}
```

## 全部通知已读

- **Method**: `PUT`
- **Path**: `/api/v1/notifications/read-all`
- **Auth**: 需要
- **角色**: `user/admin`
- **说明**: 将我所有未读通知标记为已读（幂等）

### 响应

#### 成功示例
```json
{"code":0,"message":"ok","data":{}}
```

#### data 字段结构
| 字段 | 类型 | 说明 |
|---|---|---|
| affected | number | 被更新为已读的数量 |

## 标记通知已读

- **Method**: `PUT`
- **Path**: `/api/v1/notifications/{notificationId}/read`
- **Auth**: 需要
- **角色**: `user/admin`
- **说明**: 将指定通知标记为已读（幂等）

### 响应

#### 成功示例
```json
{"code":0,"message":"ok","data":{}}
```

## 编辑帖子

- **Method**: `PUT`
- **Path**: `/api/v1/posts/{postId}`
- **Auth**: 需要
- **角色**: `user`
- **说明**: 编辑帖子（仅作者）

### 响应

#### 成功示例
```json
{"code":0,"message":"ok","data":{}}
```

## 修改密码

- **Method**: `PUT`
- **Path**: `/api/v1/users/me/password`
- **Auth**: 需要
- **角色**: `user/admin`
- **说明**: 修改当前用户密码

### 请求参数

#### Body
| 字段 | 类型 | 必填 | 默认值 | 说明 |
|---|---|---|---|---|
| oldPassword | string | 是 |  | 旧密码 |
| newPassword | string | 是 |  | 新密码 |

### 响应

#### 成功示例
```json
{"code":0,"message":"ok","data":{}}
```

## 修改资料

- **Method**: `PUT`
- **Path**: `/api/v1/users/me/profile`
- **Auth**: 需要
- **角色**: `user/admin`
- **说明**: 修改当前用户资料，支持头像文件上传

### 请求参数

#### Body
| 字段 | 类型 | 必填 | 默认值 | 说明 |
|---|---|---|---|---|
| nickname | string | 否 |  | 昵称 |
| avatarFile | file | 否 |  | 新头像文件，multipart/form-data 上传 |
| bio | string | 否 |  | 简介 |

### 响应

#### 成功示例
```json
{"code":0,"message":"ok","data":{}}
```

