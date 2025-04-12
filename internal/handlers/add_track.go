package handlers

import (
	ct "LoveMusic/internal/create_templates"
	db "LoveMusic/internal/database"
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

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Ошибка обработки формы", http.StatusBadRequest)
			return
		}

		name_music := r.FormValue("name_music")
		name_artist := r.FormValue("name_artist")

		if name_music == "" || name_artist == "" {
			tmpl.Execute(w, ct.TemplAddTrack("Не все поля заполнены", ""))
			return
		}

		dbTrack, dbArtist := db.CheckTrackAndArtist("SELECT name_music, name_artist FROM tracks WHERE name_music=$1 AND name_artist=$2", name_music, name_artist)
		if dbTrack == name_music && dbArtist == name_artist {
			tmpl.Execute(w, ct.TemplAddTrack("Такой трек уже есть", ""))
			return
		}
		db.InsertResponseDatabase("INSERT INTO tracks (name_music, name_artist) VALUES ($1,$2)", name_music, name_artist)
		tmpl.Execute(w, ct.TemplAddTrack("", "Трек добавлен"))

	}

	if err := tmpl.Execute(w, nil); err != nil {
		log.Println(err)
		return
	}
}
