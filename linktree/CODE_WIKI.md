# LinkTree Code Wiki

> **项目名称**: LinkTree
> **语言**: Go 1.26.3
> **版本**: v2.0
> **数据库**: SQLite (modernc.org/sqlite — 纯 Go 实现，无需 CGO)
> **架构模式**: 经典 MVC（Model-View-Handler）

---

## 一、项目概述

LinkTree 是一个个人链接导航站 + 文件分享 + 作品展示 + AI 对话的 All-in-One Web 应用。用户可以创建一个类似 Linktree 的个人主页，展示常用链接、分享文件、展示作品集，并内置了一个基于 OpenAI 兼容 API 的 AI 对话功能。

### 核心功能

| 功能模块 | 说明 |
|---------|------|
| 链接管理 | 创建、编辑、删除、排序个人链接 |
| 文件存储 | 上传文件、分类管理、公开/私有控制、下载计数 |
| 作品展示 | 创建作品集、封面图、多媒体内容、外链跳转 |
| AI 对话 | 多会话管理、流式输出、支持自定义模型和 API |
| 主题切换 | 亮色/暗色/自动模式，毛玻璃 UI 风格 |
| 管理后台 | Tab 导航式后台，管理所有内容 |

---

## 二、项目目录结构

```
linktree/
├── main.go                     # 程序入口，路由注册，服务启动
├── go.mod                      # Go 模块定义
├── go.sum                      # 依赖校验
├── Dockerfile                  # 多阶段构建镜像
├── docker-compose.yml          # Docker Compose 编排
├── .gitignore                  # Git 忽略规则
│
├── models/                     # 数据模型层（M）
│   ├── db.go                   # 数据库初始化 + 建表
│   ├── link.go                 # 链接模型
│   ├── category.go             # 分类模型
│   ├── file.go                 # 文件模型
│   ├── work.go                 # 作品模型
│   ├── chat.go                 # AI 对话模型
│   └── settings.go             # 系统设置 + 密码管理
│
├── handlers/                   # 请求处理层（C）
│   ├── api.go                  # 认证 + 链接 + 设置 API
│   ├── page.go                 # 页面渲染
│   ├── category_api.go         # 分类 API
│   ├── file_api.go             # 文件上传/下载/删除 API
│   ├── work_api.go             # 作品管理 API
│   └── chat_api.go             # AI 对话 API（SSE 流式）
│
├── middleware/                  # 中间件
│   └── auth.go                 # 认证中间件
│
├── templates/                  # 模板层（V）
│   ├── index.html              # 公开主页
│   ├── admin.html              # 管理后台
│   ├── login.html              # 登录页
│   └── chat.html               # AI 对话页
│
├── static/                     # 静态资源
│   ├── style.css               # 全局样式（亮/暗色主题）
│   ├── avatar.jpg              # 默认头像
│   ├── bg.jpg                  # 背景图
│   └── logo.jpg                # Logo
│
└── docs/                       # 项目文档
    ├── DEVELOPMENT.md          # 开发文档
    └── superpowers/            # 规划与设计文档
```

---

## 三、整体架构

### 3.1 架构图

```
┌─────────────────────────────────────────────────────┐
│                     浏览器客户端                       │
│  index.html / admin.html / login.html / chat.html   │
└──────────────────────┬──────────────────────────────┘
                       │ HTTP / SSE
┌──────────────────────▼──────────────────────────────┐
│                   main.go (路由层)                     │
│  ┌──────────┐  ┌──────────┐  ┌──────────────────┐   │
│  │ 公开路由  │  │ 认证路由  │  │ 受保护路由(auth)  │   │
│  └────┬─────┘  └────┬─────┘  └───────┬──────────┘   │
└───────┼──────────────┼────────────────┼──────────────┘
        │              │                │
┌───────▼──────────────▼────────────────▼──────────────┐
│               middleware/auth.go                       │
│           CheckAuth (Session 认证中间件)                │
└──────────────────────┬──────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────┐
│                  handlers/ (处理层)                    │
│  api.go │ page.go │ category_api.go │ file_api.go   │
│  work_api.go │ chat_api.go                            │
└──────────────────────┬──────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────┐
│                 models/ (数据层)                       │
│  db.go │ link.go │ category.go │ file.go │ work.go  │
│  chat.go │ settings.go                                │
└──────────────────────┬──────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────┐
│              SQLite (data/linktree.db)                │
│  links │ settings │ categories │ files │ works       │
│  work_media │ chat_sessions │ chat_messages          │
└─────────────────────────────────────────────────────┘
```

### 3.2 请求处理流程

1. 用户请求到达 `main.go` 中注册的路由
2. 受保护路由经过 `middleware.CheckAuth` 中间件验证 Session
3. 路由分发到对应的 `handlers` 函数
4. Handler 调用 `models` 层进行数据操作
5. 返回 JSON（API 请求）或渲染 HTML 模板（页面请求）

### 3.3 认证机制

- 使用 `gorilla/sessions` 基于 Cookie 的 Session 认证
- Session 名称: `linktree_session`
- Session 有效期: 24 小时 (86400s)
- 密码使用 `bcrypt` 加密存储在 `settings` 表中
- 首次运行自动生成随机密码（或通过环境变量 `ADMIN_PASSWORD` 指定）
- 中间件对 API 请求返回 401 JSON，对页面请求重定向到登录页

---

## 四、数据模型详解

### 4.1 数据库表结构

#### links（链接表）

| 字段 | 类型 | 说明 |
|-----|------|------|
| id | INTEGER PK | 自增主键 |
| title | TEXT | 链接标题 |
| url | TEXT | 链接地址 |
| sort_order | REAL | 排序权重 |
| icon | TEXT | 图标标识 |
| created_at | TEXT | 创建时间 |

#### settings（设置表）

| 字段 | 类型 | 说明 |
|-----|------|------|
| key | TEXT PK | 设置键名 |
| value | TEXT | 设置值 |

**已使用的设置键**:
- `admin_password_hash` — 管理员密码哈希
- `display_name` — 显示名称
- `bio` — 个人简介
- `avatar_data_url` — 头像 Data URL
- `theme` — 主题模式 (auto/light/dark)
- `ai_base_url` — AI API 基础 URL
- `ai_api_key` — AI API 密钥
- `ai_model` — AI 模型名称

#### categories（分类表）

| 字段 | 类型 | 说明 |
|-----|------|------|
| id | INTEGER PK | 自增主键 |
| name | TEXT | 分类名称 |
| type | TEXT | 分类类型: `file` 或 `work` |
| sort_order | REAL | 排序权重 |
| created_at | TEXT | 创建时间 |

#### files（文件表）

| 字段 | 类型 | 说明 |
|-----|------|------|
| id | INTEGER PK | 自增主键 |
| category_id | INTEGER FK | 关联分类（可为空） |
| original_name | TEXT | 原始文件名 |
| stored_name | TEXT | 存储文件名（UUID） |
| mime_type | TEXT | MIME 类型 |
| size | INTEGER | 文件大小（字节） |
| is_public | INTEGER | 是否公开 (0/1) |
| download_count | INTEGER | 下载次数 |
| sort_order | REAL | 排序权重 |
| created_at | TEXT | 创建时间 |

#### works（作品表）

| 字段 | 类型 | 说明 |
|-----|------|------|
| id | INTEGER PK | 自增主键 |
| category_id | INTEGER FK | 关联分类（可为空） |
| title | TEXT | 作品标题 |
| description | TEXT | 作品描述 |
| cover_url | TEXT | 封面图 URL |
| content_type | TEXT | 内容类型: `image`/`video`/`mixed` |
| external_url | TEXT | 外部链接 |
| sort_order | REAL | 排序权重 |
| created_at | TEXT | 创建时间 |

#### work_media（作品媒体关联表）

| 字段 | 类型 | 说明 |
|-----|------|------|
| id | INTEGER PK | 自增主键 |
| work_id | INTEGER FK | 关联作品（CASCADE 删除） |
| file_id | INTEGER FK | 关联文件（可为空） |
| media_url | TEXT | 媒体 URL |
| media_type | TEXT | 媒体类型: `image`/`video` |
| sort_order | REAL | 排序权重 |

#### chat_sessions（对话会话表）

| 字段 | 类型 | 说明 |
|-----|------|------|
| id | INTEGER PK | 自增主键 |
| title | TEXT | 会话标题 |
| model | TEXT | 使用的 AI 模型 |
| created_at | TEXT | 创建时间 |
| updated_at | TEXT | 更新时间 |

#### chat_messages（对话消息表）

| 字段 | 类型 | 说明 |
|-----|------|------|
| id | INTEGER PK | 自增主键 |
| session_id | INTEGER FK | 关联会话（CASCADE 删除） |
| role | TEXT | 角色: `user`/`assistant`/`system` |
| content | TEXT | 消息内容 |
| created_at | TEXT | 创建时间 |

### 4.2 表关系图

```
categories ──1:N──→ files
categories ──1:N──→ works
works ───────1:N──→ work_media ──→ files (可选关联)
chat_sessions ──1:N──→ chat_messages
settings (独立 KV 表)
links (独立表)
```

---

## 五、模块详解

### 5.1 models/ — 数据模型层

#### db.go — 数据库初始化

| 函数 | 说明 |
|-----|------|
| `InitDB(dbPath string) error` | 打开 SQLite 数据库，启用 WAL 模式和外键约束，创建全部 7 张表 |

**关键配置**:
- `PRAGMA journal_mode=WAL` — 启用 WAL 日志模式，提升并发读性能
- `PRAGMA foreign_keys=ON` — 启用外键约束

#### link.go — 链接模型

| 结构体/函数 | 说明 |
|------------|------|
| `Link` | 链接结构体: ID, Title, URL, SortOrder, Icon |
| `OrderEntry` | 排序条目: ID, SortOrder |
| `GetAllLinks() ([]Link, error)` | 获取所有链接，按 sort_order 升序 |
| `CreateLink(title, url string) (Link, error)` | 创建链接，自动计算下一个排序值 |
| `UpdateLink(id int, title, url string, sortOrder float64) error` | 更新链接 |
| `DeleteLink(id int) error` | 删除链接 |
| `ReorderLinks(orders []OrderEntry) error` | 批量更新排序（事务） |
| `SeedDefaultLinks()` | 首次运行时插入默认链接（CSDN, GitHub, 中转站） |

#### category.go — 分类模型

| 结构体/函数 | 说明 |
|------------|------|
| `Category` | 分类结构体: ID, Name, Type, SortOrder |
| `GetCategoriesByType(categoryType string) ([]Category, error)` | 按类型获取分类 |
| `CreateCategory(name, categoryType string) (Category, error)` | 创建分类 |
| `UpdateCategory(id int, name string, sortOrder float64) error` | 更新分类 |
| `DeleteCategory(id int) error` | 删除分类 |

#### file.go — 文件模型

| 结构体/函数 | 说明 |
|------------|------|
| `File` | 文件结构体: ID, CategoryID, OriginalName, StoredName, MimeType, Size, FormattedSize, IsPublic, DownloadCount, SortOrder, CreatedAt |
| `allowedMimeTypes` | 允许上传的 MIME 类型白名单 |
| `IsAllowedMimeType(mimeType string) bool` | 检查 MIME 类型是否允许 |
| `GetFilesByCategory(categoryID *int) ([]File, error)` | 按分类获取文件 |
| `GetPublicFiles() ([]File, error)` | 获取所有公开文件 |
| `GetFileByID(id int) (*File, error)` | 按 ID 获取文件 |
| `CreateFile(...) (File, error)` | 创建文件记录 |
| `IncrementDownloadCount(id int) error` | 增加下载计数 |
| `DeleteFile(id int) error` | 删除文件记录 |
| `GenerateStoredName(originalName string) string` | 生成 UUID 存储文件名 |
| `GetUploadDir() string` | 获取上传目录路径 |
| `GetMaxFileSize() int64` | 获取最大文件大小 (200MB) |
| `FormatFileSize(size int64) string` | 格式化文件大小显示 |

**支持的文件类型**: jpeg, png, gif, webp, svg, pdf, txt, doc, docx, xls, xlsx, zip, rar, gz, mp4, webm, ogg

#### work.go — 作品模型

| 结构体/函数 | 说明 |
|------------|------|
| `Work` | 作品结构体: ID, CategoryID, Title, Description, CoverURL, ContentType, ExternalURL, SortOrder, CreatedAt, Media |
| `WorkMedia` | 作品媒体结构体: ID, WorkID, FileID, MediaURL, MediaType, SortOrder |
| `GetWorksByCategory(categoryID *int) ([]Work, error)` | 按分类获取作品 |
| `GetPublicWorks() ([]Work, error)` | 获取所有公开作品 |
| `GetWorkByID(id int) (*Work, error)` | 按 ID 获取作品 |
| `CreateWork(...) (Work, error)` | 创建作品 |
| `UpdateWork(...) error` | 更新作品 |
| `DeleteWork(id int) error` | 删除作品 |
| `GetWorkMedia(workID int) ([]WorkMedia, error)` | 获取作品媒体列表 |
| `AddWorkMedia(workID int, fileID *int, mediaURL, mediaType string) (WorkMedia, error)` | 添加作品媒体 |
| `DeleteWorkMedia(id int) error` | 删除作品媒体 |
| `LoadWorkMedia(works []Work) error` | 批量加载作品媒体 |

#### chat.go — AI 对话模型

| 结构体/函数 | 说明 |
|------------|------|
| `ChatSession` | 对话会话: ID, Title, Model, MessageCount, CreatedAt, UpdatedAt |
| `ChatMessage` | 对话消息: ID, SessionID, Role, Content, CreatedAt |
| `GetAllChatSessions() ([]ChatSession, error)` | 获取所有会话（含消息计数） |
| `CreateChatSession(title, model string) (ChatSession, error)` | 创建会话 |
| `UpdateChatSession(id int, title, model string) error` | 更新会话 |
| `UpdateChatSessionTime(id int) error` | 更新会话时间戳 |
| `DeleteChatSession(id int) error` | 删除会话 |
| `GetChatSessionByID(id int) (*ChatSession, error)` | 按 ID 获取会话 |
| `GetChatMessages(sessionID int) ([]ChatMessage, error)` | 获取会话消息 |
| `GetRecentChatMessages(sessionID int, limit int) ([]ChatMessage, error)` | 获取最近 N 条消息（倒序查询后反转） |
| `CreateChatMessage(sessionID int, role, content string) (ChatMessage, error)` | 创建消息（同时更新会话时间） |
| `DeleteChatMessages(sessionID int) error` | 删除会话所有消息 |
| `GetChatMessageCount(sessionID int) int` | 获取消息数量 |
| `TruncateChatMessages(sessionID int, keep int) error` | 保留最近 N 条消息，删除其余 |

#### settings.go — 系统设置

| 函数 | 说明 |
|-----|------|
| `GetSetting(key string) (string, error)` | 获取设置值 |
| `GetSettingOrDefault(key, defaultVal string) string` | 获取设置值（带默认值） |
| `SetSetting(key, value string) error` | 设置值（INSERT OR REPLACE） |
| `InitAdminPassword(envPassword string) (string, error)` | 初始化管理员密码（首次运行时） |
| `VerifyPassword(password string) bool` | 验证密码（bcrypt） |
| `ChangePassword(newPassword string) error` | 修改密码 |
| `generateRandomPassword(length int) string` | 生成随机密码（crypto/rand） |
| `MaskAPIKey(key string) string` | API Key 脱敏显示 |
| `GetAIBaseURL() string` | 获取 AI API 基础 URL（默认 DeepSeek） |
| `GetAIAPIKey() string` | 获取 AI API Key |
| `GetAIModel() string` | 获取 AI 模型名称（默认 deepseek-v4-flash） |

---

### 5.2 handlers/ — 请求处理层

#### api.go — 认证 + 链接 + 设置

| 函数 | 路由 | 方法 | 说明 |
|-----|------|------|------|
| `SetStore(s *sessions.CookieStore)` | — | — | 注入 Session Store |
| `writeJSON(w, data)` | — | — | 工具函数：写入 JSON 响应 |
| `Login` | `/api/auth/login` | POST | 登录认证，设置 Session |
| `Logout` | `/api/auth/logout` | POST | 退出登录，清除 Session |
| `ListLinks` | `/api/links` | GET | 获取链接列表 |
| `CreateLink` | `/api/links` | POST | 创建链接 |
| `UpdateLink` | `/api/links/{id}` | PUT | 更新链接 |
| `DeleteLink` | `/api/links/{id}` | DELETE | 删除链接 |
| `ReorderLinks` | `/api/links/reorder` | PUT | 批量排序链接 |
| `GetSettings` | `/api/settings` | GET | 获取系统设置（API Key 脱敏） |
| `UpdateSettings` | `/api/settings` | PUT | 更新设置（支持部分更新） |

**API 响应格式**:
- 成功: `{"ok": true, ...}` 或直接返回数据数组
- 失败: `{"ok": false, "message": "错误信息"}`

#### page.go — 页面渲染

| 函数 | 路由 | 说明 |
|-----|------|------|
| `InitTemplates() error` | — | 解析 templates/*.html 模板 |
| `IndexPage` | `GET /` | 渲染公开主页（链接+作品+文件） |
| `LoginPage` | `GET /admin/login` | 渲染登录页 |
| `AdminPage` | `GET /admin` | 渲染管理后台 |
| `ChatPage` | `GET /chat` | 渲染 AI 对话页 |

**PageData 结构体**: 传递给 index.html 的数据，包含 DisplayName, Bio, AvatarURL, Links, Works, Files, Theme。

#### category_api.go — 分类 API

| 函数 | 路由 | 方法 | 说明 |
|-----|------|------|------|
| `ListCategories` | `/api/categories?type=` | GET | 获取分类列表（公开） |
| `CreateCategory` | `/api/categories` | POST | 创建分类（需认证） |
| `UpdateCategory` | `/api/categories/{id}` | PUT | 更新分类（需认证） |
| `DeleteCategory` | `/api/categories/{id}` | DELETE | 删除分类（需认证） |

#### file_api.go — 文件 API

| 函数 | 路由 | 方法 | 说明 |
|-----|------|------|------|
| `ListFiles` | `/api/files` | GET | 获取文件列表（支持分类筛选） |
| `ListPublicFiles` | `/api/files/public` | GET | 获取公开文件列表 |
| `UploadFile` | `/api/files/upload` | POST | 上传文件（multipart/form-data） |
| `DownloadFile` | `/api/files/download/{id}` | GET | 下载文件（增加下载计数） |
| `DeleteFile` | `/api/files/{id}` | DELETE | 删除文件（物理文件+数据库记录） |

**上传流程**:
1. 解析 multipart 表单
2. 检测 MIME 类型（先内容检测，后扩展名兜底）
3. 白名单验证 MIME 类型
4. 生成 UUID 存储文件名
5. 保存到 `data/uploads/` 目录
6. 写入数据库记录

**下载流程**:
1. 查询文件记录
2. 检查物理文件存在
3. 增加下载计数
4. 设置 Content-Disposition 头
5. 返回文件内容

#### work_api.go — 作品 API

| 函数 | 路由 | 方法 | 说明 |
|-----|------|------|------|
| `ListWorks` | `/api/works` | GET | 获取作品列表（含媒体） |
| `ListPublicWorks` | `/api/works/public` | GET | 获取公开作品列表 |
| `CreateWork` | `/api/works` | POST | 创建作品 |
| `UpdateWork` | `/api/works/{id}` | PUT | 更新作品 |
| `DeleteWork` | `/api/works/{id}` | DELETE | 删除作品 |
| `AddWorkMedia` | `/api/works/{id}/media` | POST | 添加作品媒体 |
| `DeleteWorkMedia` | `/api/works/media/{id}` | DELETE | 删除作品媒体 |

#### chat_api.go — AI 对话 API

| 函数 | 路由 | 方法 | 说明 |
|-----|------|------|------|
| `ListChatSessions` | `/api/chat/sessions` | GET | 获取对话会话列表 |
| `CreateChatSession` | `/api/chat/sessions` | POST | 创建新对话会话 |
| `DeleteChatSession` | `/api/chat/sessions/{id}` | DELETE | 删除对话会话 |
| `UpdateChatSession` | `/api/chat/sessions/{id}` | PUT | 更新对话会话 |
| `GetChatMessages` | `/api/chat/sessions/{id}/messages` | GET | 获取会话消息 |
| `SendChatMessage` | `/api/chat/send` | POST | 发送消息（SSE 流式响应） |

**SendChatMessage 流程**:
1. 检查 AI API Key 配置
2. 自动创建会话（如果 session_id 为 0）
3. 保存用户消息到数据库
4. 获取最近 20 条消息作为上下文
5. 构建 OpenAI 兼容 API 请求
6. 设置 SSE 响应头
7. 流式读取 AI 响应并转发给客户端
8. 保存 AI 完整回复到数据库
9. 自动生成会话标题（首次对话）

**SSE 事件类型**:
- `session_created` — 自动创建的会话信息
- `start` — 开始接收 AI 响应
- `delta` — 增量文本内容
- `error` — 错误信息
- `done` — 响应完成

---

### 5.3 middleware/ — 中间件

#### auth.go — 认证中间件

| 函数 | 说明 |
|-----|------|
| `CheckAuth(store *sessions.CookieStore) func(http.Handler) http.Handler` | 返回认证中间件，验证 Session 中的 `authenticated` 字段 |

**行为**:
- API 请求 (`/api/*`) → 返回 `401 {"ok":false,"message":"未登录"}`
- 页面请求 → 重定向到 `/admin/login`

---

### 5.4 templates/ — 模板层

| 模板 | 说明 |
|-----|------|
| `index.html` | 公开主页：头像、名称、简介、链接列表、作品网格、文件下载、小组件（时钟/日历/天气） |
| `admin.html` | 管理后台：Tab 导航（链接/文件/作品/设置），CRUD 操作界面 |
| `login.html` | 登录页：密码输入表单 |
| `chat.html` | AI 对话页：侧边栏会话列表 + 主对话区 + 流式消息 |

### 5.5 static/ — 静态资源

| 文件 | 说明 |
|-----|------|
| `style.css` | 全局样式（约 1983 行），包含亮色/暗色主题 CSS 变量、毛玻璃效果、响应式布局 |
| `avatar.jpg` | 默认头像 |
| `bg.jpg` | 背景图片 |
| `logo.jpg` | Logo 图片 |

---

## 六、路由总览

### 6.1 公开路由（无需认证）

| 方法 | 路径 | Handler | 说明 |
|-----|------|---------|------|
| GET | `/` | `IndexPage` | 公开主页 |
| GET | `/admin/login` | `LoginPage` | 登录页 |
| POST | `/api/auth/login` | `Login` | 登录接口 |
| POST | `/api/auth/logout` | `Logout` | 退出接口 |
| GET | `/static/*` | FileServer | 静态文件 |
| GET | `/api/categories` | `ListCategories` | 分类列表 |
| GET | `/api/files/public` | `ListPublicFiles` | 公开文件 |
| GET | `/api/files/download/{id}` | `DownloadFile` | 文件下载 |
| GET | `/api/works/public` | `ListPublicWorks` | 公开作品 |

### 6.2 受保护路由（需认证）

| 方法 | 路径 | Handler | 说明 |
|-----|------|---------|------|
| GET | `/admin` | `AdminPage` | 管理后台 |
| GET | `/chat` | `ChatPage` | AI 对话页 |
| GET/POST/PUT/DELETE | `/api/links/*` | 链接 CRUD | 链接管理 |
| GET/PUT | `/api/settings` | 设置管理 | 系统设置 |
| POST/PUT/DELETE | `/api/categories/*` | 分类管理 | 分类 CRUD |
| GET/POST/DELETE | `/api/files/*` | 文件管理 | 文件上传/删除 |
| GET/POST/PUT/DELETE | `/api/works/*` | 作品管理 | 作品 CRUD + 媒体 |
| GET/POST/PUT/DELETE | `/api/chat/*` | 对话管理 | AI 对话全流程 |

---

## 七、依赖关系

### 7.1 直接依赖

| 依赖 | 版本 | 说明 |
|-----|------|------|
| `github.com/gorilla/sessions` | v1.4.0 | Session 管理（Cookie Store） |
| `golang.org/x/crypto` | v0.52.0 | bcrypt 密码加密 |
| `modernc.org/sqlite` | v1.51.0 | 纯 Go SQLite 驱动（无需 CGO） |

### 7.2 间接依赖

| 依赖 | 说明 |
|-----|------|
| `github.com/google/uuid` | UUID 生成（文件存储命名） |
| `github.com/dustin/go-humanize` | 人类可读格式 |
| `github.com/gorilla/securecookie` | Cookie 安全编码 |
| `modernc.org/libc` | SQLite 依赖的 C 库模拟 |
| `modernc.org/mathutil` | 数学工具 |
| `modernc.org/memory` | 内存管理 |

### 7.3 模块依赖图

```
main.go
  ├── handlers/
  │     ├── models/ (数据操作)
  │     └── gorilla/sessions (Session)
  ├── models/
  │     ├── modernc.org/sqlite (数据库)
  │     ├── golang.org/x/crypto (bcrypt)
  │     └── google/uuid (文件命名)
  └── middleware/
        └── gorilla/sessions (认证)
```

---

## 八、项目运行方式

### 8.1 Docker Compose 部署（推荐）

```bash
docker-compose up -d
```

服务将在 `http://localhost:3001` 启动（映射容器内 8080 端口）。

### 8.2 本地开发运行

```bash
# 构建
go build -o linktree.exe .

# 运行
./linktree.exe
```

服务将在 `http://localhost:8080` 启动。

### 8.3 环境变量

| 变量 | 默认值 | 说明 |
|-----|-------|------|
| `PORT` | `8080` | 服务监听端口 |
| `DATA_DIR` | `data` | 数据存储目录（数据库+上传文件） |
| `ADMIN_PASSWORD` | 随机生成 | 管理员初始密码（仅首次运行生效） |
| `SESSION_SECRET` | 内置默认值 | Session 加密密钥（生产环境务必修改） |
| `MAX_FILE_SIZE` | `209715200` | 最大上传文件大小 (200MB) |
| `TZ` | — | 时区设置 |

### 8.4 数据目录结构

```
data/
├── linktree.db          # SQLite 数据库
└── uploads/             # 上传文件存储
    ├── xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx.jpg
    └── ...
```

### 8.5 Docker 构建说明

- **多阶段构建**: `golang:1.23-alpine` 编译 → `alpine:3.21` 运行
- **CGO_ENABLED=0**: 纯 Go 编译，无需 C 工具链
- **数据卷**: `/app/data` 挂载到宿主机 `./data`
- **镜像源**: 使用阿里云镜像加速

---

## 九、安全设计

| 安全措施 | 说明 |
|---------|------|
| 密码加密 | bcrypt 哈希存储，不保存明文 |
| Session 认证 | HttpOnly Cookie，24 小时过期 |
| 文件名安全 | UUID 重命名存储，防止路径穿越 |
| MIME 白名单 | 仅允许指定类型文件上传 |
| API Key 脱敏 | 返回前端时掩码显示 (`sk-xxxx****xxxx`) |
| 路由保护 | 管理功能全部经过认证中间件 |
| 外键约束 | SQLite 启用外键，级联删除 |

---

## 十、前端设计

### 10.1 UI 风格

- **ZYYO 风格**毛玻璃卡片设计
- CSS 变量驱动的亮色/暗色主题
- 响应式布局（移动端适配）
- Q 弹动画效果（cubic-bezier 缓动）

### 10.2 主要页面

| 页面 | 特点 |
|-----|------|
| 公开主页 | 居中卡片布局 + 右侧小组件（时钟/日历/天气） |
| 管理后台 | Tab 导航（链接/文件/作品/设置） |
| AI 对话 | 侧边栏会话列表 + 主对话区 + SSE 流式输出 |
| 登录页 | 居中表单卡片 |

### 10.3 小组件

| 组件 | 功能 |
|-----|------|
| 时钟 | 实时数字时钟 + 点击展开模拟表盘 |
| 日历 | 当前日期 + 点击展开月历 |
| 天气 | 天气信息 + 点击展开详情 |
| 主题切换 | 亮色/暗色/自动模式切换 |
