package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"linktree/models"
)

func ListExperiences(w http.ResponseWriter, r *http.Request) {
	list, err := models.GetAllExperiences()
	if err != nil {
		writeJSON(w, map[string]interface{}{"ok": false, "message": err.Error()})
		return
	}
	writeJSON(w, list)
}

func CreateExperience(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Period      string `json:"period"`
		Title       string `json:"title"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]interface{}{"ok": false, "message": "请求格式错误"})
		return
	}
	if body.Period == "" || body.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]interface{}{"ok": false, "message": "时间段和标题不能为空"})
		return
	}
	e, err := models.CreateExperience(body.Period, body.Title, body.Description)
	if err != nil {
		writeJSON(w, map[string]interface{}{"ok": false, "message": err.Error()})
		return
	}
	writeJSON(w, map[string]interface{}{"ok": true, "experience": e})
}

func UpdateExperience(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]interface{}{"ok": false, "message": "无效ID"})
		return
	}
	var body struct {
		Period      string `json:"period"`
		Title       string `json:"title"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]interface{}{"ok": false, "message": "请求格式错误"})
		return
	}
	if body.Period == "" || body.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]interface{}{"ok": false, "message": "时间段和标题不能为空"})
		return
	}
	if err := models.UpdateExperience(id, body.Period, body.Title, body.Description); err != nil {
		writeJSON(w, map[string]interface{}{"ok": false, "message": err.Error()})
		return
	}
	writeJSON(w, map[string]interface{}{"ok": true})
}

func DeleteExperience(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]interface{}{"ok": false, "message": "无效ID"})
		return
	}
	if err := models.DeleteExperience(id); err != nil {
		writeJSON(w, map[string]interface{}{"ok": false, "message": err.Error()})
		return
	}
	writeJSON(w, map[string]interface{}{"ok": true})
}