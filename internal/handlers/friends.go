package handlers

import (
	ct "LoveMusic/internal/create_templates"
	db "LoveMusic/internal/database"
	"html/template"
	"net/http"
)

func UserFriends(w http.ResponseWriter, r *http.Request) {

	tmpl, _ := template.ParseFiles("static/templates/friends.html")

	if _, err := Authorise(w, r); err != nil {
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

		DBuser_id := db.CheckUserID("SELECT user_id FROM users WHERE user_id=$1", userID)
		if DBuser_id != userID {
			tmpl.Execute(w, ct.GetExeptionOnRegister("Нет такого пользователя"))
			return
		}

		http.Redirect(w, r, "/user/"+userID, http.StatusSeeOther)
	}

	tmpl.Execute(w, nil)
}
