# 项目后续修改指南

最后更新：2026-09-27

本文档用于快速理解 `C:\Users\20506\Desktop\中转站` 文件夹内的项目结构，方便后续继续修改、排查和扩展。

## 目录总览

当前文件夹内有两个主要项目、一个独立的内网模型代理配置文档，以及两个压缩包/部署产物：

| 路径 | 类型 | 说明 |
| --- | --- | --- |
| `linktree/` | Go Web 应用 | 个人链接页、文件分享、作品展示、AI 对话、后台管理 |
| `new-api-main/` | Go + React/Vite 应用 | New API 大模型网关与管理后台源码 |
| `new-api-deploy/` | 部署产物 | New API 已构建的二进制和前端 dist |
| `linktree.tar.gz` | 压缩包 | linktree 项目归档 |
| `new-api-deploy.tar.gz` | 压缩包 | new-api 部署包归档 |
| `公司内网模型代理配置.md` | 运维文档 | opencode、mihomo、qwen-relay 的内网模型代理链路 |

后续如果是源码修改，优先改 `linktree/` 或 `new-api-main/`；通常不要直接改 `new-api-deploy/`，除非目标是修改已经打包好的部署目录。

---

## linktree 项目

### 基本信息

- 项目路径：`C:\Users\20506\Desktop\中转站\linktree`
- 技术栈：Go 1.26.3、标准库 `net/http` 路由、`gorilla/sessions`、SQLite (`modernc.org/sqlite`)
- 默认端口：`8080`
- 数据目录：默认 `data/`
- 数据库：默认 `data/linktree.db`
- 入口文件：`main.go`
- 部署方式：本地 `go run`、Docker Compose、Dockerfile；部署细节见 `linktree/DEPLOY.md`

### 功能模块

| 模块 | 说明 | 主要文件 |
| --- | --- | --- |
| 公开站点前端 | Vue 3 + Vue Router 单页应用目标；后端通过 `web/dist` 提供 SPA 入口 | `web/`、`main.go` 中的 `spaHandler` |
| 公开数据 API | 公开资料、链接、作品、经历、分类和 AI 提问接口 | `handlers/public_api.go`、`handlers/work_api.go`、`handlers/experience_api.go`、`handlers/ask_api.go` |
| 后台管理 | 登录后管理链接、设置、分类、文件、作品和经历 | `templates/admin.html`、`handlers/` |
| 登录认证 | Cookie Session 登录、登出、鉴权中间件 | `handlers/api.go`、`middleware/auth.go` |
| 链接管理 | 链接增删改查和排序 | `models/link.go`、`handlers/api.go` |
| 分类管理 | 文件/作品分类 | `models/category.go`、`handlers/category_api.go` |
| 文件分享 | 上传、下载、公开/私有控制 | `models/file.go`、`handlers/file_api.go` |
| 作品展示 | 作品 CRUD、多媒体管理、外链跳转 | `models/work.go`、`handlers/work_api.go` |
| 经历展示 | 公开时间线与后台经历 CRUD | `models/experience.go`、`handlers/experience_api.go`、`templates/experience.html` |
| AI 提问 | 公开访客问答，流式返回 OpenAI 兼容模型结果 | `handlers/ask_api.go`、`templates/ask.html` |
| AI 对话 | 多会话、消息存储、调用 OpenAI 兼容接口 | `models/chat.go`、`handlers/chat_api.go`、`templates/chat.html` |
| 样式资源 | 页面样式、头像、背景、Logo | `static/` |

### 启动方式

在 `linktree` 目录下：

```powershell
go run main.go
```

可选环境变量：

| 变量 | 默认值 | 作用 |
| --- | --- | --- |
| `PORT` | `8080` | 服务端口 |
| `DATA_DIR` | `data` | SQLite 数据和上传文件存放目录 |
| `ADMIN_PASSWORD` | 自动生成 | 首次运行时的管理员密码 |
| `SESSION_SECRET` | 无默认值 | Session 加密密钥，缺失时拒绝启动 |
| `COOKIE_SECURE` | 未设置/非 `1` | 设为 `1` 后 Cookie 仅通过 HTTPS 发送 |
| `MAX_FILE_SIZE` | 代码固定 200MB | 当前 `GetMaxFileSize()` 未读取环境变量，上传接口实际固定为 200MB；Compose 中的同名变量目前只是配置记录 |
| `TZ` | 运行环境默认值 | Docker 中设为 `Asia/Shanghai` |

注意：当前版本在 `SESSION_SECRET` 未设置时会拒绝启动，不再使用内置默认密钥。生产环境必须设置随机密钥；Docker Compose 示例中的 `ADMIN_PASSWORD` 和 `SESSION_SECRET` 也应替换为部署环境专用值。

### 路由入口

所有路由集中注册在 `main.go`。页面入口已切换为 SPA 回退：

- 前端页面：`GET /` 及未命中显式路由的 `GET` 页面路径由 `spaHandler("web/dist")` 返回前端入口
- 管理登录页：`GET /admin/login`（仍使用 Go HTML 模板）
- 后台页：`GET /admin`（鉴权后仍使用 Go HTML 模板）
- 聊天页/API：`GET /chat`、`/api/chat/...`（鉴权后仍使用 Go HTML 模板）
- 公开资料和链接：`GET /api/public/profile`、`GET /api/public/links`
- AI 提问：`POST /api/ask`
- 链接 API：`/api/links`
- 设置 API：`/api/settings`
- 分类 API：`/api/categories`
- 文件 API：`/api/files`
- 作品 API：`/api/works`
- 经历 API：`GET /api/experiences`（公开读取）；`POST/PUT/DELETE /api/experiences...`（管理员）
- 聊天 API：`/api/chat/...`

`spaHandler` 对 `/assets/` 文件设置一年 `immutable` 缓存，对其他文件和 SPA 回退设置 `no-cache, must-revalidate`；未知 `/api...` 路径返回 JSON 404，避免误回退到前端 HTML。

公开下载接口为 `GET /api/files/download/{id}`。公开文件可匿名下载；私有文件对未登录请求统一返回 404，避免通过自增 ID 判断文件是否存在。

### 常见修改定位

| 想修改的内容 | 优先查看/修改 |
| --- | --- |
| 首页布局或文案 | `templates/index.html` |
| 新 SPA 页面与路由 | `web/src/`（页面组件、路由）及 `web/src/api.js` |
| 新 SPA 开发代理/构建 | `web/vite.config.js`、`web/package.json` |
| 后台界面 | `templates/admin.html` |
| 登录页 | `templates/login.html` |
| 聊天页 | `templates/chat.html` |
| 全局样式、颜色、背景 | `static/style.css` |
| 新增链接字段 | `models/link.go`、`handlers/api.go`、`templates/admin.html` |
| 新增作品字段 | `models/work.go`、`handlers/work_api.go`、`templates/admin.html` |
| 修改数据库表结构 | `models/db.go` 以及对应 model 文件 |
| 修改鉴权逻辑 | `middleware/auth.go`、`handlers/api.go` |
| 修改 AI 对话调用 | `handlers/chat_api.go`、`models/chat.go` |
| 修改访客 AI 提问 | `handlers/ask_api.go`、`templates/ask.html` |
| 修改经历时间线 | `models/experience.go`、`handlers/experience_api.go`、`templates/experience.html`、`templates/admin.html` |
| 修改公开站点页面路由 | Vue Router 前端路由与 `main.go` 的 SPA 回退 |
| 修改公开资料/链接数据 | `handlers/public_api.go`、`models/settings.go`、`models/link.go`、`web/src/api.js` |
| 修改上传安全策略 | `handlers/file_api.go`、`models/file.go` |
| 修改登录安全策略 | `handlers/api.go`、`main.go` |

### 注意事项

- `linktree.exe` 和 `linktree.exe~` 是编译产物，不是源码修改重点。
- `data/` 里可能包含本地数据库和上传文件，修改前应避免误删。
- SQLite 使用 WAL、5 秒 busy timeout，并将连接池限制为单写连接；新增/变更表结构应同步考虑索引和迁移。
- 上传接口使用 `http.MaxBytesReader` 限制请求总量，且拒绝 `application/octet-stream`；新增允许类型时要同时检查内容检测和可执行文件上传风险。
- 登录失败按 IP 在 1 分钟内最多允许 5 次失败尝试；`X-Forwarded-For` 只适合在可信反向代理后使用。
- `linktree/DEPLOY.md`、`linktree/deploy.sh` 是当前部署流程的参考入口；`server.log` 属于运行产物，不是源码文档。
- Vue 前端由 `web/package.json` 管理，开发服务器将 `/api` 代理到 `http://localhost:8080`；构建目录为 `web/dist`。
- 当前工作区中的 `web/index.html` 引用 `/src/main.js`，但已扫描到的 `web/src/` 目前只有 `api.js`，尚未发现入口文件或完整页面源码。前端仍处于迁移/搭建状态，不能仅凭 package 脚本判断可以完成构建。
- 当前 `Dockerfile` 仅复制 Go 二进制、`templates/` 和 `static/`；SPA 运行依赖 `web/dist`，部署前需确认镜像构建包含该目录，否则首页回退会找不到前端产物。
- 项目已有 `CODE_WIKI.md`，但当前终端读取时中文编码显示异常；后续以源码和当前总览文档为准。

### 当前部署配置

Docker 构建使用 Go 1.26 Alpine，并通过 `GOPROXY=https://goproxy.cn,direct` 下载依赖；运行时基于 Alpine 3.21，暴露容器端口 8080。现有 Compose 将宿主机 `3001` 映射到容器 `8080`，数据目录挂载到 `./data:/app/data`。

```powershell
cd C:\Users\20506\Desktop\中转站\linktree
docker compose up -d --build
```

本地开发仍可运行：

```powershell
cd C:\Users\20506\Desktop\中转站\linktree
$env:SESSION_SECRET = "replace-with-a-random-secret"
go run main.go
```

Vue 前端开发（在 `linktree/web` 目录）：

```powershell
bun install
bun run dev
bun run build
```

开发服务器默认将 API 代理到本地后端 `8080`。生产部署还需将 `web/dist` 纳入运行镜像或与 Go 服务一起发布；目前 Dockerfile 尚未复制该目录。

### 优化后的数据与性能行为

- 启动时会自动创建 `experiences` 表，并调用 `SeedDefaultExperiences()` 初始化默认经历数据。
- `works` 和 `work_media` 的约束迁移支持 `audio` 类型；作品列表加载媒体时改为单次批量查询，避免逐作品查询造成的 N+1。
- `files`、`works`、`work_media`、`chat_messages` 增加外键相关索引，SQLite 连接启用 `foreign_keys=1`、WAL 和 `synchronous=NORMAL`。
- HTTP Server 设置 `ReadHeaderTimeout=10s`、`IdleTimeout=120s`；未设置全局 `WriteTimeout`，以避免 SSE 聊天和大文件下载被截断。
- 静态资源响应增加 `Cache-Control: no-cache, must-revalidate`，方便部署后即时看到样式更新。

---

## new-api-main 项目

### 基本信息

- 项目路径：`C:\Users\20506\Desktop\中转站\new-api-main`
- 后端技术栈：Go 1.25.1、Gin、GORM、SQLite/MySQL/PostgreSQL、Redis 可选
- Default 前端：React 19、TypeScript、TanStack Router/Query、Rsbuild/Rspack、Tailwind CSS 4
- Classic 前端：React 18、JavaScript、Vite、Semi UI、Tailwind CSS 3、i18next
- 默认端口：`3000`
- 后端入口：`main.go`
- 路由目录：`router/`
- 前端源码：`web/default/src/`（Default）和 `web/classic/src/`（Classic）
- 默认嵌入前端产物：`web/default/dist`、`web/classic/dist`

### 后端启动流程

`main.go` 的核心流程：

1. 加载 `.env`
2. 初始化配置、日志、数据库、Redis、OAuth、模型比例等资源
3. 启动缓存同步、渠道测试、订阅重置、上游模型更新等后台任务
4. 创建 Gin server
5. 挂载中间件：RequestId、PoweredBy、I18n、日志、Session 等
6. 调用 `router.SetRouter(...)` 注册 API、转发接口和前端页面
7. 监听 `PORT`，默认 `3000`

### 关键目录

| 目录 | 作用 |
| --- | --- |
| `common/` | 通用配置、日志、环境变量、工具函数 |
| `constant/` | 全局常量、渠道类型、开关 |
| `controller/` | HTTP 控制器，处理后台 API、用户、渠道、支付、订阅等 |
| `middleware/` | Gin 中间件，鉴权、限流、CORS、统计、分发等 |
| `model/` | GORM 数据模型、数据库初始化、迁移、缓存 |
| `relay/` | 大模型请求转发适配层 |
| `router/` | 路由注册 |
| `service/` | 业务服务和后台任务 |
| `setting/` | 系统设置模块 |
| `web/classic/` | React Classic 前端源码与构建配置 |
| `web/default/` | React Default 前端源码与构建产物 |
| `docs/` | 项目文档和图片资源 |

### 路由结构

| 文件 | 说明 |
| --- | --- |
| `router/main.go` | 总路由入口，组合 API、Dashboard、Relay、Video、Web 路由 |
| `router/api-router.go` | 管理后台 API，例如用户、渠道、设置、支付、订阅、模型等 |
| `router/relay-router.go` | OpenAI/Claude/Gemini/Midjourney/Suno 等兼容转发接口 |
| `router/web-router.go` | 前端静态页面与主题资源 |
| `router/dashboard.go` | Dashboard 相关路由 |
| `router/video-router.go` | 视频代理/视频任务相关路由 |

### API 和业务定位

| 想修改的内容 | 优先查看/修改 |
| --- | --- |
| 用户注册、登录、用户管理 | `controller/user.go`、`router/api-router.go`、`model/user.go` |
| 渠道管理 | `controller/channel.go`、`model/channel.go`、`web/classic/src/pages/Channel` |
| 渠道测试 | `controller/channel-test.go` |
| 模型列表/模型管理 | `controller/model.go`、`controller/model_sync.go`、`web/classic/src/pages/Model` |
| API Key/令牌 | `controller/token.go`、`model/token.go`、`web/classic/src/pages/Token` |
| OpenAI 兼容接口转发 | `router/relay-router.go`、`controller/relay.go`、`relay/` |
| Claude/Gemini 格式转换 | `relay/`、`types/`、`router/relay-router.go` |
| Midjourney/Suno 任务 | `controller/midjourney.go`、`controller/task.go`、`relay/` |
| 支付/充值 | `controller/topup*.go`、`controller/billing.go`、`controller/payment_*.go` |
| 订阅 | `controller/subscription*.go`、`service/`、`model/` |
| 系统设置 | `controller/option.go`、`setting/`、`web/classic/src/pages/Setting` |
| OAuth | `oauth/`、`controller/oauth.go`、`controller/custom_oauth.go` |
| 性能监控 | `controller/performance.go`、`pkg/perf_metrics` |

### 数据库

数据库初始化在 `model/main.go`：

- `SQL_DSN` 为空时使用 SQLite。
- `SQL_DSN` 以 `postgres://` 或 `postgresql://` 开头时使用 PostgreSQL。
- 其他非空 `SQL_DSN` 默认按 MySQL 处理。
- `LOG_SQL_DSN` 可单独配置日志数据库；为空时日志表和主数据库共用连接。
- 迁移通过 GORM `AutoMigrate` 和若干手写迁移函数完成。

重要环境变量：

| 变量 | 作用 |
| --- | --- |
| `PORT` | 后端监听端口 |
| `FRONTEND_BASE_URL` | 非主节点时可将 Web 路由重定向到独立前端 |
| `SQL_DSN` | 主数据库连接 |
| `LOG_SQL_DSN` | 日志数据库连接 |
| `SQLITE_PATH` | SQLite 数据库路径 |
| `REDIS_CONN_STRING` | Redis 连接 |
| `SESSION_SECRET` | Session 密钥 |
| `SYNC_FREQUENCY` | 缓存同步频率 |
| `CHANNEL_UPDATE_FREQUENCY` | 渠道上游更新频率 |
| `BATCH_UPDATE_ENABLED` | 是否启用批量更新 |
| `ENABLE_PPROF` | 是否启用 pprof |

### 前端界面与主题

New API 当前维护两套可嵌入的前端：`default` 是 React 19 + TypeScript 的 TanStack Router/Rsbuild 新版界面；`classic` 是 React 18 + JavaScript 的 Vite/Semi UI 界面。两者分别构建到 `web/default/dist` 和 `web/classic/dist`，并在 `main.go` 中通过 `go:embed` 嵌入后端。

`router/web-router.go` 创建主题感知的静态文件系统，并在 SPA 回退时根据 `common.GetTheme()` 返回对应主题的 `index.html`。如需改主题选择、静态资源或 SPA 回退行为，联查 `router/main.go`、`router/web-router.go`、`common` 中主题配置，以及两个前端各自的构建产物路径。

#### Default 前端

路径：`new-api-main/web/default`

| 路径 | 作用 |
| --- | --- |
| `src/main.tsx` | React、TanStack Router、React Query、主题/字体/方向 provider 初始化 |
| `src/routes/` | TanStack 文件路由；`routeTree.gen.ts` 为自动生成文件，不手工编辑 |
| `src/features/` | 按功能组织的页面业务模块，如渠道、聊天、仪表盘、模型、设置、用户、用量日志、钱包等 |
| `src/components/` | 布局、数据表和通用 UI 组件 |
| `src/lib/` | API、错误处理、日期及通用工具 |
| `src/stores/` | Zustand 状态管理 |
| `src/i18n/` | 国际化配置与语言资源 |
| `src/styles/` | Tailwind 与主题样式 |
| `rsbuild.config.ts` | Rsbuild/Rspack 配置、TanStack Router 插件、代理与分包 |
| `AGENTS.md` | Default 前端开发规范，改该目录前先阅读 |

常用命令（在 Default 前端目录执行）：

```powershell
bun run dev
bun run typecheck
bun run build
bun run build:check
bun run lint
```

开发服务器将 `/api`、`/mj`、`/pg` 代理到 `VITE_REACT_APP_SERVER_URL`；未设置时默认 `http://localhost:3000`。生产构建启用压缩、移除 `console.log`，并按路由分包；`build:check` 会先运行 TypeScript 项目检查再构建。

#### Classic 前端

路径：`new-api-main/web/classic`

关键文件：

| 文件/目录 | 作用 |
| --- | --- |
| `package.json` | 前端依赖与脚本 |
| `vite.config.js` | Vite 配置，含 `/api`、`/mj`、`/pg` 代理 |
| `src/index.jsx` | React 入口 |
| `src/components/layout/PageLayout` | 应用总布局和路由承载 |
| `src/pages/` | 页面级组件 |
| `src/components/table/` | 表格型管理模块 |
| `src/components/settings/` | 设置页组件 |
| `src/context/` | 用户、状态、主题上下文 |
| `src/helpers/` | API 请求、格式化、鉴权等辅助函数 |
| `src/i18n/locales/` | 多语言文案 |
| `src/index.css` | 全局样式 |

前端命令：

```powershell
cd C:\Users\20506\Desktop\中转站\new-api-main\web\classic
bun run dev
bun run build
bun run lint
bun run eslint
```

Vite 开发服务器默认把这些路径代理到后端 `http://localhost:3000`：

- `/api`
- `/mj`
- `/pg`

### New API 前端构建

Makefile 中 `build-frontend` 构建 Default 前端，`build-frontend-classic` 构建 Classic 前端，`build-all-frontends` 构建两套前端。后端使用 `go:embed`，因此部署带有前端代码变更的新版时，应先生成相应 `dist`，再重新构建后端二进制或镜像；仅更新 Go 二进制而没有新嵌入产物，不会更新前端页面。

Default 主要功能区域按 `src/features/<feature>/` 组织，查找后台新页面时先检查 `src/routes/` 的路由文件，再进入对应 feature；页面路由和功能实现通常分开维护。

### New API 启动与构建

后端本地运行：

```powershell
cd C:\Users\20506\Desktop\中转站\new-api-main
go run main.go
```

Docker Compose：

```powershell
cd C:\Users\20506\Desktop\中转站\new-api-main
docker compose up -d
```

Makefile 常用目标：

```powershell
make dev-api
make dev-web-classic
make build-frontend-classic
make build-all-frontends
```

注意：此项目的 Makefile 是类 Unix shell 风格，在 Windows PowerShell 下可能需要 Git Bash、WSL 或手动执行对应命令。

### 常见修改定位

| 修改目标 | 后端位置 | 前端位置 |
| --- | --- | --- |
| 新增后台 API | `router/api-router.go`、`controller/`、`model/` | 对应 `src/pages` 或 `src/components` |
| 新增转发接口 | `router/relay-router.go`、`controller/relay.go`、`relay/` | 通常无需前端，除非要加管理入口 |
| 改某个后台页面 | 对应 `controller` API | `web/classic/src/pages/<页面>` |
| 改表格列/操作按钮 | 对应 API | `web/classic/src/components/table/<模块>` |
| 改系统设置项 | `controller/option.go`、`setting/` | `web/classic/src/pages/Setting`、`src/components/settings` |
| 改渠道类型或适配 | `constant/`、`relay/`、`controller/channel.go` | 渠道编辑弹窗、模型选择组件 |
| 改登录注册 | `controller/user.go`、`middleware/` | `web/classic/src/pages/User`、认证组件 |
| 改多语言文案 | 后端 `i18n/` | `web/classic/src/i18n/locales/*.json` |

若目标主题为 Default，将前端位置替换为 `web/default/src/routes/`、`web/default/src/features/`、`web/default/src/i18n/`；若需支持两种主题，应分别实现或确认对应页面和翻译都存在。

### 注意事项

- `new-api.exe`、`new-api-linux`、`web/*/dist` 是构建产物，源码修改一般不要从这些文件入手；后端通过 `go:embed` 打包两个前端的 dist，发布 UI 变更时需一并重新构建。
- `.env` 可能包含本地配置或敏感信息，提交或分享前应检查。
- `one-api.db` 是本地 SQLite 数据库，修改数据库迁移前建议备份。
- `README.zh_CN.md` 在当前终端环境显示乱码，应该是编码显示问题；后续以源码和官方 README 原始文件为准。
- 项目采用 AGPL-3.0 许可，做公开部署或二次分发时需要留意许可证要求。

---

## 后续协作建议

如果要继续修改，可以直接按目标描述：

- “修改 linktree 首页样式”
- “给 linktree 作品增加标签字段”
- “修改 new-api 的渠道管理页面”
- “给 new-api 新增一个接口”
- “排查 new-api 转发 OpenAI 接口失败”

我会优先根据本文档中的路径定位代码，然后再读具体文件并实施修改。

---

## 公司内网模型代理配置文档

根目录的 `公司内网模型代理配置.md` 记录的是本机 opencode 使用公司内网 Qwen 模型的运维链路，不属于两个 Go 项目的源码模块。它描述了以下固定关系：

```text
opencode company-qwen
  -> qwen-relay 127.0.0.1:18449
  -> mihomo 127.0.0.1:18080
  -> fh Shadowsocks 网关
  -> 10.30.0.43:8449
```

后续修改代理时，优先查看该文档中的组件清单和故障排查表。重点配置位置在用户目录下，不在本工作文件夹内：mihomo `config.yaml`、`qwen-relay.ts`、opencode `opencode.json` 与 `auto-proxy.ts`。其中 SS 密码属于敏感信息，不应复制到仓库或提交到版本控制。
