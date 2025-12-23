# 接口快速测试（curl）

服务默认端口参考 GoFrame 配置，示例使用 `http://localhost:8199`。返回结构统一为 `{code, message, data}`。

前端页面：`http://localhost:8199/html/index.html`

## 1) 注册
```bash
curl -X POST http://localhost:8199/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","password":"Passw0rd!"}'
```

## 2) 登录
```bash
curl -X POST http://localhost:8199/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","password":"Passw0rd!"}'
```
记录响应中的 `data.token`。

## 3) 创建/更新 Telegram 账号（按用户隔离，token 加密存储）

- 创建（设为默认账号）
```bash
TOKEN=替换为登录得到的token
curl -X POST http://localhost:8199/api/telegram/accounts \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name":"mybot",
    "bot_username":"your_bot",
    "bot_token":"123:AA...",
    "chat_id":-1001234567890,
    "is_default":true,
    "status":1
  }'
```

- 更新（可更新 token/默认标记）
```bash
curl -X PUT http://localhost:8199/api/telegram/accounts/{id} \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "bot_token":"123:NEW...",
    "is_default":true
  }'
```

- 设为默认
```bash
curl -X POST http://localhost:8199/api/telegram/accounts/{id}/default \
  -H "Authorization: Bearer $TOKEN"
```

- 列表（不返回明文 token）
```bash
curl -X GET http://localhost:8199/api/telegram/accounts \
  -H "Authorization: Bearer $TOKEN"
```

## 4) 上传文件（需要已有默认 Telegram 账号）
```bash
curl -X POST http://localhost:8199/api/files/upload \
  -H "Authorization: Bearer $TOKEN" \
  -F "file=@/path/to/your.file"
```
- 同一用户同 MD5 的文件会直接返回已有记录，不重复上传。
- 成功返回 `data.proxy_url`。

## 5) 文件列表
```bash
curl -X GET "http://localhost:8199/api/files?page=1&page_size=20" \
  -H "Authorization: Bearer $TOKEN"
```

## 6) 用户统计
```bash
curl -X GET http://localhost:8199/api/files/stats \
  -H "Authorization: Bearer $TOKEN"
```

## 7) 代理下载（支持匿名访问）
```bash
curl -L "http://localhost:8199/api/files/proxy/<文件md5>" -o output.bin
```
常见文本/图片/PDF/音视频使用 `inline`，其它类型为下载。

## 备注
- 上传前请确保当前用户有启用且默认的 Telegram 账号（`status=1`、`is_default=1`、`chat_id` 有效），且服务可访问 Telegram。
- Bot Token 仅在创建/更新时明文提交，服务端 AES-GCM 加密存储，不在列表返回。
- 上传大小默认限制 50MB，可在配置 `upload.maxBytes` 或 `internal/consts/upload.go` 调整。
- 可用环境变量覆盖数据库/密钥：`GF_GDB_DEFAULT_LINK`、`SECURITY_ENCRYPTKEY` 等。
