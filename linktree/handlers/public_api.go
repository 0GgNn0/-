package handlers

import (
	"net/http"

	"linktree/models"
)

func PublicProfile(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, map[string]string{
		"display_name": models.GetSettingOrDefault("display_name", "啊芃"),
		"bio":          models.GetSettingOrDefault("bio", ""),
		"theme":        models.GetSettingOrDefault("theme", "auto"),
	})
}

func PublicLinks(w http.ResponseWriter, r *http.Request) {
	links, err := models.GetAllLinks()
	if err != nil {
		writeJSON(w, map[string]interface{}{"ok": false, "message": err.Error()})
		return
	}
	writeJSON(w, links)
}
