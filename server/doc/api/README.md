# API 文档编写规范（多 Agent 协作）

本目录用于前后端与多 Agent 协作对接。所有服务 API 文档必须遵循本规范，避免实现歧义。

---

## 1. 全局约定

- Base Path：`/api/v1`
- 认证方式：`Authorization: Bearer <token>`
- 字段命名：`camelCase`
- 时间格式：`RFC3339`（例如 `2026-04-21T10:30:00+08:00`）
- 分页参数：`pageNum`、`pageSize`（默认 `10`）
- 列表返回：统一 `list + total`
- 点赞规则：仅新增，不提供取消点赞；重复请求幂等

---

## 2. 统一响应结构

```json
{
  "code": 0,
  "message": "ok",
  "data": {}
}
```

- `code=0` 表示成功，非 0 表示失败
- 失败时：HTTP 状态码为 4xx/5xx，`data` 固定返回空对象 `{}`（便于前端统一处理）
- 列表建议：

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

---

## 3. 接口模板（每个接口必须使用）

~~~md
## 接口名

- **Method**: `GET/POST/PUT/DELETE`
- **Path**: `/api/v1/...`
- **Auth**: 需要 / 不需要
- **角色**: `guest` / `user` / `admin`
- **说明**: 一句话描述业务目标

### 请求参数

#### Query
| 字段 | 类型 | 必填 | 默认值 | 说明 |
|---|---|---|---|---|

#### Body
| 字段 | 类型 | 必填 | 默认值 | 说明 |
|---|---|---|---|---|

### 响应

#### 成功示例
{
  "code": 0,
  "message": "ok",
  "data": {}
}

#### data 字段结构
| 字段 | 类型 | 说明 |
|---|---|---|

### 错误码
| code | 含义 | 触发条件 |
|---|---|---|

### 业务规则（实现必须遵守）
- 规则 1
- 规则 2

### 实现备注（给后端 Agent）
- 所属服务：gateway + xxx-svc
- 涉及实体：xxx
- 事务策略：同步/异步
- 幂等约束：唯一键或去重规则
~~~

---

## 4. 文档拆分

- `user-svc-api.md`
- `forum-svc-api.md`
- `content-svc-api.md`
- `relation-svc-api.md`
- `notify-svc-api.md`
- `admin-svc-api.md`
- `rank-svc-api.md`

每个文件只写该服务负责的能力；跨服务调用在“实现备注”说明，不在接口层暴露内部 RPC 细节。
