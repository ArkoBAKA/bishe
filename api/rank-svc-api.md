# rank-svc API

> 负责热榜读取接口与热榜计算任务约定（任务本身为内部能力，不直接对前端开放）。

## 接口清单

- `GET /api/v1/rank/posts`

---

## `GET /api/v1/rank/posts`

- **Auth**: 不需要
- **角色**: `all`
- **说明**: 获取帖子热榜

### 请求参数

#### Query
| 字段 | 类型 | 必填 | 默认值 | 说明 |
|---|---|---|---|---|
| pageNum | number | 否 | 1 | 页码 |
| pageSize | number | 否 | 10 | 每页数量 |
| forumId | number | 否 | - | 可选，按贴吧过滤 |

### 业务规则（实现必须遵守）
- 热榜仅基于帖子表必填字段计算：`likeCount`、`commentCount`、`viewCount`。
- 结果写入 Redis（如 ZSET），接口从 Redis 读取并回补展示信息。

---

## 内部任务说明（非对外 API）

- 定时任务周期、打分公式、Redis Key 规范在实现阶段补充到本文件。
# rank-svc API

> 负责热榜读取接口与热榜计算任务约定（任务本身为内部能力，不直接对前端开放）。

## 接口清单

- `GET /api/v1/rank/posts`

---

## `GET /api/v1/rank/posts`

- **Auth**: 不需要
- **角色**: `all`
- **说明**: 获取帖子热榜

### 请求参数

#### Query
| 字段 | 类型 | 必填 | 默认值 | 说明 |
|---|---|---|---|---|
| pageNum | number | 否 | 1 | 页码 |
| pageSize | number | 否 | 10 | 每页数量 |
| forumId | number | 否 | - | 可选，按贴吧过滤 |

### 业务规则（实现必须遵守）
- 热榜仅基于帖子表必填字段计算：`likeCount`、`commentCount`、`viewCount`。
- 结果写入 Redis（如 ZSET），接口从 Redis 读取并回补展示信息。

---

## 内部任务说明（非对外 API）

- 定时任务周期、打分公式、Redis Key 规范在实现阶段补充到本文件。
