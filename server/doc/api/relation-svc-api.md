# relation-svc API

> 负责关注贴吧、关注用户、取消关注与我的关注列表。

## 接口清单

- `POST /api/v1/follows`
- `DELETE /api/v1/follows`
- `GET /api/v1/follows/me`

---

## 创建关注关系

- **Method**: `POST`
- **Path**: `/api/v1/follows`
- **Auth**: 需要
- **角色**: `user`
- **说明**: 创建关注关系

### 请求参数

#### Body
| 字段 | 类型 | 必填 | 默认值 | 说明 |
|---|---|---|---|---|
| targetType | string | 是 | - | `forum` / `user` |
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
| 70001 | 参数非法 | targetType/targetId 非法 |
| 70002 | 目标不存在 | forum/user 不存在 |

### 业务规则（实现必须遵守）
- 关系幂等，同一目标重复关注不重复创建。
- 若此前取消过关注，则本接口会恢复为关注中（状态置为 active）。
- 取消关注使用状态位，不物理删除。

### 实现备注（给后端 Agent）
- 所属服务：gateway + relation-svc
- 涉及实体：follows
- 事务策略：单表写
- 幂等约束：唯一键 `(user_id, target_type, target_id)`

---

## 取消关注

- **Method**: `DELETE`
- **Path**: `/api/v1/follows`
- **Auth**: 需要
- **角色**: `user`
- **说明**: 取消关注关系（幂等）

### 请求参数

#### Body
| 字段 | 类型 | 必填 | 默认值 | 说明 |
|---|---|---|---|---|
| targetType | string | 是 | - | `forum` / `user` |
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
| 70001 | 参数非法 | targetType/targetId 非法 |

### 业务规则（实现必须遵守）
- 幂等：未关注或重复取消，返回成功。
- 通过 `status=canceled` 表示取消关注，不物理删除。

---

## 我的关注列表

- **Method**: `GET`
- **Path**: `/api/v1/follows/me`
- **Auth**: 需要
- **角色**: `user`
- **说明**: 获取我的关注列表（默认仅 active，按 createdAt desc）

### 请求参数

#### Query
| 字段 | 类型 | 必填 | 默认值 | 说明 |
|---|---|---|---|---|
| pageNum | number | 否 | 1 | 页码 |
| pageSize | number | 否 | 10 | 每页数量 |
| targetType | string | 否 | - | `forum` / `user` |

### 响应

#### 成功示例
```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "list": [
      {
        "followId": 90001,
        "targetType": "forum",
        "targetId": 10001,
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
| list | array | 关注列表 |
| total | number | 总数 |

#### list item 字段结构
| 字段 | 类型 | 说明 |
|---|---|---|
| followId | number | 关注记录 ID |
| targetType | string | `forum` / `user` |
| targetId | number | 目标 ID |
| createdAt | string | 创建时间（RFC3339） |

### 错误码
| code | 含义 | 触发条件 |
|---|---|---|
| 40101 | 未登录或 Token 无效 | 鉴权失败 |
| 70003 | 参数非法 | pageNum/pageSize/targetType 非法 |

### 业务规则（实现必须遵守）
- 默认仅返回 `status=active`。
- 排序固定 `createdAt desc`。

---

## 建表/改表 SQL（MySQL 8）

### follows（新增）

```sql
CREATE TABLE IF NOT EXISTS `follows` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT UNSIGNED NOT NULL,
  `target_type` VARCHAR(16) NOT NULL,
  `target_id` BIGINT UNSIGNED NOT NULL,
  `status` VARCHAR(16) NOT NULL DEFAULT 'active',
  `canceled_at` DATETIME NULL,
  `created_at` DATETIME NOT NULL,
  `updated_at` DATETIME NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_follows_user_target` (`user_id`, `target_type`, `target_id`),
  KEY `idx_follows_user_type_status_created_at` (`user_id`, `target_type`, `status`, `created_at`),
  KEY `idx_follows_target_type_id` (`target_type`, `target_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
```

### 其他表结构变更

无。

### 数据迁移（可选）

若历史已使用 `forum_follows`，可把历史关注同步到 `follows`：

```sql
INSERT IGNORE INTO `follows` (`user_id`, `target_type`, `target_id`, `status`, `created_at`, `updated_at`)
SELECT `user_id`, 'forum', `forum_id`, 'active', `created_at`, `created_at`
FROM `forum_follows`;
```
