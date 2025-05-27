package tracks

import (
	ct "LoveMusic/internal/create_templates"
	"LoveMusic/internal/handlers/auth"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"

	"github.com/redis/go-redis/v9"
)

func AddTrack(log *slog.Logger, storageForTracks StorageForTracks, redisClient *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, err := auth.Authorise(w, r, redisClient); err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		tmpl, err := template.ParseFiles("static/templates/add_track.html")
		if err != nil {
			fmt.Println(err)
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

			dbTrack, dbArtist := storageForTracks.CheckTrackAndArtist(name_music, name_artist)
			if dbTrack == name_music && dbArtist == name_artist {
				tmpl.Execute(w, ct.TemplAddTrack("Такой трек уже есть", ""))
				return
			}
			storageForTracks.Insert("INSERT INTO tracks (name_music, name_artist) VALUES ($1,$2)", name_music, name_artist)
			tmpl.Execute(w, ct.TemplAddTrack("", "Трек добавлен"))

		}

		if err := tmpl.Execute(w, nil); err != nil {
			fmt.Println(err)
			return
		}
	}

}
