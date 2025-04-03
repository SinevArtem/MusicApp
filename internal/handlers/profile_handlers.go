package handlers

import (
	ct "LoveMusic/internal/create_templates"
	"html/template"
	"net/http"
)

func LoadProfile(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	v := ct.GetChartUser()
	tmpl, _ := template.ParseFiles("static/templates/profile.html")
	tmpl.Execute(w, v)

	//http.ServeFile(w, r, "static/templates/profile.html")

}
