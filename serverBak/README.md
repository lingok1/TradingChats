# TradingChats 服务部署与维护指南

## 快速部署

### 步骤 1：环境准备与权限设置
1. 确保位于项目根目录：`/home/TradingChats/serverBak`
2. 设置目录权限（如未设置）：
   ```bash
   chmod -R 777 /home/TradingChats/serverBak
   ```

### 步骤 2：检查当前运行状态
查看当前运行的 Docker 容器，确认服务状态：
```bash
docker ps
```
**预期输出**：如果服务已部署，应看到名为 `trading-chats-frontend` 和 `trading-chats-backend` 的容器。

### 步骤 3：清理旧环境（如需重新部署）
如果发现旧的容器正在运行，需要先停止并清理：
```bash
# 停止并删除容器（根据实际容器名称）
docker stop trading-chats-frontend trading-chats-backend
docker rm trading-chats-frontend trading-chats-backend

# 删除旧镜像（镜像名称为 serverbak-frontend/serverbak-backend）
docker rmi serverbak-frontend serverbak-backend
```
**注意**：容器名称可能为 `trading-chats-*`，而镜像名称为 `serverbak-*`，这是正常现象。

### 步骤 4：执行部署
运行自动化部署脚本：
```bash
cd /home/TradingChats/serverBak
./deploy.sh
```

## 部署脚本详解
`deploy.sh` 脚本自动执行以下流程：
1. **环境检查**：验证 Docker 和 docker compose 是否可用
2. **文件检查**：确认所有必需文件（backend/, frontend/, docker-compose.yml 等）存在
3. **清理网络**：移除旧的 Docker 网络（`serverbak_default`）
4. **构建镜像**：
   - `serverbak-backend`：基于 Alpine 的 Go 后端服务
   - `serverbak-frontend`：基于 Nginx 的前端服务
5. **启动服务**：通过 docker compose 启动容器
   - `trading-chats-backend`：监听 8080 端口
   - `trading-chats-frontend`：监听 80（HTTP）和 443（HTTPS）端口
6. **状态验证**：显示容器状态和最近日志

## 服务访问地址

### 生产环境访问
- **前端首页**：http://服务器IP/
- **Swagger API 文档**：http://服务器IP/swagger/index.html
- **后端直连**：http://服务器IP:8080/swagger/index.html

### 本地测试访问（如适用）
- 前端：http://localhost/
- 后端 API：http://localhost:8080

## 日常维护命令

### 查看服务状态
```bash
docker ps | grep trading-chats
```

### 查看服务日志
```bash
# 后端日志
docker logs trading-chats-backend --tail 50

# 前端日志
docker logs trading-chats-frontend --tail 50
```

### 重启服务
```bash
cd /home/TradingChats/serverBak
docker-compose restart
```

### 停止服务
```bash
cd /home/TradingChats/serverBak
docker-compose down
```

### 完全清理（删除所有容器、镜像、网络）
```bash
cd /home/TradingChats/serverBak
docker-compose down --rmi all --volumes
```

## 故障排除

### 常见问题
1. **权限不足**：运行 `chmod +x deploy.sh` 确保脚本可执行
2. **Docker 未安装**：请参考官方文档安装 Docker 和 docker-compose-plugin
3. **端口冲突**：检查 80、443、8080 端口是否被其他进程占用
4. **容器名称不匹配**：使用 `docker ps` 查看实际容器名称，相应调整命令

### 健康检查
- 后端健康检查：`curl -f http://localhost:8080/health`
- 前端访问测试：`curl -I http://localhost/`

## 文件结构说明
```
serverBak/
├── backend/          # 后端源代码
├── frontend/         # 前端构建文件
├── docker-compose.yml # Docker 编排配置
├── deploy.sh         # 自动化部署脚本
└── README.md         # 本文件
```

---
**最后更新**：2026-04-11
**基于实际部署经验修订**：确保文档与真实操作流程一致，可直接复用。