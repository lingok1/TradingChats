---
alwaysApply: false
description: 运行启动停止重启前端后端项目
---
前端只能在端口：5173运行，后端只能在端口：8080运行。
运行启动前端后端项目前检查一下端口占用情况并处理占用端口的情况

- 前端项目占用端口：5173
- 后端项目占用端口：8080
- 如果占用端口，需要先关闭占用端口的进程

检查一下前后端端口占用情况的命令类似为：
netstat -ano | findstr :8080
结束这个进程类似命令为：
taskkill /PID 63364 /F

后端项目路径为：trading-chats-backend\cmd\api\main.go
运行后端项目的命令为：go run cmd/api/main.go
后端接口文档地址是：
http://localhost:8080/swagger/index.html

- 前端项目目录： d:\02-AICode\TradingChats-doc\trading-chats-frontend
- 启动命令： npm run dev
- 服务状态：正在运行
- 访问地址： http://localhost:5173/
