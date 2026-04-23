# upload-svc API

> 负责上传文件，静态目录映射 `/api/v1/...`。
## 核心职责

1. **统一文件上传**  
   提供 `POST /api/v1/upload` 接口，支持 `multipart/form-data` 单文件上传，返回可直接访问的 URL。  
   上传完成后自动落库（文件元数据表），用于后续审计、去重、生命周期管理。

2. **静态目录映射**  
   所有上传成功的文件统一挂载到  
   `/api/v1/static/{bucket}/{year}/{month}/{day}/{uuid}.{ext}`  
   网关直接透传，upload-svc 仅做「只读」映射，不经过业务逻辑，保障大文件下载性能。  不需要token即可查看文件
   支持 HTTP 缓存头（`ETag`/`Last-Modified`/`Cache-Control: max-age=31536000, immutable`）。

3. **安全策略**  
   - 仅允许登录后的用户上传（`Authorization: Bearer`）。  
   - 文件白名单：`jpg|jpeg|png|gif|webp|mp4|mov|pdf|zip|txt|md`。  
   - 单文件 ≤ 50 MB；单次请求 ≤ 100 MB。  
   - 文件名使用 `UUIDv4 + 原扩展名`，屏蔽用户原始路径信息。  
   - 桶（bucket）隔离：公开文件放 `public`，敏感文件放 `private`（后续可接私有签名下载）。

4. **高可用与可观测**  
   - 上传接口先写「临时文件」，校验完成后再原子移动到正式目录，失败可自动清理。  
   - 记录 `upload_log`：用户 ID、文件大小、MIME、耗时、客户端 IP。  
   - Prometheus 指标：`upload_total`、`upload_duration_seconds`、`upload_failures_total`。

## 接口清单（upload-svc 自身）

- `POST /api/v1/upload`                单文件上传（返回 URL）
- `GET  /api/v1/static/{bucket}/**`    静态映射（只读）
- `GET  /api/v1/upload/{fileId}`       查询文件元数据
- `DELETE /api/v1/upload/{fileId}`     删除（软删，仅本人或管理员）

## 上传接口示例

### `POST /api/v1/upload`

- **Content-Type**: `multipart/form-data; boundary=...`  
- **请求头**: `Authorization: Bearer <token>`  
- **表单字段**:
  | 字段       | 类型   | 必填 | 说明                     |
  | ---------- | ------ | ---- | ------------------------ |
  | file       | file   | 是   | 文件本身                 |
  | bucket     | string | 否   | `public`(默认)/`private` |
  | scene      | string | 否   | 业务场景，如 `avatar`/`post`/`comment`，仅做标记 |

- **成功 200**

