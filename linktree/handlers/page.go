package handlers

import (
	"html/template"
	"net/http"

	"linktree/models"
)

var tmpl *template.Template

func InitTemplates() error {
	var err error
	tmpl, err = template.ParseGlob("templates/*.html")
	return err
}

type PageData struct {
	DisplayName string
	Bio         string
	AvatarURL   string
	Links       []models.Link
	Theme       string
}

func IndexPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	links, _ := models.GetAllLinks()
	data := PageData{
		DisplayName: models.GetSettingOrDefault("display_name", "我的主页"),
		Bio:         models.GetSettingOrDefault("bio", ""),
		AvatarURL:   models.GetSettingOrDefault("avatar_data_url", ""),
		Links:       links,
		Theme:       models.GetSettingOrDefault("theme", "auto"),
	}
	tmpl.ExecuteTemplate(w, "index.html", data)
}

func LoginPage(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "login.html", nil)
}

func AdminPage(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "admin.html", nil)
}
