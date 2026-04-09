# docker-package

该目录用于在腾讯云 OpenCloudOS 8 服务器上直接构建并运行前后端镜像。

## 目录说明

- backend/build/server：本地编译好的 Linux amd64 后端二进制
- backend/docs：后端 Swagger 文档资源
- backend/Dockerfile：后端运行时镜像
- frontend/dist：前端已构建静态资源
- frontend/nginx.conf：前端 Nginx 反向代理配置
- frontend/Dockerfile：前端运行时镜像
- docker-compose.yml：前后端编排文件

## 在服务器上使用

进入该目录后执行：

```bash
docker compose build
docker compose up -d
```

## 访问地址

- 前端首页：<http://服务器IP/>
- 后端 Swagger：<http://服务器IP/swagger/index.html>

<br />

## 你在腾讯云 OpenCloudOS 8 上怎么用

上传整个 docker-package 到服务器后，进入目录执行：

```
chmod +x deploy.sh
./deploy.sh
```

<br />

## 说明

当前 compose 继续使用外部 MongoDB 和 Redis：

- MongoDB：150.158.18.92:27017
- Redis：150.158.18.92:6379

