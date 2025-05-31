package tracks

import (
	"LoveMusic/internal/database/pgsql/tracks"
	"LoveMusic/internal/handlers/auth"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"

	"github.com/redis/go-redis/v9"
)

type StorageForTracks interface {
	GetTracks(track, artist string) []tracks.Tracks
	CheckTrackAndArtist(args ...any) (string, string)
	Insert(response string, args ...any) error
	GetUserID(login string) int
	GetTrackID(name_music, name_artist string) int
	AddToCollection(user_id, track_id int) error
}

type SearchData struct {
	Tracks []tracks.Tracks
}

func SearchTrack(log *slog.Logger, storageForTracks StorageForTracks, redisClient *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, err := auth.Authorise(w, r, redisClient); err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		tmpl, _ := template.ParseFiles("static/templates/search_track.html")
		if r.Method == http.MethodPost {
			err := r.ParseForm()
			if err != nil {
				http.Error(w, "Ошибка обработки формы", http.StatusBadRequest)
				return
			}

			name_music := r.FormValue("name_music")
			name_artist := r.FormValue("name_artist")

			data := SearchData{}
			data.Tracks = storageForTracks.GetTracks(name_music, name_artist)

			err = tmpl.Execute(w, data)
			if err != nil {
				fmt.Println(err)
			}
			return
		}

		tmpl.Execute(w, nil)
	}

}

func AddToCollection(log *slog.Logger, storageForTracks StorageForTracks, redisClient *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		login, err := auth.Authorise(w, r, redisClient)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		user_id := storageForTracks.GetUserID(login)

		if err := r.ParseForm(); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		name_music := r.FormValue("n_music")
		name_artist := r.FormValue("n_artist")

		fmt.Println(name_music, name_artist)

		track_id := storageForTracks.GetTrackID(name_music, name_artist)

		if err := storageForTracks.AddToCollection(user_id, track_id); err != nil {
			log.Error("Failed to add track to collection", "error", err)
			http.Error(w, "Failed to add track to collection", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status": "success"}`))
	}
}
