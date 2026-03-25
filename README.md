# 凌期AI辅助期货挑选品种与持仓止盈止损分析

## 项目简介

这是一个面向期货交易分析场景的前后端分离项目，通过提示词模板、动态参数、模型配置与批次任务调度，调用多个 AI 模型生成期货品种筛选与持仓止盈止损分析结果。

项目目标包括：

- 按批次生成多模型 AI 分析结果
- 以 Markdown 表格与结构化信号形式展示分析报告
- 支持提示词模板、模型配置、定时任务与系统配置管理
- 支持 PC 与移动端差异化展示
- 通过 Swagger 驱动前后端接口对接

## 仓库结构

```text
TradingChats-doc/
├── trading-chats-backend/         # Go + Gin 后端服务
├── trading-chats-frontend/        # Vue 3 + Vite 前端应用
├── futures-functional-demand-backend.md
├── futures-functional-demand-frontend.md
└── README.md
```

## 技术架构

### 后端

- Golang
- Gin
- MongoDB
- Redis
- Swagger

### 前端

- Vue 3
- Vite
- TypeScript
- Element Plus
- Axios
- markdown-it

## 核心业务能力

### 1. AI 分析报告

- 默认获取最新一次成功批次：`GET /api/ai-responses/latest`
- 支持按批次查询历史结果：`GET /api/ai-responses/batch?batch_id=...`
- 支持手动触发批次生成：`POST /api/ai-responses/generate`
- 多模型结果按模型维度展示与对比
- PC 端优先展示 Markdown 表格
- 移动端优先展示结构化信号卡片与详情抽屉

### 2. 提示词模板管理

- 支持模板增删改查
- 支持标签分类，例如股票、期货、期权等
- 支持基于模板与参数动态生成提示词
- 支持将外部 HTTPS JSON 数据拼接进提示词

### 3. 模型与 API 配置管理

- 支持配置自定义 API URL、API Key、模型列表
- 支持 OpenAI 与 Anthropic 协议
- 支持模型连通性测试

### 4. 定时任务调度

- 支持交易时段内定时执行 AI 分析任务
- 支持任务状态管理与执行日志查看
- 支持按批次存储与追踪分析结果

### 5. 系统配置管理

- 支持系统标题与 Logo 配置
- 支持动态参数配置
- 前后端已拆分为独立接口：
  - `GET /api/system-config`
  - `PUT /api/system-config/basic`
  - `PUT /api/system-config/parameters`

## 前端展示约定

### PC 端

- 首页按模型卡片纵向排列
- 卡片内渲染 AI 返回的 Markdown 表格
- 支持多模型上下滚动对比

### 移动端

- 先展示信号层列表
- 点击单条信号后，通过抽屉或全屏弹层查看详情
- 若结构化信号不存在，则降级展示 Markdown 内容

### 未完成页面策略

- 首页展示已开发完成内容
- 当日分析、多品种、持仓、新闻等未开发页面显示空白提示页

## 数据契约与接口约定

### Swagger 驱动开发

前端接口对接必须以 Swagger 为准：

- `trading-chats-backend/docs/swagger.yaml`
- `trading-chats-backend/docs/swagger.json`

统一约定：

- `basePath: /api`
- 统一响应结构：`models.Response`
- 若 Swagger 与实现不一致，应先修正 Swagger 再联调前端

### AI 响应结构

当前前端兼容两类后端返回：

1. 现状：直接返回 `AIResponse[]`
2. 目标：返回包含 `batch_id`、`batch_status`、`models[]` 的聚合结构

目标结构中每个模型建议包含：

- `response_markdown`
- `signals`
- `status`
- `error`
- `model_name`
- `provider`

## 主要接口

### AI 响应

- `GET /api/ai-responses/latest`
- `GET /api/ai-responses/batch?batch_id=...`
- `POST /api/ai-responses/generate`
- `GET /api/ai-responses`
- `GET /api/ai-responses/{id}`

### 提示词模板

- `POST /api/prompt-templates`
- `GET /api/prompt-templates`
- `GET /api/prompt-templates/tag?tag=...`
- `GET /api/prompt-templates/{id}`
- `PUT /api/prompt-templates/{id}`
- `DELETE /api/prompt-templates/{id}`
- `POST /api/prompt-templates/generate`

### 模型配置

- `POST /api/model-api-configs`
- `GET /api/model-api-configs`
- `GET /api/model-api-configs/provider?provider=...`
- `GET /api/model-api-configs/{id}`
- `PUT /api/model-api-configs/{id}`
- `DELETE /api/model-api-configs/{id}`
- `POST /api/model-api-configs/{id}/test`

### 定时任务

- `POST /api/schedules`
- `GET /api/schedules`
- `PUT /api/schedules/{id}/status`
- `GET /api/schedules/{id}/logs`
- `DELETE /api/schedules/{id}`

### 系统配置

- `GET /api/system-config`
- `PUT /api/system-config/basic`
- `PUT /api/system-config/parameters`

## 本地开发

## 环境要求

### 后端

- Go 1.20+
- MongoDB
- Redis

### 前端

- Node.js
- npm

## 后端启动

目录：`trading-chats-backend`

### 安装依赖

```bash
go mod download
```

### 配置环境变量

参考后端目录下的 `.env` 文件，主要包括：

- `PORT`
- `MONGODB_URI`
- `MONGODB_DATABASE`
- `REDIS_URI`
- `JWT_SECRET`
- `API_TIMEOUT`

### 开发启动

```bash
go run cmd/api/main.go
```

### Swagger 文档地址

```text
http://localhost:8080/swagger/index.html
```

### 重新生成 Swagger 文档

```bash
swag init -g cmd/api/main.go
```

## 前端启动

目录：`trading-chats-frontend`

### 安装依赖

```bash
npm install
```

### 开发启动

```bash
npm run dev
```

默认地址：

```text
http://localhost:5173/
```

### 构建

```bash
npm run build
```

## 主题与交互规范

- 支持明暗主题切换
- 主题选择需要持久化
- 默认优先使用用户配置，其次跟随系统主题
- 做多使用红色系，做空使用绿色系
- 移动端设置入口采用右侧抽屉
- 移动端详情使用抽屉或全屏弹层

## 异常与降级策略

- 无最新成功批次：展示空态与刷新入口
- 部分模型失败：展示错误摘要并保留成功模型内容
- Markdown 渲染失败：降级为纯文本展示
- 网络超时：展示重试入口

## 安全说明

- 文档中的 MongoDB 与 Redis 连接信息仅用于需求说明，不建议在生产环境中使用明文配置
- 模型 API Key 应加密存储并脱敏展示
- Swagger 与日志中不应输出密钥
- 后续应完善鉴权、限流、审计与 CORS 策略

## 当前实现状态摘要

当前项目已具备以下基础能力：

- 前后端基础工程可运行
- Swagger 文档可访问
- 提示词模板、模型配置、定时任务、系统配置基础管理界面与接口已存在
- 首页已支持展示最新一批 AI 分析结果
- 未开发完成的主标签页已显示占位提示
- 首页模型结果已调整为纵向排列，便于展示更多表格细节

## 参考文档

- [futures-functional-demand-backend.md](file:///d:/02-AICode/TradingChats-doc/futures-functional-demand-backend.md)
- [futures-functional-demand-frontend.md](file:///d:/02-AICode/TradingChats-doc/futures-functional-demand-frontend.md)
- [trading-chats-backend/README.md](file:///d:/02-AICode/TradingChats-doc/trading-chats-backend/README.md)
