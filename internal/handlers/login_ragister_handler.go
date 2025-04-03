package handlers

import (
	"LoveMusic/internal/database"
	"fmt"
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
	tmpl.Execute(w, nil)

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
			fmt.Println("не все введено")
			return
		}

		if strings.Compare(password, check_password) == 0 {
			password_hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

			database.InsertResponseDatabase("INSERT INTO users (username, login, password) VALUES ($1, $2, $3);", username, login, password_hash)
		} else {

			fmt.Println("пароли не совпадают")
			return
		}

	}

}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/login" {
		http.NotFound(w, r)
		return
	}
	tmpl, _ := template.ParseFiles("static/templates/login.html")
	tmpl.Execute(w, nil)

}
