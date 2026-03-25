# 凌期AI辅助期货挑选品种和持仓止盈止损分析程序

## 项目简介

这是一个基于Golang和Gin框架的后端应用程序，用于辅助期货挑选品种和持仓止盈止损分析。该程序通过调用AI模型生成分析结果，帮助用户做出更明智的交易决策。

## 技术栈

- **语言**: Go 1.20
- **框架**: Gin
- **数据库**: MongoDB, Redis
- **API文档**: Swagger

## 安装步骤

1. **克隆项目**

   ```bash
   git clone <repository-url>
   cd trading-chats-backend
   ```

2. **安装依赖**

   ```bash
   go mod download
   ```

3. **配置环境变量**

   复制 `.env.example` 文件为 `.env` 并根据实际情况修改配置：

   ```bash
   cp .env.example .env
   ```

   主要配置项：
   - `PORT`: 服务器端口
   - `MONGODB_URI`: MongoDB连接地址
   - `MONGODB_DATABASE`: MongoDB数据库名称
   - `REDIS_URI`: Redis连接地址

## 运行方法

### 开发环境

```bash
go run cmd/api/main.go
```

### 生产环境

1. **构建可执行文件**

   ```bash
   go build -o trading-chats-api cmd/api/main.go
   ```

2. **运行可执行文件**

   ```bash
   ./trading-chats-api
   ```

## API端点文档

启动服务后，可以通过以下地址访问Swagger API文档：

```
http://localhost:8080/swagger/index.html
```

主要API端点：

### 提示词模版

- `POST /api/prompt-templates`: 创建提示词模版
- `GET /api/prompt-templates`: 获取所有提示词模版
- `GET /api/prompt-templates/tag`: 根据标签获取提示词模版
- `GET /api/prompt-templates/{id}`: 根据ID获取提示词模版
- `PUT /api/prompt-templates/{id}`: 更新提示词模版
- `DELETE /api/prompt-templates/{id}`: 删除提示词模版
- `POST /api/prompt-templates/generate`: 动态生成提示词

### 模型与API配置

- `POST /api/model-api-configs`: 创建模型与API配置
- `GET /api/model-api-configs`: 获取所有模型与API配置
- `GET /api/model-api-configs/provider`: 根据提供商获取模型与API配置
- `GET /api/model-api-configs/{id}`: 根据ID获取模型与API配置
- `PUT /api/model-api-configs/{id}`: 更新模型与API配置
- `DELETE /api/model-api-configs/{id}`: 删除模型与API配置
- `POST /api/model-api-configs/{id}/test`: 测试模型的连通性

### AI响应信息

- `GET /api/ai-responses`: 获取所有AI响应信息
- `GET /api/ai-responses/batch`: 根据批次ID获取AI响应信息
- `GET /api/ai-responses/{id}`: 根据ID获取AI响应信息
- `POST /api/ai-responses/generate`: 生成批次AI响应

## 测试说明

### 运行单元测试

```bash
go test ./...
```

### 测试数据库连接

启动服务后，可以通过以下端点测试数据库连接：

```
http://localhost:8080/health
```

## 项目结构

```
trading-chats-backend/
├── cmd/
│   └── api/            # 应用程序入口
├── internal/
│   ├── api/            # HTTP处理程序和路由
│   ├── config/         # 配置管理
│   ├── db/             # 数据库连接
│   ├── models/         # 数据模型
│   ├── repository/     # 数据访问层
│   ├── service/        # 业务逻辑层
│   └── utils/          # 工具函数
├── pkg/                # 可重用包
├── .env                # 环境变量
├── go.mod              # Go模块文件
└── README.md           # 项目说明
```
