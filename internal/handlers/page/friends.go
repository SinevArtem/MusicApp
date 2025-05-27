package page

import (
	ct "LoveMusic/internal/create_templates"
	"LoveMusic/internal/handlers/auth"
	"html/template"
	"log/slog"
	"net/http"

	"github.com/redis/go-redis/v9"
)

func UserFriends(log *slog.Logger, homepageGetter HomepageGetter, redisClient *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, _ := template.ParseFiles("static/templates/friends.html")

		if _, err := auth.Authorise(w, r, redisClient); err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		if r.Method == http.MethodPost {
			err := r.ParseForm()
			if err != nil {
				http.Error(w, "Ошибка обработки формы", http.StatusBadRequest)
				return
			}

			userID := r.FormValue("user_id")

			if userID == "" {
				tmpl.Execute(w, ct.GetExeptionOnRegister("Не все поля заполнены"))
				return
			}

			DBuser_id := homepageGetter.CheckUserID("SELECT user_id FROM users WHERE user_id=$1", userID)
			if DBuser_id != userID {
				tmpl.Execute(w, ct.GetExeptionOnRegister("Нет такого пользователя"))
				return
			}

			http.Redirect(w, r, "/user/"+userID, http.StatusSeeOther)
		}

		tmpl.Execute(w, nil)
	}

}
