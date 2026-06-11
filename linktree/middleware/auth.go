package middleware

import (
	"net/http"

	"github.com/gorilla/sessions"
)

func CheckAuth(store *sessions.CookieStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, _ := store.Get(r, "linktree_session")
			if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
				// For API requests, return JSON error
				if len(r.URL.Path) >= 4 && r.URL.Path[:4] == "/api" {
					w.Header().Set("Content-Type", "application/json; charset=utf-8")
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte(`{"ok":false,"message":"未登录"}`))
					return
				}
				// For page requests, redirect to login
				http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
