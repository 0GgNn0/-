# LinkTree 前台重构报告（Vue3 + Naive UI + Vite SPA）

## 版本信息
- **日期**: 2026-10-01
- **状态**: ✅ 已部署线上并全量验证通过（本地 + 生产环境）
- **Git**: 本地已改、**尚未提交**（上一次提交为 `ee1aa3a`）

---

## 一、改动概述

将前台 7 个公开页面从 Go HTML 模板重构为 **Vue 3 + Naive UI + Vite** 单页应用（SPA），
后台（admin/login/chat）保留原有 Go 模板，架构不变、API 不变。

```
浏览器
 ├─ /admin/*  /chat        → Go 模板（原样，登录态管理）
 ├─ /api/*                 → Go JSON API（新增 2 个公开接口）
 └─ /  /works  /about ...  → SPA（web/dist，Go spaHandler 兜底回 index.html）
```

---

## 二、新增文件清单

### Go 后端
| 文件 | 说明 |
|-----|------|
| `handlers/public_api.go` | 公开接口 `GET /api/public/profile`、`GET /api/public/links` |
| `.dockerignore` | 排除 data/、node_modules/、dist/、linktree.exe、日志等 |

### 前端工程 `web/`（全新目录）
| 文件 | 说明 |
|-----|------|
| `package.json` / `package-lock.json` | Vue3 + vue-router + naive-ui + Vite6 |
| `vite.config.js` | 开发代理 `/api`、`/static` → `:8080` |
| `index.html` | Geist / Geist Mono 字体引入 |
| `src/main.js` | 7 条路由 + 兜底 `/:pathMatch(.*)* → /` |
| `src/App.vue` | 顶栏导航 + 主题切换 + footer + n-config-provider |
| `src/api.js` | fetch 封装 + `/api/ask` SSE 流式解析 |
| `src/theme.js` | 三态主题（light/dark/auto），localStorage `linktree-theme` |
| `src/profile.js` | 共享 profile 状态（display_name/bio） |
| `src/style.css` | 2026 视觉系统（噪点、聚光边框、reveal 交错入场、mono kicker） |
| `src/pages/HomePage.vue` | 92px 超大标题 + Open to Work 徽章 + 精选作品 |
| `src/pages/WorksPage.vue` | 作品网格 + 分类筛选 chips |
| `src/pages/WorkDetailPage.vue` | 作品详情（图/视频/音频 + 外链按钮） |
| `src/pages/AboutPage.vue` | 简介 + 数据统计 + 技术栈 |
| `src/pages/ExperiencePage.vue` | 经历时间线（读 `/api/experiences`） |
| `src/pages/ContactPage.vue` | 联系方式（读 `/api/public/links`） |
| `src/pages/AskPage.vue` | AI 问答，SSE 流式输出 + 推荐问题 |

---

## 三、修改文件清单

| 文件 | 修改内容 |
|-----|---------|
| `main.go` | 移除 7 个 public 页面模板路由；新增 `spaHandler("web/dist")`；注册 `/api/public/profile\|links` |
| `Dockerfile` | 改为三阶段构建：`node:20-alpine` 构建 dist → `golang` 编译 → `alpine` 运行；`COPY --from=webbuilder /web/dist ./web/dist` |
| `web/src/App.vue` | **Bug 修复**：页脚"管理后台"由 `<router-link>` 改为 `<a href>`（详见第五节） |

---

## 四、路由与接口对照

### 前台路由（SPA）
| 路径 | 页面 |
|-----|------|
| `/` | 首页 |
| `/works` | 作品列表 |
| `/work/:id` | 作品详情 |
| `/about` | 关于我 |
| `/experience` | 经历 |
| `/contact` | 联系 |
| `/ask` | AI 问答（SSE 流式） |

### 公开 API
| 接口 | 说明 |
|-----|------|
| `GET /api/public/profile` | display_name / bio / theme |
| `GET /api/public/links` | 对外链接（联系页） |
| `GET /api/works/public` | 作品 + 媒体（原已有） |
| `GET /api/experiences` | 经历时间线（原已有） |
| `POST /api/ask` | AI 问答 SSE（原已有） |

未匹配的 `/api/*` GET 请求统一返回 **404 JSON**（不会回退成 index.html）。

---

## 五、修复的 Bug

### 页脚"管理后台"点击无效
- **现象**：点击页脚"管理后台"停留在首页（或从内页被弹回首页）。
- **根因**：`App.vue` 使用 `<router-link to="/admin/login">`，被 `main.js` 兜底路由
  `/:pathMatch(.*)* → redirect '/'` 劫持，Vue Router 直接重定向回首页。
- **修复**：改为普通 `<a href="/admin/login">`，浏览器整页跳转到 Go 模板登录页。
- **状态**：本地 + 线上点击后均正确到达 `/admin/login` ✅

---

## 六、部署方式

```bash
# 本地打包上传（排除 data、node_modules、dist、可执行文件、日志）
tar -czf linktree-deploy.tar.gz --exclude=./data --exclude=./web/node_modules \
    --exclude=./web/dist --exclude=./linktree.exe --exclude=./server.log .
scp linktree-deploy.tar.gz root@62.234.92.232:/root/

# 服务器解压 + 重建
ssh root@62.234.92.232
cd /opt/linktree && tar -xzf /root/linktree-deploy.tar.gz -C /opt/linktree
docker compose up -d --build linktree
```

- 线上地址：http://62.234.92.232:3001/
- 后台入口：首页页脚"管理后台" → `/admin/login`（密码 `admine`）
- Docker 构建内置 `npm` registry 为 npmmirror、`GOPROXY` 为 goproxy.cn（国内加速）

---

## 七、验证报告（本地 :8080 + 线上 :3001 双环境）

| 验证项 | 结果 |
|-------|------|
| 7 个 SPA 路由（直接访问 + 刷新） | ✅ 全部 200，回退 index.html |
| 首页 DOM（hero 92px、Open to Work 徽章、精选卡、5 顶栏链接） | ✅ |
| 作品页（卡片 + 分类 chips，本地 3 / 线上 4 条） | ✅ |
| 作品详情（分类标签、描述、外链按钮、返回链接） | ✅ |
| 关于页（3 数据块 + 8 技术栈 chips） | ✅ |
| 经历页（时间线 3 条，读 API） | ✅ |
| 联系页（3 条外链，读 API） | ✅ |
| ask 页 SSE 流式问答（本地 + 线上各实测 1 次） | ✅ |
| 未登录访问 `/admin` → 302 → `/admin/login` | ✅ |
| 后台登录 `admine` → `/admin` 6 面板 | ✅ |
| 后台 6 个桌面 tab 切换 + 数据加载 | ✅ |
| 后台"经历添加"表单、"新建作品"弹窗 | ✅ |
| 后台移动端 375px（侧栏隐藏、6 移动 tab、无横向溢出） | ✅ |
| 主题切换（light/dark、localStorage 持久化） | ✅ |
| 全链路：首页 → 管理后台 → 登录 → 后台 → 返回首页 | ✅ |
| 网络请求（全部 200/304，零 404） | ✅ |
| 浏览器控制台（零报错） | ✅ |

---

## 八、待办 / 注意事项

1. **Git 未提交**：本轮改动（`web/`、`main.go`、`handlers/public_api.go`、
   `Dockerfile`、`.dockerignore`、`templates` 相关、`static/style.css`）待提交推送。
   仓库：`https://github.com/0GgNn0/-.git`（master，上一提交 `ee1aa3a`）。
2. **字体依赖 Google Fonts**：Geist 字体来自 fonts.googleapis.com，
   国内访问可能偏慢（有 fallback 字体兜底），后续可考虑自托管。
3. **`web/dist` 不入库**：由 Docker 构建时生成；本地预览用 `npm run build` 或
   `npm run dev`（dev 模式代理到 :8080）。
4. 后台密码 `admine`（设置在 SQLite `settings.admin_password_hash`，
   环境变量 `ADMIN_PASSWORD` 仅首次生效）。
