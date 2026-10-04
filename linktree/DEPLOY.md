# LinkTree 部署指南（服务器 62.234.92.232）

## 服务器连接
- 连接：`ssh root@62.234.92.232`（SSH 密钥免密，密钥 `~/.ssh/id_ed25519`）
- 系统：OpenCloudOS 9.4
- 部署目录：`/opt/linktree`
- 容器名：`linktree`

## 背景
服务器上运行的是旧版 CSS（V5 之前的），因为 Docker 镜像在构建时把 templates/ 和
static/ 复制进镜像，**必须重新构建镜像**才能让线上生效。另外原 Dockerfile 用
golang:1.23-alpine 构建 go.mod 要求 1.26.3 的项目会失败，现已修复为
golang:1.26-alpine + GOPROXY。

## 部署步骤

### 1. 上传最新代码到服务器（本地 Windows 执行）
```powershell
# 只传源码变更，不要 scp -r：web/node_modules 会把几万个文件塞过去（.gitignore/.dockerignore 已排除）
scp linktree\main.go linktree\models\settings.go linktree\docker-compose.yml root@62.234.92.232:/opt/linktree/
```
前端 `web/dist` 由镜像内 `node:20-alpine` 阶段构建，**不需要**本地上传。

### 2. 服务器凭据（首次或轮换时执行一次）
线上敏感配置外置在 `/opt/linktree/.env`（`chmod 600`，**不进 git**）：
```
SESSION_SECRET=<64 位随机串>
ADMIN_PASSWORD=<强随机密码>
```
`docker-compose.yml` 已改为 `env_file: .env`，仓库内不再有任何明文密钥。
⚠️ 顺序必须是**先写 `.env`，再改 compose/重建**；`main.go` 在 `SESSION_SECRET` 为空时 `log.Fatal` 拒绝启动。

### 3. SSH 登录并重建
```bash
ssh root@62.234.92.232
cd /opt/linktree && chmod +x deploy.sh && ./deploy.sh
```
脚本自动：`docker compose up -d --build` → 等待 → 健康检查 → 容器状态。
数据目录 `/opt/linktree/data`（含 SQLite 数据库和上传文件）由 docker-compose
volume 挂载，重建镜像不会丢失。

### 4. 后台密码轮换（改 `.env` 里的 `ADMIN_PASSWORD` 后必须显式执行一次）
环境变量只在数据库首次初始化时生效，因此**光改 `.env` 不会改密码**。轮换步骤：
```bash
cd /opt/linktree
docker compose stop linktree
set -a; . ./.env; set +a
# 该容器是常驻服务，用 --rm + 限时杀掉，不要让它前台挂着
timeout 20 docker run --rm --name linktree-rotate \
  -e PORT=8080 -e DATA_DIR=/app/data \
  -e SESSION_SECRET="$SESSION_SECRET" -e ADMIN_PASSWORD="$ADMIN_PASSWORD" \
  -v /opt/linktree/data:/app/data \
  linktree-linktree:latest ./linktree --rotate-admin-password
docker compose up -d
```
正常输出 `[admin-password-rotate] rotated: 已按 ADMIN_PASSWORD 重置后台密码`；
再次执行会输出 `skipped: 与当前密码一致`（幂等）。
`SESSION_SECRET` 变更会让所有已登录 Cookie 失效，需重新登录。

### 5. 验证
- 公开页：http://62.234.92.232:3001/
- 管理后台：http://62.234.92.232:3001/admin/login
- 登录接口自测：`POST /api/auth/login {"password":"<新密码>"}` 应 200，旧密码应 401
- 注意登录限流：同 IP 每分钟 5 次，连续试错会拿到 429

## 本地预览
```powershell
cd C:\Users\20506\Desktop\中转站\linktree
$env:PORT="8080"; $env:ADMIN_PASSWORD="devpass"; $env:SESSION_SECRET="dev-secret"
go run .
# 浏览器打开 http://localhost:8080/
```
沙箱内构建需把 Go 缓存指向工作区：`$env:GOCACHE="<工作区>\.dsh-tmp\gocache"`。

## 后续内容替换（TODO 占位内容位置）
| 内容 | 文件 | 位置 |
|---|---|---|
| 职位头衔/slogan/技能标签 | web/src/pages/HomePage.vue | hero 区 |
| 简历下载 | web/src/pages/HomePage.vue | 需先在后台「文件」上传 PDF 并设为公开 |
| 关于我文案/统计数字 | web/src/pages/AboutPage.vue | 数据块 |
| 经历时间线 | 后台「经历」面板 → `/api/experiences` | — |
| 作品分类 | 后台「作品」面板 | — |

## 敏感信息位置
- 线上密钥：`/opt/linktree/.env`（仅服务器，600 权限）
- 服务器连接信息：本地 `C:\Users\20506\.config\opencode\servers.json`（不进入 git）
- **历史遗留**：`SESSION_SECRET` 与 `ADMIN_PASSWORD=ChangeMe123!` 曾以明文提交（提交 `e0e9580` 起），
  已在 2026-10-04 轮换作废；`ai_api_key` 是明文存于 SQLite `settings` 表（未进 git）。
