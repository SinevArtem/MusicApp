package page

import (
	"LoveMusic/internal/database/pgsql/page"
	"LoveMusic/internal/handlers/auth"
	"html/template"
	"log/slog"
	"net/http"

	"github.com/redis/go-redis/v9"
)

type CollectionGetter interface {
	GetUserID(login string) int
	GetTopTracksUser(user_id int) []page.TopTracksUser
}

func CollectionHandler(log *slog.Logger, collectionGetter CollectionGetter, redisClient *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		login, err := auth.Authorise(w, r, redisClient)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		user_id := collectionGetter.GetUserID(login)
		topTracks := collectionGetter.GetTopTracksUser(user_id)

		data := struct {
			TopTracks []page.TopTracksUser
		}{
			TopTracks: topTracks,
		}
		tmpl, _ := template.ParseFiles("static/templates/collection.html")
		tmpl.Execute(w, data)

	}

}
