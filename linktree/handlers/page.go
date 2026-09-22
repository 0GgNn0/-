package handlers

import (
	"html/template"
	"net/http"
	"strconv"

	"linktree/models"
)

var tmpl *template.Template

func InitTemplates() error {
	var err error
	funcs := template.FuncMap{
		"mul": func(a, b interface{}) float64 {
			af, _ := strconv.ParseFloat(toStr(a), 64)
			bf, _ := strconv.ParseFloat(toStr(b), 64)
			return af * bf
		},
	}
	tmpl, err = template.New("").Funcs(funcs).ParseGlob("templates/*.html")
	return err
}

func toStr(v interface{}) string {
	switch x := v.(type) {
	case int:
		return strconv.Itoa(x)
	case float64:
		return strconv.FormatFloat(x, 'f', -1, 64)
	case string:
		return x
	default:
		return ""
	}
}

type PageData struct {
	DisplayName string
	Bio         string
	AvatarURL   string
	Links       []models.Link
	Works       []models.Work
	Files       []models.File
	Theme       string
}

func IndexPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	links, _ := models.GetAllLinks()
	works, _ := models.GetPublicWorks()
	models.LoadWorkMedia(works)
	files, _ := models.GetPublicFiles()
	data := PageData{
		DisplayName: models.GetSettingOrDefault("display_name", "我的主页"),
		Bio:         models.GetSettingOrDefault("bio", ""),
		AvatarURL:   models.GetSettingOrDefault("avatar_data_url", ""),
		Links:       links,
		Works:       works,
		Files:       files,
		Theme:       models.GetSettingOrDefault("theme", "auto"),
	}
	tmpl.ExecuteTemplate(w, "index.html", data)
}

func LoginPage(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "login.html", nil)
}

// ---- 极简多页结构 handler ----

type WorksPageData struct {
	DisplayName string
	Works       []models.Work
	Theme       string
}

func WorksPage(w http.ResponseWriter, r *http.Request) {
	works, _ := models.GetPublicWorks()
	models.LoadWorkMedia(works)
	data := WorksPageData{
		DisplayName: models.GetSettingOrDefault("display_name", "啊芃"),
		Works:       works,
		Theme:       models.GetSettingOrDefault("theme", "auto"),
	}
	tmpl.ExecuteTemplate(w, "works.html", data)
}

type WorkDetailData struct {
	DisplayName string
	Work        *models.Work
	Theme       string
}

func WorkDetailPage(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	work, err := models.GetWorkByID(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	media, _ := models.GetWorkMedia(id)
	work.Media = media
	data := WorkDetailData{
		DisplayName: models.GetSettingOrDefault("display_name", "啊芃"),
		Work:        work,
		Theme:       models.GetSettingOrDefault("theme", "auto"),
	}
	tmpl.ExecuteTemplate(w, "work.html", data)
}

func AboutPage(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"DisplayName": models.GetSettingOrDefault("display_name", "啊芃"),
		"Bio":         models.GetSettingOrDefault("bio", ""),
		"Theme":       models.GetSettingOrDefault("theme", "auto"),
	}
	tmpl.ExecuteTemplate(w, "about.html", data)
}

func ExperiencePage(w http.ResponseWriter, r *http.Request) {
	experiences, _ := models.GetAllExperiences()
	data := map[string]interface{}{
		"DisplayName": models.GetSettingOrDefault("display_name", "啊芃"),
		"Experiences": experiences,
		"Theme":       models.GetSettingOrDefault("theme", "auto"),
	}
	tmpl.ExecuteTemplate(w, "experience.html", data)
}

func ContactPage(w http.ResponseWriter, r *http.Request) {
	links, _ := models.GetAllLinks()
	data := map[string]interface{}{
		"DisplayName": models.GetSettingOrDefault("display_name", "啊芃"),
		"Links":       links,
		"Theme":       models.GetSettingOrDefault("theme", "auto"),
	}
	tmpl.ExecuteTemplate(w, "contact.html", data)
}

func AdminPage(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "admin.html", nil)
}

func ChatPage(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "chat.html", nil)
}
