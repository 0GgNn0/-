package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"linktree/models"

	"github.com/gorilla/sessions"
)

var store *sessions.CookieStore

func SetStore(s *sessions.CookieStore) {
	store = s
}

func writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(data)
}

// ---- Auth handlers ----

func Login(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]interface{}{"ok": false, "message": "请求格式错误"})
		return
	}

	if !models.VerifyPassword(body.Password) {
		w.WriteHeader(http.StatusUnauthorized)
		writeJSON(w, map[string]interface{}{"ok": false, "message": "密码错误"})
		return
	}

	session, _ := store.Get(r, "linktree_session")
	session.Values["authenticated"] = true
	session.Save(r, w)

	writeJSON(w, map[string]interface{}{"ok": true, "message": "登录成功"})
}

func Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "linktree_session")
	session.Values["authenticated"] = false
	session.Save(r, w)
	writeJSON(w, map[string]interface{}{"ok": true})
}

// ---- Links handlers ----

func ListLinks(w http.ResponseWriter, r *http.Request) {
	links, err := models.GetAllLinks()
	if err != nil {
		writeJSON(w, map[string]interface{}{"ok": false, "message": err.Error()})
		return
	}
	writeJSON(w, links)
}

func CreateLink(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Title string `json:"title"`
		URL   string `json:"url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]interface{}{"ok": false, "message": "请求格式错误"})
		return
	}
	if body.Title == "" || body.URL == "" {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]interface{}{"ok": false, "message": "标题和URL不能为空"})
		return
	}

	link, err := models.CreateLink(body.Title, body.URL)
	if err != nil {
		writeJSON(w, map[string]interface{}{"ok": false, "message": err.Error()})
		return
	}
	writeJSON(w, map[string]interface{}{"ok": true, "link": link})
}

func UpdateLink(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]interface{}{"ok": false, "message": "无效ID"})
		return
	}

	var body struct {
		Title     string  `json:"title"`
		URL       string  `json:"url"`
		SortOrder float64 `json:"sort_order"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]interface{}{"ok": false, "message": "请求格式错误"})
		return
	}

	if err := models.UpdateLink(id, body.Title, body.URL, body.SortOrder); err != nil {
		writeJSON(w, map[string]interface{}{"ok": false, "message": err.Error()})
		return
	}
	writeJSON(w, map[string]interface{}{"ok": true})
}

func DeleteLink(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]interface{}{"ok": false, "message": "无效ID"})
		return
	}

	if err := models.DeleteLink(id); err != nil {
		writeJSON(w, map[string]interface{}{"ok": false, "message": err.Error()})
		return
	}
	writeJSON(w, map[string]interface{}{"ok": true})
}

func ReorderLinks(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Orders []models.OrderEntry `json:"orders"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]interface{}{"ok": false, "message": "请求格式错误"})
		return
	}

	if err := models.ReorderLinks(body.Orders); err != nil {
		writeJSON(w, map[string]interface{}{"ok": false, "message": err.Error()})
		return
	}
	writeJSON(w, map[string]interface{}{"ok": true})
}

// ---- Settings handlers ----

func GetSettings(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"display_name":    models.GetSettingOrDefault("display_name", "我的主页"),
		"bio":             models.GetSettingOrDefault("bio", ""),
		"avatar_data_url": models.GetSettingOrDefault("avatar_data_url", ""),
		"theme":           models.GetSettingOrDefault("theme", "auto"),
		"ai_base_url":     models.GetAIBaseURL(),
		"ai_api_key":      models.MaskAPIKey(models.GetAIAPIKey()),
		"ai_model":        models.GetAIModel(),
	}
	writeJSON(w, data)
}

func UpdateSettings(w http.ResponseWriter, r *http.Request) {
	var body struct {
		DisplayName *string `json:"display_name"`
		Bio         *string `json:"bio"`
		AvatarURL   *string `json:"avatar_data_url"`
		NewPassword *string `json:"new_password"`
		AIBaseURL   *string `json:"ai_base_url"`
		AIAPIKey    *string `json:"ai_api_key"`
		AIModel     *string `json:"ai_model"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]interface{}{"ok": false, "message": "请求格式错误"})
		return
	}

	if body.DisplayName != nil {
		models.SetSetting("display_name", *body.DisplayName)
	}
	if body.Bio != nil {
		models.SetSetting("bio", *body.Bio)
	}
	if body.AvatarURL != nil {
		models.SetSetting("avatar_data_url", *body.AvatarURL)
	}
	if body.NewPassword != nil && *body.NewPassword != "" {
		if err := models.ChangePassword(*body.NewPassword); err != nil {
			writeJSON(w, map[string]interface{}{"ok": false, "message": "密码修改失败: " + err.Error()})
			return
		}
	}
	if body.AIBaseURL != nil {
		models.SetSetting("ai_base_url", *body.AIBaseURL)
	}
	if body.AIAPIKey != nil {
		// Skip if masked placeholder
		if *body.AIAPIKey != "" && !strings.Contains(*body.AIAPIKey, "****") {
			models.SetSetting("ai_api_key", *body.AIAPIKey)
		}
	}
	if body.AIModel != nil {
		models.SetSetting("ai_model", *body.AIModel)
	}

	writeJSON(w, map[string]interface{}{"ok": true})
}
