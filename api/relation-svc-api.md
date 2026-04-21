# relation-svc API

> 负责关注贴吧、关注用户、取消关注与我的关注列表。

## 接口清单

- `POST /api/v1/follows`
- `DELETE /api/v1/follows`
- `GET /api/v1/follows/me`

---

## `POST /api/v1/follows`

- **Auth**: 需要
- **角色**: `user`
- **说明**: 创建关注关系

### 请求参数

#### Body
| 字段 | 类型 | 必填 | 默认值 | 说明 |
|---|---|---|---|---|
| targetType | string | 是 | - | `forum` / `user` |
| targetId | number | 是 | - | 目标 ID |

### 业务规则（实现必须遵守）
- 关系幂等，同一目标重复关注不重复创建。
- 取消关注使用状态位，不物理删除。

---

## 其余接口模板

其余接口按 `api/README.md` 的固定模板补充详细字段、示例与错误码。
