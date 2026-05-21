---
alwaysApply: false
description: 打包前端后端项目
---
## 后端项目打包

后端打包的路径为：D:\02-AICode\TradingChats-doc\trading-chats-backend\cmd\api

后端打包为linux版本命令：

```
bash
$env:GOOS="linux"; $env:GOARCH="amd64"; $env:CGO_ENABLED="0"; go build -o output/servers ./cmd/api
```

打包好后将文件名一定为server文件放到指定目录里面：D:\02-AICode\TradingChats-doc\serverBak\backend\build，

注意server文件一定要能在linux系统CPU 架构amd64中运行。

将swagger的文档目录里面文件D:\02-AICode\TradingChats-doc\trading-chats-backend\cmd\api\docs，放到指定目录里面：D:\02-AICode\TradingChats-doc\serverBak\backend\build\docs。


## 前端项目打包

路径为：D:\02-AICode\TradingChats-doc\trading-chats-frontend

打包命令为：npm run build

如果遇到编译错误进行排查问题并修复。

打包完成后将dist文件夹放到目录里面：D:\02-AICode\TradingChats-doc\serverBak\frontend
