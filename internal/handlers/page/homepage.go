package page

import (
	"LoveMusic/internal/database/pgsql/page"
	"LoveMusic/internal/handlers/auth"
	"html/template"
	"log/slog"
	"net/http"

	"github.com/redis/go-redis/v9"
)

type HomepageGetter interface {
	GetTopTracksUser(user_id int) []page.TopTracksUser
	GetUserID(login string) int
	SelectUser(user_id int) string
	CheckUserID(response string, args ...any) string
}

func New(log *slog.Logger, hompageGetter HomepageGetter, redisClient *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		login, err := auth.Authorise(w, r, redisClient)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		user_id := hompageGetter.GetUserID(login)

		username := hompageGetter.SelectUser(user_id)
		if username == "" {
			http.Error(w, "Page not found", http.StatusNotFound)
			return
		}

		viewprofile := struct {
			Username string
			User_ID  int
		}{
			Username: username,
			User_ID:  user_id,
		}

		topTracks := hompageGetter.GetTopTracksUser(user_id)

		data := struct {
			Profile struct {
				Username string
				User_ID  int
			}
			TopTracks []page.TopTracksUser
		}{
			Profile:   viewprofile,
			TopTracks: topTracks,
		}

		tmpl, err := template.ParseFiles("static/templates/user_profile.html")
		if err != nil {
			log.Error("error parse files")

			return
		}

		if err := tmpl.Execute(w, data); err != nil {
			log.Error("error execute data")
			return
		}

		//http.ServeFile(w, r, "static/templates/profile.html")
	}
}

// func UserProfileHandler(w http.ResponseWriter, r *http.Request) {
// 	if _, err := Authorise(w, r); err != nil {
// 		http.Redirect(w, r, "/login", http.StatusSeeOther)
// 		return
// 	}

// 	myurl := strings.Split(r.URL.Path, "/")
// 	if len(myurl) < 3 || myurl[2] == "" {
// 		fmt.Println("f")
// 		http.Error(w, "Not found", http.StatusNotFound)
// 		return
// 	}
// 	user_id, _ := strconv.Atoi(myurl[2])

// 	username := db.SelectUser(user_id)
// 	if username == "" {
// 		http.Error(w, "Page not found", http.StatusNotFound)
// 		return
// 	}

// 	viewprofile := struct {
// 		Username string
// 		User_ID  int
// 	}{
// 		Username: username,
// 		User_ID:  user_id,
// 	}

// 	topTracks := db.GetTopTracksUser(user_id)

// 	data := struct {
// 		Profile struct {
// 			Username string
// 			User_ID  int
// 		}
// 		TopTracks []db.TopTracksUser
// 	}{
// 		Profile:   viewprofile,
// 		TopTracks: topTracks,
// 	}

// 	tmpl, err := template.ParseFiles("static/templates/user_profile.html")
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// 	fmt.Println(data)

// 	if err := tmpl.Execute(w, data); err != nil {
// 		log.Println(err)
// 		return
// 	}

// }
