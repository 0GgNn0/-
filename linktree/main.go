package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"linktree/handlers"
	"linktree/middleware"
	"linktree/models"

	"github.com/gorilla/sessions"
)

// ---------- 压缩 ----------

var gzipPool = sync.Pool{New: func() any { return gzip.NewWriter(io.Discard) }}

// compressibleType 判断该 Content-Type 是否值得压缩（SSE 的 text/event-stream 不在其列）
func compressibleType(ct string) bool {
	ct = strings.ToLower(strings.TrimSpace(strings.Split(ct, ";")[0]))
	switch ct {
	case "text/html", "text/css", "text/plain", "text/xml", "text/javascript",
		"application/javascript", "application/x-javascript", "application/json",
		"application/xml", "application/rss+xml", "image/svg+xml":
		return true
	}
	return false
}

// gzipResponseWriter 只压缩 compressibleType 命中的响应，跳过大文件下载与 SSE
type gzipResponseWriter struct {
	http.ResponseWriter
	gz          *gzip.Writer
	compress    bool
	wroteHeader bool
}

func (w *gzipResponseWriter) WriteHeader(code int) {
	if !w.wroteHeader {
		w.wroteHeader = true
		h := w.Header()
		if compressibleType(h.Get("Content-Type")) && h.Get("Content-Encoding") == "" {
			w.compress = true
			h.Del("Content-Length") // 长度在压缩后变化，交给 chunked
			h.Set("Content-Encoding", "gzip")
		}
		// 无论是否压缩都要 Vary，避免缓存把两种编码混用
		h.Add("Vary", "Accept-Encoding")
	}
	w.ResponseWriter.WriteHeader(code)
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	if w.compress {
		if w.gz == nil {
			w.gz = gzipPool.Get().(*gzip.Writer)
			w.gz.Reset(w.ResponseWriter)
		}
		return w.gz.Write(b)
	}
	return w.ResponseWriter.Write(b)
}

func (w *gzipResponseWriter) Flush() {
	if w.gz != nil {
		w.gz.Flush()
	}
	if f, ok := w.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

func (w *gzipResponseWriter) finish() {
	if w.gz != nil {
		w.gz.Close()
		gzipPool.Put(w.gz)
		w.gz = nil
	}
}

// gzipMiddleware 全站压缩；SSE(/api/ask) 与大文件下载因 Content-Type 不匹配自动跳过
func gzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") || r.Method == http.MethodHead {
			next.ServeHTTP(w, r)
			return
		}
		gw := &gzipResponseWriter{ResponseWriter: w}
		defer gw.finish()
		next.ServeHTTP(gw, r)
	})
}

// ---------- 静态资源 ----------

// serveStaticFile 带强缓存 + 条件请求（ETag/304）地发送文件。
// maxAge > 0 下发 immutable 长缓存；maxAge == 0 下发 no-cache（index.html 用）。
func serveStaticFile(w http.ResponseWriter, r *http.Request, fullPath string, maxAge int) {
	st, err := os.Stat(fullPath)
	if err != nil || st.IsDir() {
		http.NotFound(w, r)
		return
	}
	etag := `W/"` + strconv.FormatInt(st.ModTime().UnixNano(), 36) + "-" + strconv.FormatInt(st.Size(), 36) + `"`
	h := w.Header()
	if maxAge > 0 {
		h.Set("Cache-Control", "public, max-age="+strconv.Itoa(maxAge)+", immutable")
	} else {
		h.Set("Cache-Control", "no-cache")
	}
	h.Set("ETag", etag)
	if match := r.Header.Get("If-None-Match"); match != "" && strings.Contains(match, etag) {
		w.WriteHeader(http.StatusNotModified)
		return
	}
	f, err := os.Open(fullPath)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	defer f.Close()
	// ServeContent 负责 Content-Type / Content-Length / Range
	http.ServeContent(w, r, filepath.Base(fullPath), st.ModTime(), f)
}

// spaHandler 服务 Vite 构建产物：真实文件直出；未知路径返回真 404（不再静默回退首页）。
func spaHandler(distDir string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := filepath.Clean(r.URL.Path)
		full := filepath.Join(distDir, p)
		if st, err := os.Stat(full); err == nil && !st.IsDir() {
			if strings.HasPrefix(p, "/assets/") {
				serveStaticFile(w, r, full, 31536000)
			} else {
				serveStaticFile(w, r, full, 0)
			}
			return
		}
		http.NotFound(w, r)
	})
}

// 每个前端路由的 meta：服务端直接注入 index.html，让不执行 JS 的爬虫/分享卡片也能拿到真实标题
var routeMeta = map[string][2]string{
	"/":           {"", "全栈 / AI 开发者个人主页：作品集、经历与 AI 问答"},
	"/works":      {"作品", "我做过的东西：Web 应用、AI 应用与工具类项目"},
	"/about":      {"关于我", "技术栈、个人简介与职业方向"},
	"/experience": {"经历", "工作与项目经历时间线"},
	"/contact":    {"联系", "通过邮箱或社交账号联系我"},
	"/ask":        {"问 AI", "用 AI 问答了解我的技能与项目经历"},
}

// escapeMeta 处理注入 HTML 属性的文本（标题均来自本站数据库，仍做防御性转义）
func escapeMeta(s string) string {
	r := strings.NewReplacer(
		`&`, "&amp;", `<`, "&lt;", `>`, "&gt;", `"`, "&quot;", `'`, "&#39;",
		"\r", " ", "\n", " ",
	)
	return r.Replace(s)
}

// renderIndex 读取 index.html 并替换 meta 占位符
func renderIndex(status int, path string, w http.ResponseWriter) {
	b, err := os.ReadFile(filepath.Join("web", "dist", "index.html"))
	if err != nil {
		http.Error(w, "frontend not built", http.StatusServiceUnavailable)
		return
	}
	site := "啊芃"
	title := site + " · 全栈 / AI 开发者"
	desc := "全栈 / AI 开发者个人主页：作品集、经历与 AI 问答"

	if path == "/ask-ai" { // 兼容旧入口
		path = "/ask"
	}
	if m, ok := routeMeta[path]; ok {
		if m[0] != "" {
			title = m[0] + " - " + site
		}
		desc = m[1]
	} else if strings.HasPrefix(path, "/work/") {
		idStr := strings.TrimPrefix(path, "/work/")
		if id, err := strconv.Atoi(idStr); err == nil {
			if wk, err := models.GetWorkByID(id); err == nil && wk != nil {
				title = wk.Title + " - " + site
				if wk.Description != "" {
					desc = wk.Description
				}
			} else {
				// 不存在的作品 ID：按真 404 返回，避免软 404 被搜索引擎收录
				title = "作品不存在 - " + site
				status = http.StatusNotFound
			}
		}
	}

	base := siteBaseURL()
	ogImage := base + "/static/avatar.jpg"
	html := string(b)
	html = strings.ReplaceAll(html, "{{TITLE}}", escapeMeta(title))
	html = strings.ReplaceAll(html, "{{DESCRIPTION}}", escapeMeta(desc))
	html = strings.ReplaceAll(html, "{{CANONICAL}}", escapeMeta(base+path))
	html = strings.ReplaceAll(html, "{{OG_IMAGE}}", escapeMeta(ogImage))
	if status == http.StatusNotFound {
		// 未知路径：状态码保持 404（搜索引擎不会收录），但仍返回 SPA 外壳让用户看到设计过的 404 页
		html = strings.ReplaceAll(html, `<meta name="theme-color"`,
			`<meta name="robots" content="noindex,follow">`+"\n  "+`<meta name="theme-color"`)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.WriteHeader(status)
	io.WriteString(w, html)
}

// spaIndex 返回注入 meta 后的 index.html：200=已知前端路由，404=未知路径（软 404 修复）
func spaIndex(status int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		renderIndex(status, r.URL.Path, w)
	}
}

// ---------- robots / sitemap ----------

func siteBaseURL() string {
	if v := os.Getenv("SITE_BASE_URL"); v != "" {
		return strings.TrimRight(v, "/")
	}
	return "http://62.234.92.232:3001"
}

func robotsTxt(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=3600")
	fmt.Fprintf(w, "User-agent: *\nAllow: /\nDisallow: /admin\nDisallow: /chat\nDisallow: /api/\n\nSitemap: %s/sitemap.xml\n", siteBaseURL())
}

func sitemapXML(w http.ResponseWriter, r *http.Request) {
	base := siteBaseURL()
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=3600")
	io.WriteString(w, `<?xml version="1.0" encoding="UTF-8"?>`+"\n")
	io.WriteString(w, `<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`+"\n")
	for _, p := range []string{"/", "/works", "/about", "/experience", "/contact", "/ask"} {
		fmt.Fprintf(w, "  <url><loc>%s%s</loc><changefreq>weekly</changefreq></url>\n", base, p)
	}
	if works, err := models.GetPublicWorks(); err == nil {
		for _, wk := range works {
			fmt.Fprintf(w, "  <url><loc>%s/work/%d</loc><changefreq>monthly</changefreq></url>\n", base, wk.ID)
		}
	}
	io.WriteString(w, "</urlset>\n")
}

// noCache 供 /static/ 下非字体资源使用（字体走长缓存，单独注册）
func noCache(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache, must-revalidate")
		next.ServeHTTP(w, r)
	})
}

func main() {
	// --rotate-admin-password：显式把 ADMIN_PASSWORD 重置为库中密码（用于凭据轮换）。
	// 默认不执行：环境变量只在数据库首次初始化时生效，避免每次重启都改密码。
	rotateAdminPassword := false
	for _, a := range os.Args[1:] {
		switch a {
		case "--rotate-admin-password":
			rotateAdminPassword = true
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	dataDir := os.Getenv("DATA_DIR")
	if dataDir == "" {
		dataDir = "data"
	}
	dbPath := filepath.Join(dataDir, "linktree.db")

	// Ensure data dir exists
	os.MkdirAll(dataDir, 0755)

	// Init DB
	if err := models.InitDB(dbPath); err != nil {
		log.Fatalf("Failed to init database: %v", err)
	}

	// Seed default data
	models.SeedDefaultLinks()
	models.SeedDefaultExperiences()

	// Init admin password (first run only)
	generatedPassword, err := models.InitAdminPassword(os.Getenv("ADMIN_PASSWORD"))
	if err != nil {
		log.Printf("Warning: failed to initialize admin password: %v", err)
	}
	if generatedPassword != "" {
		fmt.Println("========================================")
		fmt.Printf("  初始管理密码: %s\n", generatedPassword)
		fmt.Println("  请首次登录后立即修改！")
		fmt.Println("========================================")
	}

	if rotateAdminPassword {
		msg, err := models.RotateAdminPasswordFromEnv(os.Getenv("ADMIN_PASSWORD"))
		if err != nil {
			log.Fatalf("管理密码轮换失败: %v", err)
		}
		log.Printf("[admin-password-rotate] %s", msg)
	}

	// Init templates
	if err := handlers.InitTemplates(); err != nil {
		log.Fatalf("Failed to parse templates: %v", err)
	}

	// Init session store
	sessionKey := os.Getenv("SESSION_SECRET")
	if sessionKey == "" {
		log.Fatal("SESSION_SECRET 未设置，拒绝启动（会话密钥缺失会导致认证可伪造）")
	}
	store := sessions.NewCookieStore([]byte(sessionKey))
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400, // 24 hours
		HttpOnly: true,
		Secure:   os.Getenv("COOKIE_SECURE") == "1", // HTTPS 部署时设为 1
		SameSite: http.SameSiteLaxMode,
	}
	handlers.SetStore(store)

	// Build routes
	mux := http.NewServeMux()

	// Public routes (SPA frontend served by spaHandler at "GET /")
	mux.HandleFunc("GET /api/public/profile", handlers.PublicProfile)
	mux.HandleFunc("GET /api/public/links", handlers.PublicLinks)
	mux.HandleFunc("GET /api/public/resume", handlers.PublicResume)
	mux.HandleFunc("POST /api/ask", handlers.AskAI)
	mux.HandleFunc("GET /admin/login", handlers.LoginPage)
	mux.HandleFunc("POST /api/auth/login", handlers.Login)
	mux.HandleFunc("POST /api/auth/logout", handlers.Logout)

	// Static files：字体是不可变资源给长缓存，其余保持 no-cache
	mux.Handle("GET /static/fonts/", http.StripPrefix("/static/fonts/",
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			serveStaticFile(w, r, filepath.Join("static", "fonts", filepath.Base(r.URL.Path)), 31536000)
		})))
	mux.Handle("GET /static/", http.StripPrefix("/static/", noCache(http.FileServer(http.Dir("static")))))

	// SEO：robots / sitemap
	mux.HandleFunc("GET /robots.txt", robotsTxt)
	mux.HandleFunc("GET /sitemap.xml", sitemapXML)

	// Protected routes
	auth := middleware.CheckAuth(store)

	mux.Handle("GET /admin", auth(http.HandlerFunc(handlers.AdminPage)))
	mux.Handle("GET /api/links", auth(http.HandlerFunc(handlers.ListLinks)))
	mux.Handle("POST /api/links", auth(http.HandlerFunc(handlers.CreateLink)))
	mux.Handle("PUT /api/links/reorder", auth(http.HandlerFunc(handlers.ReorderLinks)))
	mux.Handle("PUT /api/links/{id}", auth(http.HandlerFunc(handlers.UpdateLink)))
	mux.Handle("DELETE /api/links/{id}", auth(http.HandlerFunc(handlers.DeleteLink)))
	mux.Handle("GET /api/settings", auth(http.HandlerFunc(handlers.GetSettings)))
	mux.Handle("PUT /api/settings", auth(http.HandlerFunc(handlers.UpdateSettings)))

	// Experiences
	mux.HandleFunc("GET /api/experiences", handlers.ListExperiences)
	mux.Handle("POST /api/experiences", auth(http.HandlerFunc(handlers.CreateExperience)))
	mux.Handle("PUT /api/experiences/{id}", auth(http.HandlerFunc(handlers.UpdateExperience)))
	mux.Handle("DELETE /api/experiences/{id}", auth(http.HandlerFunc(handlers.DeleteExperience)))

	// Categories
	mux.HandleFunc("GET /api/categories", handlers.ListCategories)
	mux.Handle("POST /api/categories", auth(http.HandlerFunc(handlers.CreateCategory)))
	mux.Handle("PUT /api/categories/{id}", auth(http.HandlerFunc(handlers.UpdateCategory)))
	mux.Handle("DELETE /api/categories/{id}", auth(http.HandlerFunc(handlers.DeleteCategory)))

	// Files
	mux.HandleFunc("GET /api/files/public", handlers.ListPublicFiles)
	mux.HandleFunc("GET /api/files/download/{id}", handlers.DownloadFile)
	mux.Handle("GET /api/files", auth(http.HandlerFunc(handlers.ListFiles)))
	mux.Handle("POST /api/files/upload", auth(http.HandlerFunc(handlers.UploadFile)))
	mux.Handle("DELETE /api/files/{id}", auth(http.HandlerFunc(handlers.DeleteFile)))

	// Works
	mux.HandleFunc("GET /api/works/public", handlers.ListPublicWorks)
	mux.Handle("GET /api/works", auth(http.HandlerFunc(handlers.ListWorks)))
	mux.Handle("POST /api/works", auth(http.HandlerFunc(handlers.CreateWork)))
	mux.Handle("PUT /api/works/{id}", auth(http.HandlerFunc(handlers.UpdateWork)))
	mux.Handle("DELETE /api/works/{id}", auth(http.HandlerFunc(handlers.DeleteWork)))
	mux.Handle("POST /api/works/{id}/media", auth(http.HandlerFunc(handlers.AddWorkMedia)))
	mux.Handle("GET /api/works/{id}/media", auth(http.HandlerFunc(handlers.GetWorkMedia)))
	mux.Handle("DELETE /api/works/media/{id}", auth(http.HandlerFunc(handlers.DeleteWorkMedia)))

	// Chat
	mux.Handle("GET /chat", auth(http.HandlerFunc(handlers.ChatPage)))
	mux.Handle("GET /api/chat/sessions", auth(http.HandlerFunc(handlers.ListChatSessions)))
	mux.Handle("POST /api/chat/sessions", auth(http.HandlerFunc(handlers.CreateChatSession)))
	mux.Handle("DELETE /api/chat/sessions/{id}", auth(http.HandlerFunc(handlers.DeleteChatSession)))
	mux.Handle("PUT /api/chat/sessions/{id}", auth(http.HandlerFunc(handlers.UpdateChatSession)))
	mux.Handle("GET /api/chat/sessions/{id}/messages", auth(http.HandlerFunc(handlers.GetChatMessages)))
	mux.Handle("POST /api/chat/send", auth(http.HandlerFunc(handlers.SendChatMessage)))

	// SPA：已知前端路由返回 index.html（200），构建产物直出，其余一律真 404
	index200 := spaIndex(http.StatusOK)
	mux.HandleFunc("GET /{$}", index200)
	for _, p := range []string{"/works", "/work/{id}", "/about", "/experience", "/contact", "/ask"} {
		mux.HandleFunc("GET "+p, index200)
	}
	mux.Handle("GET /assets/", spaHandler("web/dist"))
	mux.HandleFunc("GET /", spaIndex(http.StatusNotFound))

	// 统一套一层 gzip（SSE 与大文件因 Content-Type 不匹配自动跳过）
	root := gzipMiddleware(mux)

	// Start
	addr := ":" + port
	srv := &http.Server{
		Addr:              addr,
		Handler:           root,
		ReadHeaderTimeout: 10 * time.Second,
		IdleTimeout:       120 * time.Second,
		// WriteTimeout 不设全局值：SSE 流式聊天与大文件下载会被全局写超时掐断
	}
	fmt.Printf("LinkTree server starting on http://0.0.0.0%s\n", addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
