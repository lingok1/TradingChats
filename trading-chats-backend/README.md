# Trading Chats Backend

Trading Chats 后端服务，提供提示词模板、模型 API 配置、AI 批量生成、定时任务、系统配置、登录鉴权与基础多租户隔离能力。

## 功能概览

- 用户登录、Token 刷新与会话管理
- 管理员与租户角色控制
- 多租户数据隔离（业务集合按 tenant_id 过滤）
- 提示词模板管理与动态 Prompt 生成
- 模型 API 配置管理与连接测试
- AI 批量生成与批次查询
- 定时任务配置、手动触发与执行日志
- 系统基础配置、动态参数、运行时参数管理
- Swagger 文档与健康检查接口

## 运行环境

- Go 1.23+
- MongoDB
- Redis

## 环境变量

| 变量名 | 默认值 | 说明 |
|---|---|---|
| PORT | 8080 | 服务端口 |
| MONGODB_URI | mongodb://localhost:27017 | MongoDB 连接地址 |
| MONGODB_DATABASE | trading_chats | MongoDB 数据库名 |
| REDIS_URI | redis://localhost:6379/0 | Redis 连接地址 |
| JWT_SECRET | change_me_in_production | JWT 密钥，生产环境必须覆盖 |
| JWT_EXPIRATION | 24h | Access Token 过期时间 |
| API_TIMEOUT | 60s | 外部模型 API 超时时间 |
| SEED_ADMIN_PASSWORD | Admin@123456 | 管理员种子密码，建议生产覆盖 |
| SEED_TENANT_ALPHA_PASSWORD | TenantAlpha@123456 | 租户 Alpha 种子密码 |
| SEED_TENANT_BETA_PASSWORD | TenantBeta@123456 | 租户 Beta 种子密码 |

## 本地启动

```bash
go mod tidy
go run cmd/api/main.go
```

服务默认地址：

- API: `http://localhost:8080`
- Swagger: `http://localhost:8080/swagger/index.html`
- Health: `http://localhost:8080/health`

## 默认账户

系统启动时会初始化 1 个管理员账号和 2 个租户账号：

- 管理员：`admin`
- 租户 Alpha：`tenant_alpha`
- 租户 Beta：`tenant_beta`

默认密码可通过环境变量覆盖，生产环境必须修改。

## 鉴权说明

### 公开接口

- `POST /api/auth/login`
- `POST /api/auth/refresh`
- `POST /api/auth/logout`
- 查询类 GET 接口
- `/health`
- `/swagger/*`

### 受保护接口

以下接口需要 `Authorization: Bearer <token>`：

- 所有创建类 POST 接口
- 所有修改类 PUT 接口
- 所有删除类 DELETE 接口
- Prompt 生成、模型连接测试、AI 生成、定时任务手动触发

## 多租户说明

当前业务集合已接入 `tenant_id` 过滤：

- prompt_templates
- model_api_configs
- ai_responses
- schedule_configs
- schedule_logs

规则：

- 管理员可跨租户访问
- 普通租户仅可访问本租户数据

注意：`system_configs` 当前仍为全局单例配置，不区分租户。

## 测试

当前建议执行：

```bash
go test ./...
```

如本机缺少对应 Go toolchain，请先安装匹配版本后再运行。
