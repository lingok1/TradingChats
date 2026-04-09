#!/usr/bin/env bash
set -e

PROJECT_DIR="$(cd "$(dirname "$0")" && pwd)"

echo "==> 进入项目目录"
cd "$PROJECT_DIR"

echo "==> 检查 Docker 是否安装"
if ! command -v docker >/dev/null 2>&1; then
  echo "Docker 未安装，请先安装 Docker 和 Docker Compose"
  exit 1
fi

echo "==> 检查 docker compose 是否可用"
if ! docker compose version >/dev/null 2>&1; then
  echo "docker compose 不可用，请确认已安装 docker-compose-plugin"
  exit 1
fi

echo "==> 当前目录文件检查"
ls

echo "==> 检查关键文件"
if test -f docker-compose.yml; then echo "File exists" ; else echo "File missing - continuing anyway" ; fi
if test -f backend/Dockerfile; then echo "File exists" ; else echo "File missing - continuing anyway" ; fi
if test -f backend/build/server; then echo "File exists" ; else echo "File missing - continuing anyway" ; fi
if test -f frontend/Dockerfile; then echo "File exists" ; else echo "File missing - continuing anyway" ; fi
if test -f frontend/nginx.conf; then echo "File exists" ; else echo "File missing - continuing anyway" ; fi
if test -f frontend/dist/index.html; then echo "File exists" ; else echo "File missing - continuing anyway" ; fi

echo "==> 停止旧容器（如存在）"
docker compose down || true

echo "==> 构建镜像"
docker compose build

echo "==> 启动服务"
docker compose up -d

echo "==> 查看容器状态"
docker compose ps

echo
echo "==> 后端日志（最近 50 行）"
docker compose logs --tail=50 backend || true

echo
echo "==> 前端日志（最近 50 行）"
docker compose logs --tail=50 frontend || true

echo
echo "==> 部署完成"
echo "前端首页:  http://服务器IP/"
echo "Swagger:   http://服务器IP/swagger/index.html"
echo "后端直连: http://服务器IP:8080/swagger/index.html"
echo
echo "如果需要执行权限，请运行：chmod +x deploy.sh"
