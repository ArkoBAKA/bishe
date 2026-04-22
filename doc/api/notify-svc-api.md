# notify-svc API

> 负责站内通知查询与已读。采用拉取模型，无实时推送。

## 接口清单

- `GET /api/v1/notifications`
- `PUT /api/v1/notifications/{notificationId}/read`
- `PUT /api/v1/notifications/read-all`

---

## `GET /api/v1/notifications`

- **Auth**: 需要
- **角色**: `user` / `admin`
- **说明**: 获取我的通知列表（按创建时间倒序）

### 请求参数

#### Query
| 字段 | 类型 | 必填 | 默认值 | 说明 |
|---|---|---|---|---|
| pageNum | number | 否 | 1 | 页码 |
| pageSize | number | 否 | 10 | 每页数量 |
| isRead | boolean | 否 | - | 是否已读筛选 |

### 业务规则（实现必须遵守）
- 排序固定 `createdAt desc`。

---

## 其余接口模板

其余接口按 `api/README.md` 的固定模板补充详细字段、示例与错误码。
