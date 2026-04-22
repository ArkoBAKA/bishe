# content-svc API

> 负责帖子、评论（楼中楼）、点赞（仅新增）、Feed（方案 A）。

## 接口清单

- `POST /api/v1/posts`
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

## `POST /api/v1/posts`

- **Auth**: 需要
- **角色**: `user`
- **说明**: 发布帖子，初始可见状态为 `待审核`

### 业务规则（实现必须遵守）
- 帖子状态采用单一可见状态流转：`待审核` / `可见` / `审核拒绝` / `已隐藏`。
- 热榜输入字段 `likeCount`、`commentCount`、`viewCount` 必填，创建时初始化为 0。

---

## `POST /api/v1/likes`

- **Auth**: 需要
- **角色**: `user`
- **说明**: 点赞帖子或评论（仅新增，不取消）

### 请求参数

#### Body
| 字段 | 类型 | 必填 | 默认值 | 说明 |
|---|---|---|---|---|
| targetType | string | 是 | - | `post` / `comment` |
| targetId | number | 是 | - | 目标 ID |

### 业务规则（实现必须遵守）
- 幂等：同一用户对同一目标只能成功一次。

---

## 其余接口模板

其余接口按 `api/README.md` 的固定模板补充详细字段、示例与错误码。
