# 第二轮优化记录：2026-10-06（性能 / SEO / 可访问性）

承接 `HANDOVER-2026-10-04.md`。本轮按用户选定的 **P0（性能+SEO基础）+ P3（体验与可访问性）** 执行；
**P2（作品集内容重写）因用户暂无素材而暂缓**，见文末。

## 一、性能（实测，本地 :8084 + 线上 :3001 双环境）

| 项 | 优化前 | 优化后 |
|---|---|---|
| JS 传输 | 390,623 B（无压缩） | **122,557 B**（gzip，↓69%） |
| CSS 传输 | 25,958 B | **4,965 B**（↓81%） |
| 首页 HTML | 1,942 B | **854 B**（↓56%） |
| 第三方字体请求 | fonts.googleapis.com 208ms + fonts.gstatic.com 226/86ms | **0 个** |
| 字体缓存 | 每次协商 | `public, max-age=31536000, immutable` |
| 构建产物 | 强缓存但无 ETag | 强缓存 + `ETag` → 304 |
| 生产首页加载 | — | **load 58ms** |

### 实现要点

- **gzip 中间件**（`main.go`）：自研 `gzipResponseWriter`，`sync.Pool` 复用 writer；
  只压缩 `compressibleType()` 命中的类型，**图片/字体/大文件下载自动跳过**；
  `Vary: Accept-Encoding` 无条件设置。**SSE 不受影响**（`/api/ask` 实测仍是
  `text/event-stream`、无 `Content-Encoding`、`data:` 事件逐个到达）。
- **字体自托管**：Geist / Geist Mono（SIL OFL 1.1）。保留 Google 的 48 个 `@font-face`
  子集逻辑，只把 URL 换成 `/static/fonts/*`；实际文件仅 **11 个 / 147KB**，
  浏览器按 `unicode-range` 只取命中页面的子集（中文页面只取拉丁子集，其余走系统字体，
  Geist 本身不含 CJK 字形）。授权说明在 `static/fonts/OFL.txt`。
  生成脚本思路保留在 `.dsh-tmp/fetch-fonts.ps1`、`gen-fonts-css.ps1`（临时目录，不在仓库）。
- **静态资源**：`serveStaticFile()` 统一处理 ETag / 304 / Range；
  `/assets/` 与 `/static/fonts/` 走一年 immutable，其余 `no-cache`。

## 二、SEO / 分享

优化前爬虫只拿到 **773 字节空壳**、0 个 og 标签、`/robots.txt` 与 `/sitemap.xml` 都返回首页 HTML、
任意未知路径返回 200（软 404）。

- **服务端按路由注入 meta**（`renderIndex()`，不依赖 JS）：
  `/` → `啊芃 · 全栈 / AI 开发者`、`/works` → `作品 - 啊芃`、
  `/work/2` → `LinkTree 个人主页 - 啊芃`（读数据库真实标题）。
- 补齐 **Open Graph + Twitter 卡片**（微信/QQ/脉脉分享不再是空白卡）。
- 新增 **`/robots.txt`**（Disallow `/admin`、`/chat`、`/api/`）与 **`/sitemap.xml`**（10 条 URL，作品来自数据库）。
- **真 404**：`/nope`、`/work/9999`、`/foo/bar` 均返回 404；body 仍是 SPA 外壳并带
  `noindex,follow`，用户看到设计过的"这个页面不存在"页（不再是 Go 默认的 `404 page not found` 纯文本）。
- 路由级 `<title>`（`main.js` 的 `router.afterEach`），标签页/历史记录不再全叫"啊芃"。
- 作品不存在时（`/work/9999`）服务端标题为 `作品不存在 - 啊芃` + 404。

## 三、体验与可访问性（P3）

- **skip link**（`.skip-link`，聚焦后出现在 16px 处）跳到 `#main`。
- 显式补齐 `a / button / [tabindex]` 的 `:focus-visible` 焦点环（WCAG 2.4.7）。
- **筛选 chips 改为读后台分类**（`api.categories()`，带数量角标）：
  线上实测 `全部 3 / Web 应用 2 / AI 应用 1`，点击筛选实时生效。
  此前是硬编码 `web/ai/tool`，后台新增分类前台看不到。
- **移动端首屏**：hero 542px → **449px**，精选作品进入首屏，375px 无横向溢出。
- 精选作品 2 → 3 张（`slice(0,3)`，桌面三栏 / 平板两栏）。
- 新增 `NotFoundPage.vue`，兜底路由不再是 `redirect: '/'`（避免未知路径被静默吞成首页）。

## 四、⚠️ 事件：后台密码被改回 `admine`

**现象**：10-06 做回归时新密码登录失败（401）。
**取证**：用 bcrypt 直接比对线上哈希（不试登录、避开限流）：

```
线上 admin_password_hash = $2a$10$L3fOz2GKVlTcOZmh20Si7uKCAzvJnqKP4iTsxO4.WKgcSWaKuvPf.
  Bsh4e6JQye1OQaaNKzjEkqvo（10-04 轮换的）  match=false
  admine                                      match=true   ← 当时有效密码
  ChangeMe123!                                match=false
```

**结论**：10-04 轮换之后，密码被改回了排查过程中公开过的弱密码 `admine`。
数据库是 **WAL 模式**，主库文件 mtime 停在 10-04 20:53，之后的写入都在 `-wal` 里，
**无法从时间戳定位改动时刻与来源**（可能是用户在后台点过"修改密码"，也可能是别的原因）。

**处置**：
1. 重新轮换，`.env` 与库内哈希重新一致（已实测 200；`admine` 现在 401）。
2. `models/settings.go:ChangePassword` 增加审计字段，任何改密都会留痕：
   `admin_password_changed_at`（时间）、`admin_password_changed_by`（`env_rotation` / `admin_panel`）。
   当前值：`2026-10-06 14:38:02` / `env_rotation`。

**教训**：任何凭据变更都必须落库留痕，否则"密码为什么不对"无从查证。

## 五、仍未完成

1. **P2 作品集内容重写**（用户已选，但暂无素材）：计划给 `works` 表加
   `tech_stack / role / highlights / repo_url / demo_url` 字段 + 后台输入框 + 详情页结构化渲染；
   或等用户口述素材后由 AI 组织文案。**不编造用户经历**。
2. 公开列表第 1 个作品标题是「测试」、描述为空（HR 第一眼就看到它）；
   关于页「20+ 项目作品」与实际 4 个不符。
3. 简历 PDF 仍未上传，`/api/public/resume` 返回 `available:false`。
4. 三个作品的真实封面图（3/4 显示渐变占位）。
5. `ai_api_key` 仍明文存于 SQLite `settings`（`sk-27966be…`，已泄露），需用户换 key。
6. git 历史仍留有已作废的旧密钥明文（`e0e9580` 起）；轮换后已无害，彻底清除需重写历史。
7. 部署升级未动：**HTTPS / 反代去掉 3001 端口**、data 定期备份。
   用户目标是**国内岗位**，`http://62.234.92.232:3001` 这种地址放简历上是减分项。
8. 百度对纯客户端渲染的索引能力有限；如需强索引要做预渲染（prerender）或 SSR，成本较高。

## 六、本轮提交

| 提交 | 内容 |
|---|---|
| `9475501` | perf+seo: gzip、immutable 缓存、字体自托管、服务端 meta、真 404、a11y |
