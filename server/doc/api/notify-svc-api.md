# notify-svc API

> 负责站内通知查询与已读。采用拉取模型，无实时推送。

## 通知类型约定

| type | 含义 | data 示例 |
|---|---|---|
| like | 被点赞 | `{ "targetType": "post/comment", "targetId": 1, "fromUserId": 2 }` |
| comment | 被评论/被回复 | `{ "postId": 1, "commentId": 2, "fromUserId": 3 }` |
| system | 系统通知 | `{}` |

## 接口清单

- `GET /api/v1/notifications`
- `PUT /api/v1/notifications/{notificationId}/read`
- `PUT /api/v1/notifications/read-all`

---

## 获取我的通知列表

- **Method**: `GET`
- **Path**: `/api/v1/notifications`
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

### 响应

#### 成功示例
```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "list": [
      {
        "notificationId": 50001,
        "type": "like",
        "title": "收到点赞",
        "content": "你的帖子被点赞了",
        "data": {
          "targetType": "post",
          "targetId": 10001,
          "fromUserId": 10002
        },
        "isRead": false,
        "createdAt": "2026-04-21T10:30:00+08:00",
        "readAt": ""
      }
    ],
    "total": 1
  }
}
```

#### data 字段结构
| 字段 | 类型 | 说明 |
|---|---|---|
| list | array | 通知列表 |
| total | number | 总数 |

#### list item 字段结构
| 字段 | 类型 | 说明 |
|---|---|---|
| notificationId | number | 通知 ID |
| type | string | `like` / `comment` / `system` |
| title | string | 标题 |
| content | string | 内容摘要 |
| data | object | 扩展数据（可空） |
| isRead | boolean | 是否已读 |
| createdAt | string | 创建时间（RFC3339） |
| readAt | string | 已读时间（RFC3339；未读为空字符串） |

### 错误码
| code | 含义 | 触发条件 |
|---|---|---|
| 40101 | 未登录或 Token 无效 | 鉴权失败 |
| 80001 | 参数非法 | pageNum/pageSize/isRead 非法 |

### 业务规则（实现必须遵守）
- 排序固定 `createdAt desc`。
- 仅返回当前用户的通知。
- data 字段为 JSON 对象（由后端按通知 type 写入）。

### 实现备注（给后端 Agent）
- 所属服务：gateway + notify-svc
- 涉及实体：notifications
- 事务策略：读操作
- 幂等约束：无

---

## 标记单条通知已读

- **Method**: `PUT`
- **Path**: `/api/v1/notifications/{notificationId}/read`
- **Auth**: 需要
- **角色**: `user` / `admin`
- **说明**: 将指定通知标记为已读（幂等）

### 请求参数

#### Path
| 字段 | 类型 | 必填 | 默认值 | 说明 |
|---|---|---|---|---|
| notificationId | number | 是 | - | 通知 ID |

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
| 80002 | 通知不存在 | notificationId 不存在或不属于当前用户 |

### 业务规则（实现必须遵守）
- 幂等：重复标记同一条通知已读返回成功。
- 已读时间 `readAt` 仅在首次标记已读时写入。

### 实现备注（给后端 Agent）
- 所属服务：gateway + notify-svc
- 涉及实体：notifications
- 事务策略：单表写
- 幂等约束：`(notificationId, userId)` 约束（重复操作返回成功）

---

## 全部标记已读

- **Method**: `PUT`
- **Path**: `/api/v1/notifications/read-all`
- **Auth**: 需要
- **角色**: `user` / `admin`
- **说明**: 将我所有未读通知标记为已读（幂等）

### 请求参数

无

### 响应

#### 成功示例
```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "affected": 3
  }
}
```

#### data 字段结构
| 字段 | 类型 | 说明 |
|---|---|---|
| affected | number | 被更新为已读的数量 |

### 错误码
| code | 含义 | 触发条件 |
|---|---|---|
| 40101 | 未登录或 Token 无效 | 鉴权失败 |

### 业务规则（实现必须遵守）
- 幂等：重复调用返回成功，affected 可能为 0。
- 仅更新当前用户未读通知（`isRead=false`）。

### 实现备注（给后端 Agent）
- 所属服务：gateway + notify-svc
- 涉及实体：notifications
- 事务策略：单表批量写
- 幂等约束：条件更新（只更新未读）

---

## 建表/改表 SQL（MySQL 8）

### notifications（新增）

```sql
CREATE TABLE IF NOT EXISTS `notifications` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT UNSIGNED NOT NULL,
  `type` VARCHAR(16) NOT NULL,
  `title` VARCHAR(64) NOT NULL DEFAULT '',
  `content` VARCHAR(255) NOT NULL DEFAULT '',
  `data_json` VARCHAR(2000) NOT NULL DEFAULT '',
  `is_read` TINYINT(1) NOT NULL DEFAULT 0,
  `read_at` DATETIME NULL,
  `created_at` DATETIME NOT NULL,
  `updated_at` DATETIME NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_notifications_user_read_created_at` (`user_id`, `is_read`, `created_at`),
  KEY `idx_notifications_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
```

### 其他表结构变更

无。
