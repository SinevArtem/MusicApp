package handlers

import (
	ct "LoveMusic/internal/create_templates"
	db "LoveMusic/internal/database"
	"html/template"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// if r.URL.Path != "/register" {
	// 	http.NotFound(w, r)
	// 	return
	// }

	tmpl, _ := template.ParseFiles("static/templates/register.html")

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Ошибка обработки формы", http.StatusBadRequest)
			return
		}

		username := r.FormValue("username")
		login := r.FormValue("login")
		password := r.FormValue("password")
		check_password := r.FormValue("check_password")

		if username == "" || login == "" || password == "" || check_password == "" {
			tmpl.Execute(w, ct.GetExeptionOnRegister("Не все поля заполнены"))
			return
		}

		if len(username) < 4 && len(login) < 4 {
			tmpl.Execute(w, ct.GetExeptionOnRegister("Username и Логин должны состоять минимум из 4 символов"))
			return
		}

		if strings.Compare(password, check_password) != 0 {
			tmpl.Execute(w, ct.GetExeptionOnRegister("Пароли не совпадают"))
			return
		}

		if len(password) < 8 {
			tmpl.Execute(w, ct.GetExeptionOnRegister("Пароль должен состоять минимум из 8 символов"))
			return
		}

		password_hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		db.InsertResponseDatabase("INSERT INTO users (username, login, password) VALUES ($1, $2, $3);", username, login, password_hash)

		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return

	}

	tmpl.Execute(w, nil)

}

func LoginHandler(w http.ResponseWriter, r *http.Request) {

	tmpl, _ := template.ParseFiles("static/templates/login.html")

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Ошибка обработки формы", http.StatusBadRequest)
			return
		}

		login := r.FormValue("login")
		password := r.FormValue("password")

		if login == "" || password == "" {
			tmpl.Execute(w, ct.GetExeptionOnRegister("Не все поля заполнены"))
			return
		}

		login_and_password := db.SelectLoginOrPasswordOnDatabase(login)
		if login_and_password.Login != login {
			tmpl.Execute(w, ct.GetExeptionOnRegister("Неправильный логин или пароль"))
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(login_and_password.Password), []byte(password))
		if err != nil {
			tmpl.Execute(w, ct.GetExeptionOnRegister("Неправильный логин или пароль"))
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	tmpl.Execute(w, nil)

}
