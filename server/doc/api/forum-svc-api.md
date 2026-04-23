# forum-svc API

> 负责贴吧创建、编辑、列表、详情。

## 接口清单

- `POST /api/v1/forums`
- `PUT /api/v1/forums/{forumId}`
- `GET /api/v1/forums`
- `GET /api/v1/forums/{forumId}`

---

## 创建贴吧

- **Method**: `POST`
- **Path**: `/api/v1/forums`
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

### 响应

#### 成功示例
```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "forumId": 10001
  }
}
```

#### data 字段结构
| 字段 | 类型 | 说明 |
|---|---|---|
| forumId | number | 新创建贴吧 ID |

### 错误码
| code | 含义 | 触发条件 |
|---|---|---|
| 40101 | 未登录或 Token 无效 | 鉴权失败 |
| 20001 | 参数非法 | name 为空/过长，description/coverUrl 格式不合法 |
| 20002 | 贴吧已存在 | name 重复（全局唯一冲突） |

### 业务规则（实现必须遵守）
- `name` 全局唯一。
- `name` 入库前需标准化（去首尾空格）。
- 默认创建者为贴吧 owner（仅 owner 可编辑）。

### 实现备注（给后端 Agent）
- 所属服务：gateway + forum-svc
- 涉及实体：forums
- 事务策略：单表写入事务
- 幂等约束：`uk_forums_name(name)`

---

## 编辑贴吧

- **Method**: `PUT`
- **Path**: `/api/v1/forums/{forumId}`
- **Auth**: 需要
- **角色**: `user`
- **说明**: 编辑指定贴吧（仅创建者/owner）

### 请求参数

#### Path
| 字段 | 类型 | 必填 | 默认值 | 说明 |
|---|---|---|---|---|
| forumId | number | 是 | - | 贴吧 ID |

#### Body
| 字段 | 类型 | 必填 | 默认值 | 说明 |
|---|---|---|---|---|
| name | string | 否 | - | 贴吧名称（修改时仍需全局唯一） |
| description | string | 否 | - | 贴吧简介 |
| coverUrl | string | 否 | - | 封面图地址 |

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
| 20001 | 参数非法 | 字段长度/格式不合法，或无任何字段变更 |
| 20002 | 贴吧已存在 | 修改 name 时发生唯一冲突 |
| 20003 | 贴吧不存在 | forumId 无对应记录 |
| 20004 | 无权限 | 非 owner 试图编辑 |

### 业务规则（实现必须遵守）
- 至少一个字段有变更才执行更新。
- `name` 修改时同样需要标准化（去首尾空格）并保证全局唯一。
- 仅创建者/owner 可编辑。

### 实现备注（给后端 Agent）
- 所属服务：gateway + forum-svc
- 涉及实体：forums
- 事务策略：单表更新事务
- 幂等约束：重复提交同样内容不报错（返回成功）

---

## 贴吧列表

- **Method**: `GET`
- **Path**: `/api/v1/forums`
- **Auth**: 不需要
- **角色**: `guest`
- **说明**: 分页查询贴吧列表（按创建时间倒序）

### 请求参数

#### Query
| 字段 | 类型 | 必填 | 默认值 | 说明 |
|---|---|---|---|---|
| pageNum | number | 否 | 1 | 页码，从 1 开始 |
| pageSize | number | 否 | 10 | 每页条数 |
| keyword | string | 否 | - | 名称关键字（模糊匹配） |

### 响应

#### 成功示例
```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "list": [
      {
        "forumId": 10001,
        "name": "Golang",
        "description": "Go 语言讨论",
        "coverUrl": "",
        "ownerId": 10001,
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
| list | array | 贴吧列表 |
| total | number | 总数 |

#### list item 字段结构
| 字段 | 类型 | 说明 |
|---|---|---|
| forumId | number | 贴吧 ID |
| name | string | 贴吧名称 |
| description | string | 贴吧简介 |
| coverUrl | string | 封面图 |
| ownerId | number | 创建者/owner 用户 ID |
| createdAt | string | 创建时间（RFC3339） |

### 错误码
| code | 含义 | 触发条件 |
|---|---|---|
| 20001 | 参数非法 | pageNum/pageSize 非法 |

### 业务规则（实现必须遵守）
- 按 `createdAt desc` 排序。
- `keyword` 走简单模糊匹配（SQL `LIKE`）。

### 实现备注（给后端 Agent）
- 所属服务：gateway + forum-svc
- 涉及实体：forums
- 事务策略：读操作
- 幂等约束：无

---

## 贴吧详情

- **Method**: `GET`
- **Path**: `/api/v1/forums/{forumId}`
- **Auth**: 不需要
- **角色**: `guest`
- **说明**: 获取贴吧详情

### 请求参数

#### Path
| 字段 | 类型 | 必填 | 默认值 | 说明 |
|---|---|---|---|---|
| forumId | number | 是 | - | 贴吧 ID |

### 响应

#### 成功示例
```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "forumId": 10001,
    "name": "Golang",
    "description": "Go 语言讨论",
    "coverUrl": "",
    "ownerId": 10001,
    "createdAt": "2026-04-21T10:30:00+08:00"
  }
}
```

#### data 字段结构
| 字段 | 类型 | 说明 |
|---|---|---|
| forumId | number | 贴吧 ID |
| name | string | 贴吧名称 |
| description | string | 贴吧简介 |
| coverUrl | string | 封面图 |
| ownerId | number | 创建者/owner 用户 ID |
| createdAt | string | 创建时间（RFC3339） |

### 错误码
| code | 含义 | 触发条件 |
|---|---|---|
| 20003 | 贴吧不存在 | forumId 无对应记录 |

### 业务规则（实现必须遵守）
- 仅返回公开字段。

### 实现备注（给后端 Agent）
- 所属服务：gateway + forum-svc
- 涉及实体：forums
- 事务策略：读操作
- 幂等约束：无

---

## 建表语句（MySQL 8）

```sql
CREATE TABLE IF NOT EXISTS `forums` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(64) NOT NULL,
  `description` VARCHAR(255) NOT NULL DEFAULT '',
  `cover_url` VARCHAR(512) NOT NULL DEFAULT '',
  `owner_id` BIGINT UNSIGNED NOT NULL,
  `created_at` DATETIME NOT NULL,
  `updated_at` DATETIME NOT NULL,
  `deleted_at` DATETIME NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_forums_name` (`name`),
  KEY `idx_forums_owner_id` (`owner_id`),
  KEY `idx_forums_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
```
