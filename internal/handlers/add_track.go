package handlers

import (
	"html/template"
	"log"
	"net/http"
)

func AddTrack(w http.ResponseWriter, r *http.Request) {
	if _, err := Authorise(w, r); err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	tmpl, err := template.ParseFiles("static/templates/add_track.html")
	if err != nil {
		log.Println(err)
		return
	}

	if err := tmpl.Execute(w, r); err != nil {
		log.Println(err)
		return
	}
}
