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
scp -r C:\Users\20506\Desktop\中转站\linktree\* root@62.234.92.232:/opt/linktree/
```
必传文件：`templates/index.html`、`static/style.css`、`Dockerfile`（其余可一并传）。

### 2. SSH 登录并执行部署脚本
```bash
ssh root@62.234.92.232
cd /opt/linktree && chmod +x deploy.sh && ./deploy.sh
```
脚本自动：`docker compose up -d --build` → 等待 → 健康检查 → 容器状态。
数据目录 `/opt/linktree/data`（含 SQLite 数据库和上传文件）由 docker-compose
volume 挂载，重建镜像不会丢失。

### 3. 验证
- 公开页：http://62.234.92.232:3001/
- 管理后台：http://62.234.92.232:3001/admin/login
- 检查页面出现：关于我（含统计数字）、经历时间线、技能标签、作品过滤栏

## 本地预览
```powershell
cd C:\Users\20506\Desktop\中转站\linktree
$env:PORT="8080"; $env:ADMIN_PASSWORD="ChangeMe123!"; $env:SESSION_SECRET="dev"
go run .
# 浏览器打开 http://localhost:8080/
```

## 后续内容替换（TODO 占位内容位置）
| 内容 | 文件 | 位置 |
|---|---|---|
| 职位头衔/slogan/技能标签 | templates/index.html | tile-profile 内 |
| 简历下载路径 | templates/index.html | btn-resume 的 href（可改为 /api/files/download/{ID}） |
| 关于我文案/统计数字 | templates/index.html | tile-about 内 |
| 经历时间线 | templates/index.html | tile-timeline 内 |
| 作品分类 | templates/index.html | work-card 的 data-category 属性（web/ai/tool） |

## 敏感信息位置
服务器密码/管理密码/SSH 密钥路径保存在本地
`C:\Users\20506\.config\opencode\servers.json`（不进入 git 仓库）。