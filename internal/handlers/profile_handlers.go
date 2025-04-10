package handlers

import (
	ct "LoveMusic/internal/create_templates"
	"html/template"
	"net/http"
)

func LoadProfile(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := Authorise(w, r); err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	v := ct.GetChartUser()
	tmpl, _ := template.ParseFiles("static/templates/profile.html")
	tmpl.Execute(w, v)

	//http.ServeFile(w, r, "static/templates/profile.html")

}

func UserFriends(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := Authorise(w, r); err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	v := ct.GetChartUser()
	tmpl, _ := template.ParseFiles("static/templates/friends.html")
	tmpl.Execute(w, v)
}
