# forum-svc API

> 负责贴吧创建、编辑、列表、详情。

## 接口清单

- `POST /api/v1/forums`
- `PUT /api/v1/forums/{forumId}`
- `GET /api/v1/forums`
- `GET /api/v1/forums/{forumId}`

---

## `POST /api/v1/forums`

- **Auth**: 需要
- **角色**: `user`
- **说明**: 创建贴吧

### 请求参数

#### Body
| 字段 | 类型 | 必填 | 默认值 | 说明 |
|---|---|---|---|---|
| name | string | 是 | - | 贴吧名称（全局唯一） |
| description | string | 否 | - | 贴吧简介 |
| coverUrl | string | 否 | - | 封面图地址 |

### 业务规则（实现必须遵守）
- `name` 全局唯一。

---

## 其余接口模板

其余接口按 `api/README.md` 的固定模板补充详细字段、示例与错误码。
