#!/bin/bash
# ============================================================
# LinkTree 部署脚本（在服务器 62.234.92.232 上执行）
# 部署目录：/opt/linktree
# 作用：重新构建镜像（打包最新模板+静态文件）并重启容器
# 前提：最新源码已上传到 /opt/linktree（templates/ static/ 等）
# 注意：数据目录 /opt/linktree/data 由 docker-compose 挂载，不会丢失
# ============================================================
set -e

cd /opt/linktree

echo "==> 1/4 重新构建镜像（--build 会使用最新 Dockerfile + 源码）"
docker compose up -d --build linktree

echo "==> 2/4 等待服务就绪"
sleep 5

echo "==> 3/4 健康检查"
curl -s -o /dev/null -w "HTTP %{http_code}\n" http://localhost:3001/ || echo "检查失败：请确认服务状态"

echo "==> 4/4 容器状态"
docker ps --filter name=linktree --format "{{.Names}} {{.Status}}"

echo "==> 完成。访问 http://62.234.92.232:3001/"