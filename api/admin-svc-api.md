# admin-svc API

> 负责内容审核与举报审计闭环处理。

## 接口清单

- `GET /api/v1/admin/posts/pending`
- `PUT /api/v1/admin/posts/{postId}/review`
- `POST /api/v1/reports`
- `GET /api/v1/admin/reports`
- `PUT /api/v1/admin/reports/{reportId}/process`
- `DELETE /api/v1/admin/posts/{postId}`
- `DELETE /api/v1/admin/comments/{commentId}`
- `PUT /api/v1/admin/users/{userId}/ban`

---

## `PUT /api/v1/admin/posts/{postId}/review`

- **Auth**: 需要
- **角色**: `admin`
- **说明**: 审核帖子并更新可见状态

### 请求参数

#### Body
| 字段 | 类型 | 必填 | 默认值 | 说明 |
|---|---|---|---|---|
| action | string | 是 | - | `approve` / `reject` / `hide` |
| reviewRemark | string | 否 | - | 处理备注 |

### 业务规则（实现必须遵守）
- 审核动作必须映射到帖子可见状态流转。

---

## `PUT /api/v1/admin/reports/{reportId}/process`

- **Auth**: 需要
- **角色**: `admin`
- **说明**: 处理举报审计单并执行动作

### 业务规则（实现必须遵守）
- 报告处理需记录处理人、处理时间、处理动作、处理备注。

---

## 其余接口模板

其余接口按 `api/README.md` 的固定模板补充详细字段、示例与错误码。
