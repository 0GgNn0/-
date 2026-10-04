package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"linktree/handlers"
	"linktree/middleware"
	"linktree/models"

	"github.com/gorilla/sessions"
)

func noCache(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache, must-revalidate")
		next.ServeHTTP(w, r)
	})
}

// spaHandler 服务 Vite 构建产物：真实文件直出，其余回退 index.html（前端路由）
func spaHandler(distDir string) http.Handler {
	fileServer := http.FileServer(http.Dir(distDir))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := filepath.Clean(r.URL.Path)
		full := filepath.Join(distDir, p)
		if st, err := os.Stat(full); err == nil && !st.IsDir() {
			if len(p) >= len("/assets/") && p[:8] == "/assets/" {
				w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
			} else {
				w.Header().Set("Cache-Control", "no-cache, must-revalidate")
			}
			fileServer.ServeHTTP(w, r)
			return
		}
		// 未知 API 路径返回 404 JSON，避免前端 index 污染 API 语义
		if len(r.URL.Path) >= 4 && r.URL.Path[:4] == "/api" {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"ok":false,"message":"not found"}`))
			return
		}
		w.Header().Set("Cache-Control", "no-cache, must-revalidate")
		http.ServeFile(w, r, filepath.Join(distDir, "index.html"))
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

	// Static files
	mux.Handle("GET /static/", http.StripPrefix("/static/", noCache(http.FileServer(http.Dir("static")))))

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

	// SPA fallback：非上述模式的 GET 请求交给 Vue 前端（web/dist）
	mux.Handle("GET /", spaHandler("web/dist"))

	// Start
	addr := ":" + port
	srv := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
		IdleTimeout:       120 * time.Second,
		// WriteTimeout 不设全局值：SSE 流式聊天与大文件下载会被全局写超时掐断
	}
	fmt.Printf("LinkTree server starting on http://0.0.0.0%s\n", addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
