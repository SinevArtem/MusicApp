package handlers

import (
	db "LoveMusic/internal/database"
	"html/template"
	"net/http"
)

func CollectionHandler(w http.ResponseWriter, r *http.Request) {
	login, err := Authorise(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	user_id := db.GetUserID(login)
	topTracks := db.GetTopTracksUser(user_id)

	data := struct {
		TopTracks []db.TopTracksUser
	}{
		TopTracks: topTracks,
	}
	tmpl, _ := template.ParseFiles("static/templates/collection.html")
	tmpl.Execute(w, data)

}
