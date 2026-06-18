# LinkTree v2 开发文档

## 版本信息
- **版本**: v2.0
- **新增功能**: 文件存储系统 + 作品展示系统
- **编译状态**: ✅ 通过

---

## 一、新增文件清单

| 文件 | 说明 | 行数 |
|-----|------|------|
| `models/category.go` | 分类 CRUD 模型 | 72 |
| `models/file.go` | 文件模型 + 上传逻辑 + 类型白名单 | 190 |
| `models/work.go` | 作品模型 + 媒体关联 | 168 |
| `handlers/category_api.go` | 分类 API handler | 92 |
| `handlers/file_api.go` | 文件上传/下载/删除 API | 195 |
| `handlers/work_api.go` | 作品管理 API | 172 |

---

## 二、修改文件清单

| 文件 | 修改内容 |
|-----|---------|
| `models/db.go` | 新增 4 张表建表语句（categories, files, works, work_media） |
| `main.go` | 新增 18 条路由（分类、文件、作品相关 API） |
| `handlers/page.go` | 新增 Works/Files 数据传递到模板 |
| `templates/admin.html` | 重写为 Tab 导航结构（链接/文件/作品/设置） |
| `templates/index.html` | 新增作品展示网格 + 文件下载列表 + 作品详情弹窗 |
| `static/style.css` | 新增约 300 行样式（Tab、作品卡片、文件列表、弹窗等） |
| `Dockerfile` | 新增 `mkdir -p /app/data/uploads` |
| `docker-compose.yml` | 新增 DATA_DIR、MAX_FILE_SIZE 环境变量 |

---

## 三、数据库新增表

### categories（分类表）
```sql
CREATE TABLE categories (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    name       TEXT NOT NULL,
    type       TEXT NOT NULL CHECK(type IN ('file', 'work')),
    sort_order REAL NOT NULL DEFAULT 0,
    created_at TEXT DEFAULT (datetime('now','localtime'))
);
```

### files（文件表）
```sql
CREATE TABLE files (
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    category_id    INTEGER,
    original_name  TEXT NOT NULL,
    stored_name    TEXT NOT NULL,
    mime_type      TEXT NOT NULL,
    size           INTEGER NOT NULL,
    is_public      INTEGER NOT NULL DEFAULT 1,
    download_count INTEGER NOT NULL DEFAULT 0,
    sort_order     REAL NOT NULL DEFAULT 0,
    created_at     TEXT DEFAULT (datetime('now','localtime')),
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE SET NULL
);
```

### works（作品表）
```sql
CREATE TABLE works (
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    category_id  INTEGER,
    title        TEXT NOT NULL,
    description  TEXT DEFAULT '',
    cover_url    TEXT DEFAULT '',
    content_type TEXT NOT NULL CHECK(content_type IN ('image', 'video', 'mixed')),
    external_url TEXT DEFAULT '',
    sort_order   REAL NOT NULL DEFAULT 0,
    created_at   TEXT DEFAULT (datetime('now','localtime')),
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE SET NULL
);
```

### work_media（作品媒体表）
```sql
CREATE TABLE work_media (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    work_id    INTEGER NOT NULL,
    file_id    INTEGER,
    media_url  TEXT NOT NULL,
    media_type TEXT NOT NULL CHECK(media_type IN ('image', 'video')),
    sort_order REAL NOT NULL DEFAULT 0,
    FOREIGN KEY (work_id) REFERENCES works(id) ON DELETE CASCADE,
    FOREIGN KEY (file_id) REFERENCES files(id) ON DELETE SET NULL
);
```

---

## 四、API 路由总览

### 公开路由（无需登录）
| 方法 | 路径 | 说明 |
|-----|------|------|
| `GET` | `/` | 公开主页（含链接、作品、文件） |
| `GET` | `/api/files/public` | 获取公开文件列表 |
| `GET` | `/api/files/download/{id}` | 下载/预览文件 |
| `GET` | `/api/works/public` | 获取公开作品列表 |

### 管理路由（需登录）
| 方法 | 路径 | 说明 |
|-----|------|------|
| `GET` | `/api/categories?type=` | 获取分类列表 |
| `POST` | `/api/categories` | 创建分类 |
| `PUT` | `/api/categories/{id}` | 更新分类 |
| `DELETE` | `/api/categories/{id}` | 删除分类 |
| `GET` | `/api/files` | 获取文件列表 |
| `POST` | `/api/files/upload` | 上传文件（multipart） |
| `DELETE` | `/api/files/{id}` | 删除文件 |
| `GET` | `/api/works` | 获取作品列表 |
| `POST` | `/api/works` | 创建作品 |
| `PUT` | `/api/works/{id}` | 更新作品 |
| `DELETE` | `/api/works/{id}` | 删除作品 |
| `POST` | `/api/works/{id}/media` | 添加作品媒体 |
| `DELETE` | `/api/works/media/{id}` | 删除作品媒体 |

---

## 五、文件上传规格

| 项目 | 规格 |
|-----|------|
| 最大文件大小 | 200MB |
| 存储位置 | `data/uploads/` 目录 |
| 文件命名 | UUID + 原始扩展名（防路径穿越） |
| MIME 检查 | 白名单验证（图片/文档/压缩包/视频） |
| 公开访问 | 默认公开，无需密码 |

**支持的文件类型：**
- 图片：jpeg, png, gif, webp, svg
- 文档：pdf, txt, doc, docx, xls, xlsx
- 压缩：zip, rar, gz
- 视频：mp4, webm, ogg

---

## 六、Admin 后台新功能

### Tab 导航
```
[链接] [文件] [作品] [设置]
```

### 文件管理 Tab
- 分类筛选下拉框
- 新建分类按钮
- 多文件上传（支持同时选多个）
- 文件列表：文件名、大小、分类、下载次数、操作

### 作品管理 Tab
- 分类筛选下拉框
- 新建分类按钮
- 新建作品弹窗：标题、描述、封面图、内容类型、外链、分类
- 作品卡片网格展示
- 点击卡片编辑：修改信息 + 管理媒体列表

---

## 七、公开主页新功能

### 作品展示区
- 2列网格布局
- 封面图 + 标题 + 描述
- 点击打开详情弹窗（显示封面、描述、媒体列表、外链）

### 文件下载区
- 列表布局
- 文件图标（根据类型自动匹配）
- 文件名 + 大小
- 点击直接下载

---

## 八、部署说明

### 方式一：Docker Compose（推荐）
```bash
docker-compose up -d
```

### 方式二：本地运行
```bash
go build -o linktree.exe .
./linktree.exe
```

### 环境变量
| 变量 | 默认值 | 说明 |
|-----|-------|------|
| `PORT` | 8080 | 服务端口 |
| `DATA_DIR` | data | 数据目录 |
| `ADMIN_PASSWORD` | 随机生成 | 管理密码 |
| `SESSION_SECRET` | 默认值 | Session 密钥 |
| `MAX_FILE_SIZE` | 209715200 | 最大文件大小（200MB） |

---

## 九、文件结构变更图

```
linktree/
├── main.go                    # +18条路由
├── models/
│   ├── db.go                  # +4张新表
│   ├── link.go                # (不变)
│   ├── settings.go            # (不变)
│   ├── category.go            # [新增] 分类 CRUD
│   ├── file.go                # [新增] 文件管理
│   └── work.go                # [新增] 作品管理
├── handlers/
│   ├── api.go                 # (不变)
│   ├── page.go                # +传递 Works/Files 数据
│   ├── category_api.go        # [新增] 分类 API
│   ├── file_api.go            # [新增] 文件 API
│   └── work_api.go            # [新增] 作品 API
├── templates/
│   ├── admin.html             # 重写：Tab导航+文件/作品管理
│   ├── index.html             # +作品展示+文件下载区
│   └── login.html             # (不变)
├── static/
│   ├── style.css              # +300行新样式
│   └── uploads/               # 文件上传目录
├── Dockerfile                 # +创建uploads目录
└── docker-compose.yml         # +环境变量
```

---

## 十、后续可选增强

| 功能 | 说明 | 优先级 |
|-----|------|--------|
| 图片缩略图 | 上传时自动生成缩略图 | 中 |
| 拖拽排序 | 文件/作品拖拽排序 | 低 |
| 批量删除 | 一次删除多个文件 | 低 |
| 存储统计 | 显示已用空间/文件数 | 低 |
| 文件预览 | 图片/PDF 在线预览 | 中 |
| 外部存储 | 支持 S3/OSS | 低 |

---

**文档完成时间**: 2026-06-15
**代码状态**: ✅ 编译通过
