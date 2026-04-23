# user-svc API

> 负责用户注册、登录、登出、资料与密码管理。对外路径统一经网关暴露为 `/api/v1/...`。

## 接口清单

- `POST /api/v1/users/register`
- `POST /api/v1/users/login`
- `POST /api/v1/users/logout`
- `GET /api/v1/users/me`
- `GET /api/v1/users/{userId}`
- `PUT /api/v1/users/me/profile`
- `PUT /api/v1/users/me/password`

***

## `POST /api/v1/users/register`

- **Method**: `POST`
- **Path**: `/api/v1/users/register`
- **Auth**: 不需要
- **角色**: `guest`
- **说明**: 账号密码注册，支持头像文件上传（后端上传 本地路径 并落库 `avatarUrl`）

### 请求参数

#### Body

| 字段         | 类型     | 必填 | 默认值 | 说明                            |
| ---------- | ------ | -- | --- | ----------------------------- |
| account    | string | 是  | -   | 登录账号（唯一）                      |
| password   | string | 是  | -   | 明文密码（服务端哈希）                   |
| nickname   | string | 否  | -   | 昵称，不传则由服务端生成默认昵称              |
| avatarFile | file   | 否  | -   | 头像文件，`multipart/form-data` 上传 |

### 响应

#### 成功示例

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "userId": 10001,
    "avatarUrl": "https://本地路径.example.com/tieba/avatar/10001.png"
  }
}
```

#### data 字段结构

| 字段        | 类型     | 说明                |
| --------- | ------ | ----------------- |
| userId    | number | 新注册用户 ID          |
| avatarUrl | string | 头像 URL（未上传则为空字符串） |

### 错误码

| code  | 含义     | 触发条件             |
| ----- | ------ | ---------------- |
| 10001 | 参数非法   | 账号或密码格式不合法       |
| 10002 | 账号已存在  | account 重复       |
| 10003 | 密码强度不足 | 未满足最小长度或复杂度      |
| 10004 | 头像上传失败 | 上传 本地路径 失败或文件格式非法 |

### 业务规则（实现必须遵守）

- 账号唯一。
- 密码仅保存哈希。
- account 入库前需标准化（去首尾空格）。
- 注册接口 `Content-Type` 使用 `multipart/form-data`。
- 有 `avatarFile` 时，后端先上传 本地路径，得到 `avatarUrl` 后再写入用户表。

### 实现备注（给后端 Agent）

- 所属服务：gateway + user-svc
- 涉及实体：用户
- 事务策略：本地路径 上传 + 单表写入事务（写库失败时可回收已上传对象，失败可仅记录日志）
- 幂等约束：账号唯一索引

***

## `POST /api/v1/users/login`

- **Method**: `POST`
- **Path**: `/api/v1/users/login`
- **Auth**: 不需要
- **角色**: `guest`
- **说明**: 账号密码登录并签发 JWT

### 请求参数

#### Body

| 字段       | 类型     | 必填 | 默认值 | 说明   |
| -------- | ------ | -- | --- | ---- |
| account  | string | 是  | -   | 登录账号 |
| password | string | 是  | -   | 登录密码 |

### 响应

#### 成功示例

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "token": "jwt-token",
    "tokenType": "Bearer",
    "expiresIn": 7200,
    "user": {
      "userId": 10001,
      "account": "arkovo",
      "nickname": "Ark",
      "avatarUrl": "",
      "bio": "",
      "role": "user",
      "status": "normal"
    }
  }
}
```

#### data 字段结构

| 字段        | 类型     | 说明          |
| --------- | ------ | ----------- |
| token     | string | JWT         |
| tokenType | string | 固定 `Bearer` |
| expiresIn | number | 秒           |
| user      | object | 当前用户摘要      |

### 错误码

| code  | 含义      | 触发条件     |
| ----- | ------- | -------- |
| 10011 | 账号或密码错误 | 认证失败     |
| 10012 | 账号已封禁   | 用户状态不可登录 |

### 业务规则（实现必须遵守）

- 登录成功返回 JWT，JWT 中至少包含 `userId`、`role`、`exp`。
- 账号封禁状态下不允许登录。

### 实现备注（给后端 Agent）

- 所属服务：gateway + user-svc
- 涉及实体：用户
- 事务策略：读操作 + token 签发
- 幂等约束：无

***

## `POST /api/v1/users/logout`

- **Method**: `POST`
- **Path**: `/api/v1/users/logout`
- **Auth**: 需要
- **角色**: `user` / `admin`
- **说明**: 登出当前会话

### 请求参数

无请求体。

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

| code  | 含义            | 触发条件 |
| ----- | ------------- | ---- |
| 40101 | 未登录或 Token 无效 | 鉴权失败 |

### 业务规则（实现必须遵守）

- 项目不做 token 黑名单，登出为客户端删除 token。

### 实现备注（给后端 Agent）

- 所属服务：gateway + user-svc
- 涉及实体：无
- 事务策略：无
- 幂等约束：天然幂等

***

## `GET /api/v1/users/me`

- **Method**: `GET`
- **Path**: `/api/v1/users/me`
- **Auth**: 需要
- **角色**: `user` / `admin`
- **说明**: 查询当前登录用户信息

### 请求参数

无。

### 响应

#### 成功示例

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "userId": 10001,
    "account": "arkovo",
    "nickname": "Ark",
    "avatarUrl": "",
    "bio": "",
    "role": "user",
    "status": "normal",
    "createdAt": "2026-04-21T10:30:00+08:00"
  }
}
```

#### data 字段结构

| 字段        | 类型     | 说明                  |
| --------- | ------ | ------------------- |
| userId    | number | 用户 ID               |
| account   | string | 登录账号                |
| nickname  | string | 昵称                  |
| avatarUrl | string | 头像                  |
| bio       | string | 简介                  |
| role      | string | `user` / `admin`    |
| status    | string | `normal` / `banned` |
| createdAt | string | 注册时间（RFC3339）       |

### 错误码

| code  | 含义            | 触发条件         |
| ----- | ------------- | ------------ |
| 40101 | 未登录或 Token 无效 | 鉴权失败         |
| 10021 | 用户不存在         | token 中用户不存在 |

### 业务规则（实现必须遵守）

- 仅返回当前 token 对应用户，不支持查他人详情。

### 实现备注（给后端 Agent）

- 所属服务：gateway + user-svc
- 涉及实体：用户
- 事务策略：读操作
- 幂等约束：无

***

## `GET /api/v1/users/{userId}`

- **Method**: `GET`
- **Path**: `/api/v1/users/{userId}`
- **Auth**: 不需要
- **角色**: `all`
- **说明**: 获取指定用户公开信息（用于帖子/评论区头像、昵称展示）

### 请求参数

#### Path

| 字段     | 类型     | 必填 | 默认值 | 说明      |
| ------ | ------ | -- | --- | ------- |
| userId | number | 是  | -   | 目标用户 ID |

### 响应

#### 成功示例

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "userId": 10002,
    "nickname": "Alice",
    "avatarUrl": "https://example.com/avatar.png",
    "bio": "喜欢分享技术",
    "role": "user"
  }
}
```

#### data 字段结构

| 字段        | 类型     | 说明                       |
| --------- | ------ | ------------------------ |
| userId    | number | 用户 ID                    |
| nickname  | string | 昵称                       |
| avatarUrl | string | 头像 URL                   |
| bio       | string | 个人简介                     |
| role      | string | `user` / `admin`（公开展示可选） |

### 错误码

| code  | 含义    | 触发条件              |
| ----- | ----- | ----------------- |
| 10051 | 用户不存在 | userId 无对应记录      |
| 10052 | 用户不可见 | 用户被封禁或账号状态不允许公开展示 |

### 业务规则（实现必须遵守）

- 仅返回公开字段，不返回 `account`、`status`、`passwordHash` 等敏感信息。
- 被封禁用户是否可公开展示，由业务配置决定；默认可返回基础资料。

### 实现备注（给后端 Agent）

- 所属服务：gateway + user-svc
- 涉及实体：用户
- 事务策略：读操作
- 幂等约束：无

***

## `PUT /api/v1/users/me/profile`

- **Method**: `PUT`
- **Path**: `/api/v1/users/me/profile`
- **Auth**: 需要
- **角色**: `user` / `admin`
- **说明**: 修改当前用户资料，支持头像文件上传（后端上传 本地路径 并回写 `avatarUrl`）

### 请求参数

#### Body

| 字段         | 类型     | 必填 | 默认值 | 说明                             |
| ---------- | ------ | -- | --- | ------------------------------ |
| nickname   | string | 否  | -   | 昵称                             |
| avatarFile | file   | 否  | -   | 新头像文件，`multipart/form-data` 上传 |
| bio        | string | 否  | -   | 简介                             |

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

| code  | 含义            | 触发条件             |
| ----- | ------------- | ---------------- |
| 40101 | 未登录或 Token 无效 | 鉴权失败             |
| 10031 | 参数非法          | 长度或格式不合法         |
| 10032 | 头像上传失败        | 上传 本地路径 失败或文件格式非法 |

### 业务规则（实现必须遵守）

- 至少一个字段有变更才执行更新。
- 接口 `Content-Type` 使用 `multipart/form-data`。
- 有 `avatarFile` 时，后端上传 本地路径 成功后回写 `avatarUrl`。

### 实现备注（给后端 Agent）

- 所属服务：gateway + user-svc
- 涉及实体：用户
- 事务策略：本地路径 上传 + 单表更新事务（更新失败时可回收已上传对象，失败可仅记录日志）
- 幂等约束：重复提交同样内容不报错

***

## `PUT /api/v1/users/me/password`

- **Method**: `PUT`
- **Path**: `/api/v1/users/me/password`
- **Auth**: 需要
- **角色**: `user` / `admin`
- **说明**: 修改当前用户密码

### 请求参数

#### Body

| 字段          | 类型     | 必填 | 默认值 | 说明  |
| ----------- | ------ | -- | --- | --- |
| oldPassword | string | 是  | -   | 旧密码 |
| newPassword | string | 是  | -   | 新密码 |

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

| code  | 含义            | 触发条件             |
| ----- | ------------- | ---------------- |
| 40101 | 未登录或 Token 无效 | 鉴权失败             |
| 10041 | 旧密码错误         | oldPassword 校验失败 |
| 10042 | 新密码不合法        | 强度或长度不满足         |

### 业务规则（实现必须遵守）

- 新旧密码不可相同。
- 修改成功后，当前会话是否继续有效由网关策略控制（本项目可保持有效，简化实现）。

### 实现备注（给后端 Agent）

- 所属服务：gateway + user-svc
- 涉及实体：用户
- 事务策略：单表更新事务
- 幂等约束：无

