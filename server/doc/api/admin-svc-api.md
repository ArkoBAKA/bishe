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
{
  "code": 0,
  "message": "ok",
  "data": {
    "list": [
      {
        "postId": 10001,
        "forumId": 10001,
        "title": "Hello",
        "authorId": 10002,
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
| 40101 | 未登录或 Token 无效 | 鉴权失败 |
| 90001 | 需要管理员权限 | 非 admin 访问 |
| 20001 | 参数非法 | pageNum/pageSize 非法 |

### 业务规则（实现必须遵守）
- 仅返回 `pending`（待审核）状态的帖子。
- 按 `createdAt asc` 排序（先审先处理）。

---

## 审核帖子并更新状态

- **Method**: `PUT`
- **Path**: `/api/v1/admin/posts/{postId}/review`
- **Auth**: 需要
- **角色**: `admin`
- **说明**: 审核帖子并更新可见状态

### 请求参数

#### Path
| 字段 | 类型 | 必填 | 默认值 | 说明 |
|---|---|---|---|---|
| postId | number | 是 | - | 帖子 ID |

#### Body
| 字段 | 类型 | 必填 | 默认值 | 说明 |
|---|---|---|---|---|
| action | string | 是 | - | `approve` / `reject` / `hide` |
| reviewRemark | string | 否 | - | 处理备注 |

### 响应

#### 成功示例
```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "status": "visible"
  }
}
```

#### data 字段结构
| 字段 | 类型 | 说明 |
|---|---|---|
| status | string | 更新后的帖子状态 |

### 错误码
| code | 含义 | 触发条件 |
|---|---|---|
| 40101 | 未登录或 Token 无效 | 鉴权失败 |
| 90001 | 需要管理员权限 | 非 admin 操作 |
| 90002 | 参数非法 | action 非法 / reviewRemark 过长 |
| 90003 | 帖子不存在 | postId 无对应记录 |

### 业务规则（实现必须遵守）
- 审核动作必须映射到帖子可见状态流转。
- `approve` -> `visible`
- `reject` -> `rejected`
- `hide` -> `hidden`
- 需记录 `reviewedBy`、`reviewedAt`、`reviewRemark`。

---

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
| targetType | string | 是 | - | `post` / `comment` / `user` |
| targetId | number | 是 | - | 目标 ID |
| reason | string | 是 | - | 简短原因（1~64） |
| detail | string | 否 | - | 详细描述（<=1000） |

### 响应

#### 成功示例
```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "reportId": 30001
  }
}
```

#### data 字段结构
| 字段 | 类型 | 说明 |
|---|---|---|
| reportId | number | 举报单 ID |

### 错误码
| code | 含义 | 触发条件 |
|---|---|---|
| 40101 | 未登录或 Token 无效 | 鉴权失败 |
| 91001 | 参数非法 | targetType/targetId/reason/detail 不合法 |

### 业务规则（实现必须遵守）
- 同一用户对同一目标的重复举报允许存在（本项目不做去重）。

---

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
| status | string | 否 | - | `pending` / `processed` |

### 响应

#### 成功示例
```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "list": [
      {
        "reportId": 30001,
        "reporterId": 10001,
        "targetType": "post",
        "targetId": 10001,
        "reason": "广告",
        "detail": "",
        "status": "pending",
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
| 40101 | 未登录或 Token 无效 | 鉴权失败 |
| 90001 | 需要管理员权限 | 非 admin 访问 |
| 20001 | 参数非法 | pageNum/pageSize/status 非法 |

---

## 处理举报并执行动作

- **Method**: `PUT`
- **Path**: `/api/v1/admin/reports/{reportId}/process`
- **Auth**: 需要
- **角色**: `admin`
- **说明**: 处理举报审计单并执行动作

### 请求参数

#### Path
| 字段 | 类型 | 必填 | 默认值 | 说明 |
|---|---|---|---|---|
| reportId | number | 是 | - | 举报单 ID |

#### Body
| 字段 | 类型 | 必填 | 默认值 | 说明 |
|---|---|---|---|---|
| action | string | 是 | - | `close` / `deletePost` / `deleteComment` / `hidePost` / `banUser` |
| processRemark | string | 否 | - | 处理备注 |
| banUntil | string | 否 | - | 当 action=banUser 时可填（RFC3339） |
| durationSeconds | number | 否 | - | 当 action=banUser 时可填（秒）；与 banUntil 二选一 |

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
| 90001 | 需要管理员权限 | 非 admin 操作 |
| 92001 | 参数非法 | action/processRemark/banUntil/durationSeconds 不合法 |
| 92002 | 举报单不存在 | reportId 无对应记录 |
| 92003 | 状态冲突 | 举报单已处理 |
| 92004 | 目标不匹配 | action 与 report.targetType 不匹配 |

### 业务规则（实现必须遵守）
- 报告处理需记录处理人、处理时间、处理动作、处理备注。
- 默认只允许处理一次：`pending` -> `processed`。
- action 执行动作需与举报目标类型一致：
  - `deletePost`/`hidePost` 仅适用于 `post`
  - `deleteComment` 仅适用于 `comment`
  - `banUser` 仅适用于 `user`

---

## 管理员删除帖子

- **Method**: `DELETE`
- **Path**: `/api/v1/admin/posts/{postId}`
- **Auth**: 需要
- **角色**: `admin`
- **说明**: 删除帖子（软删）

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
| 90001 | 需要管理员权限 | 非 admin 操作 |
| 90003 | 帖子不存在 | postId 无对应记录 |

---

## 管理员删除评论

- **Method**: `DELETE`
- **Path**: `/api/v1/admin/comments/{commentId}`
- **Auth**: 需要
- **角色**: `admin`
- **说明**: 删除评论（软删）

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
| 90001 | 需要管理员权限 | 非 admin 操作 |
| 90013 | 评论不存在 | commentId 无对应记录 |

---

## 封禁用户

- **Method**: `PUT`
- **Path**: `/api/v1/admin/users/{userId}/ban`
- **Auth**: 需要
- **角色**: `admin`
- **说明**: 封禁用户（设置 status=banned 与 banUntil）

### 请求参数

#### Path
| 字段 | 类型 | 必填 | 默认值 | 说明 |
|---|---|---|---|---|
| userId | number | 是 | - | 用户 ID |

#### Body
| 字段 | 类型 | 必填 | 默认值 | 说明 |
|---|---|---|---|---|
| banUntil | string | 否 | - | 封禁截止时间（RFC3339） |
| durationSeconds | number | 否 | - | 封禁时长（秒）；与 banUntil 二选一 |
| remark | string | 否 | - | 备注（<=255） |

### 响应

#### 成功示例
```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "status": "banned",
    "banUntil": "2026-04-30T10:30:00+08:00"
  }
}
```

### 错误码
| code | 含义 | 触发条件 |
|---|---|---|
| 40101 | 未登录或 Token 无效 | 鉴权失败 |
| 90001 | 需要管理员权限 | 非 admin 操作 |
| 93001 | 参数非法 | banUntil/durationSeconds/remark 不合法 |
| 93002 | 用户不存在 | userId 无对应记录 |

### 业务规则（实现必须遵守）
- 若未提供 `banUntil` 与 `durationSeconds`，默认封禁 7 天。

---

## 建表/改表 SQL（MySQL 8）

### reports（新增）

```sql
CREATE TABLE IF NOT EXISTS `reports` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `reporter_id` BIGINT UNSIGNED NOT NULL,
  `target_type` VARCHAR(16) NOT NULL,
  `target_id` BIGINT UNSIGNED NOT NULL,
  `reason` VARCHAR(64) NOT NULL,
  `detail` VARCHAR(1000) NOT NULL DEFAULT '',
  `status` VARCHAR(16) NOT NULL DEFAULT 'pending',
  `processed_by` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `processed_at` DATETIME NULL,
  `process_action` VARCHAR(32) NOT NULL DEFAULT '',
  `process_remark` VARCHAR(255) NOT NULL DEFAULT '',
  `created_at` DATETIME NOT NULL,
  `updated_at` DATETIME NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_reports_status_created_at` (`status`, `created_at`),
  KEY `idx_reports_target` (`target_type`, `target_id`),
  KEY `idx_reports_reporter_id_created_at` (`reporter_id`, `created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
```

### posts（增加审核字段，如已有可忽略）

```sql
ALTER TABLE `posts`
  ADD COLUMN `reviewed_by` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  ADD COLUMN `reviewed_at` DATETIME NULL,
  ADD COLUMN `review_remark` VARCHAR(255) NOT NULL DEFAULT '';
```

### users（确保存在 ban_until 字段，如已有可忽略）

```sql
ALTER TABLE `users`
  ADD COLUMN `ban_until` DATETIME NULL;
```

其余接口按 `api/README.md` 的固定模板补充详细字段、示例与错误码。
