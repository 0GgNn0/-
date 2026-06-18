package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"linktree/models"
)

func ListCategories(w http.ResponseWriter, r *http.Request) {
	categoryType := r.URL.Query().Get("type")
	if categoryType == "" {
		categoryType = "file"
	}

	categories, err := models.GetCategoriesByType(categoryType)
	if err != nil {
		writeJSON(w, map[string]interface{}{"ok": false, "message": err.Error()})
		return
	}
	writeJSON(w, categories)
}

func CreateCategory(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name string `json:"name"`
		Type string `json:"type"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]interface{}{"ok": false, "message": "请求格式错误"})
		return
	}
	if body.Name == "" || (body.Type != "file" && body.Type != "work") {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]interface{}{"ok": false, "message": "名称无效或类型错误"})
		return
	}

	category, err := models.CreateCategory(body.Name, body.Type)
	if err != nil {
		writeJSON(w, map[string]interface{}{"ok": false, "message": err.Error()})
		return
	}
	writeJSON(w, map[string]interface{}{"ok": true, "category": category})
}

func UpdateCategory(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]interface{}{"ok": false, "message": "无效ID"})
		return
	}

	var body struct {
		Name      string  `json:"name"`
		SortOrder float64 `json:"sort_order"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]interface{}{"ok": false, "message": "请求格式错误"})
		return
	}

	if err := models.UpdateCategory(id, body.Name, body.SortOrder); err != nil {
		writeJSON(w, map[string]interface{}{"ok": false, "message": err.Error()})
		return
	}
	writeJSON(w, map[string]interface{}{"ok": true})
}

func DeleteCategory(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]interface{}{"ok": false, "message": "无效ID"})
		return
	}

	if err := models.DeleteCategory(id); err != nil {
		writeJSON(w, map[string]interface{}{"ok": false, "message": err.Error()})
		return
	}
	writeJSON(w, map[string]interface{}{"ok": true})
}
