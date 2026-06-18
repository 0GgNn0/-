package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"linktree/models"
)

func ListWorks(w http.ResponseWriter, r *http.Request) {
	var categoryID *int
	if cidStr := r.URL.Query().Get("category_id"); cidStr != "" {
		cid, err := strconv.Atoi(cidStr)
		if err == nil {
			categoryID = &cid
		}
	}

	works, err := models.GetWorksByCategory(categoryID)
	if err != nil {
		writeJSON(w, map[string]interface{}{"ok": false, "message": err.Error()})
		return
	}
	models.LoadWorkMedia(works)
	writeJSON(w, works)
}

func ListPublicWorks(w http.ResponseWriter, r *http.Request) {
	works, err := models.GetPublicWorks()
	if err != nil {
		writeJSON(w, map[string]interface{}{"ok": false, "message": err.Error()})
		return
	}
	models.LoadWorkMedia(works)
	writeJSON(w, works)
}

func CreateWork(w http.ResponseWriter, r *http.Request) {
	var body struct {
		CategoryID  *int   `json:"category_id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		CoverURL    string `json:"cover_url"`
		ContentType string `json:"content_type"`
		ExternalURL string `json:"external_url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]interface{}{"ok": false, "message": "请求格式错误"})
		return
	}
	if body.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]interface{}{"ok": false, "message": "标题不能为空"})
		return
	}
	if body.ContentType == "" {
		body.ContentType = "image"
	}

	work, err := models.CreateWork(body.CategoryID, body.Title, body.Description, body.CoverURL, body.ContentType, body.ExternalURL)
	if err != nil {
		writeJSON(w, map[string]interface{}{"ok": false, "message": err.Error()})
		return
	}
	writeJSON(w, map[string]interface{}{"ok": true, "work": work})
}

func UpdateWork(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]interface{}{"ok": false, "message": "无效ID"})
		return
	}

	var body struct {
		CategoryID  *int   `json:"category_id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		CoverURL    string `json:"cover_url"`
		ContentType string `json:"content_type"`
		ExternalURL string `json:"external_url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]interface{}{"ok": false, "message": "请求格式错误"})
		return
	}

	if err := models.UpdateWork(id, body.CategoryID, body.Title, body.Description, body.CoverURL, body.ContentType, body.ExternalURL); err != nil {
		writeJSON(w, map[string]interface{}{"ok": false, "message": err.Error()})
		return
	}
	writeJSON(w, map[string]interface{}{"ok": true})
}

func DeleteWork(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]interface{}{"ok": false, "message": "无效ID"})
		return
	}

	if err := models.DeleteWork(id); err != nil {
		writeJSON(w, map[string]interface{}{"ok": false, "message": err.Error()})
		return
	}
	writeJSON(w, map[string]interface{}{"ok": true})
}

func AddWorkMedia(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	workID, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]interface{}{"ok": false, "message": "无效ID"})
		return
	}

	var body struct {
		FileID    *int   `json:"file_id"`
		MediaURL  string `json:"media_url"`
		MediaType string `json:"media_type"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]interface{}{"ok": false, "message": "请求格式错误"})
		return
	}
	if body.MediaURL == "" {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]interface{}{"ok": false, "message": "媒体URL不能为空"})
		return
	}
	if body.MediaType == "" {
		body.MediaType = "image"
	}

	media, err := models.AddWorkMedia(workID, body.FileID, body.MediaURL, body.MediaType)
	if err != nil {
		writeJSON(w, map[string]interface{}{"ok": false, "message": err.Error()})
		return
	}
	writeJSON(w, map[string]interface{}{"ok": true, "media": media})
}

func DeleteWorkMedia(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]interface{}{"ok": false, "message": "无效ID"})
		return
	}

	if err := models.DeleteWorkMedia(id); err != nil {
		writeJSON(w, map[string]interface{}{"ok": false, "message": err.Error()})
		return
	}
	writeJSON(w, map[string]interface{}{"ok": true})
}

func GetWorkMedia(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]interface{}{"ok": false, "message": "无效ID"})
		return
	}

	media, err := models.GetWorkMedia(id)
	if err != nil {
		writeJSON(w, map[string]interface{}{"ok": false, "message": err.Error()})
		return
	}
	writeJSON(w, media)
}
