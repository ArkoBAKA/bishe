# Go 后端服务

## 1. 依赖环境

- Go：建议 1.22+
- MySQL：8.x（或兼容版本）
- Redis：6.x/7.x

## 2. 配置说明

项目默认读取仓库根目录下的 [config.yaml](file:///Users/wang/project/jz/bishe/server/config.yaml)。

- 默认配置文件路径：`./config.yaml`
- 可通过环境变量指定配置文件：`CONFIG_PATH=/path/to/config.yaml`

### 2.1 config.yaml 字段

```yaml
app:
  ginMode: debug        # gin 模式：debug / release / test
  autoMigrate: true     # 启动时自动迁移（示例表：users）
http:
  addr: ":8080"         # 服务监听地址
mysql:
  host: "127.0.0.1"
  port: 3306
  user: "root"
  password: "password"
  dbName: "app"
  params:
    charset: "utf8mb4"
    parseTime: "True"
    loc: "Local"
redis:
  host: "127.0.0.1"
  port: 6379
  password: ""
  db: 0
```

### 2.2 环境变量覆盖（可选）

启动时会在读取 `config.yaml` 后，再应用环境变量覆盖（方便线上注入）。

- `CONFIG_PATH`：配置文件路径（默认 `config.yaml`）
- `HTTP_ADDR`：覆盖 `http.addr`
- `GIN_MODE`：覆盖 `app.ginMode`
- `AUTO_MIGRATE`：覆盖 `app.autoMigrate`（true/false）
- `MYSQL_DSN`：覆盖 MySQL 连接（优先级高于 yaml 的 mysql 段）
- `REDIS_ADDR`：覆盖 Redis 地址（host:port）
- `REDIS_PASSWORD`：覆盖 Redis 密码
- `REDIS_DB`：覆盖 Redis DB（整数）

示例：

```bash
export HTTP_ADDR=":8080"
export GIN_MODE="debug"
export AUTO_MIGRATE="true"
export MYSQL_DSN='root:password@tcp(127.0.0.1:3306)/app?charset=utf8mb4&parseTime=True&loc=Local'
export REDIS_ADDR="127.0.0.1:6379"
export REDIS_PASSWORD=""
export REDIS_DB="0"
```

## 3. 启动方式

### 3.1 本地启动（推荐）

在 `server` 目录执行：

```bash
go mod tidy
go run .
```

默认会监听 `:8080`，可在 `config.yaml` 中修改 `http.addr`。

### 3.2 使用自定义配置文件启动

```bash
CONFIG_PATH=/absolute/path/to/config.yaml go run .
```

### 3.3 仅用环境变量启动（不改 yaml）

```bash
HTTP_ADDR=":8080" \
GIN_MODE="debug" \
AUTO_MIGRATE="true" \
MYSQL_DSN='root:password@tcp(127.0.0.1:3306)/app?charset=utf8mb4&parseTime=True&loc=Local' \
REDIS_ADDR="127.0.0.1:6379" \
REDIS_PASSWORD="" \
REDIS_DB="0" \
go run .
```

## 4. 健康检查与接口

### 4.1 健康检查

- URL：`GET /health`
- 含义：同时检查 MySQL/Redis 的连通性

```bash
curl -s http://127.0.0.1:8080/health | cat
```

可能返回：

```json
{"mysql":true,"ok":true,"redis":true}
```

如果 MySQL 或 Redis 不可用会返回 503，并带上错误信息字段（`mysqlError` / `redisError`）。

### 4.2 用户示例接口

创建用户：

```bash
curl -s -X POST http://127.0.0.1:8080/users \
  -H 'Content-Type: application/json' \
  -d '{"username":"alice"}' | cat
```

查询用户：

```bash
curl -s http://127.0.0.1:8080/users/1 | cat
```

## 5. 常见问题

### 5.1 启动时报 MySQL 连接失败

- 检查 MySQL 是否启动、端口是否正确
- 检查 `mysql.user/mysql.password/mysql.dbName` 或 `MYSQL_DSN`
- 确认数据库已创建（示例默认 `app`）

### 5.2 启动时报 Redis 连接失败

- 检查 Redis 是否启动、端口是否正确
- 检查 `redis.password/redis.db` 或环境变量覆盖项
