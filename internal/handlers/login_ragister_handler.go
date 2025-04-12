package handlers

import (
	ct "LoveMusic/internal/create_templates"
	db "LoveMusic/internal/database"
	"context"
	"html/template"
	"net/http"
	"strings"
	"time"

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

		if db.CheckLoginDatabase("SELECT login FROM users WHERE login=$1", login) == login {
			tmpl.Execute(w, ct.GetExeptionOnRegister("Такой пользователь уже есть"))
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

		sessionToken := generateToken(32)
		csrfToken := generateToken(32)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		for {
			exists, err := db.RedisDB.Exists(ctx, "session:"+sessionToken).Result() // Result возвращает результат и ошибку, если 1 - такой ключ есть ,0 - нет
			if err != nil || exists == 0 {
				break
			}
			sessionToken = generateToken(32)
		}

		for {
			exists, err := db.RedisDB.Exists(ctx, "csrf:"+csrfToken).Result()
			if err != nil || exists == 0 {
				break
			}
			csrfToken = generateToken(32)
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    sessionToken,
			Expires:  time.Now().Add(30 * time.Minute),
			HttpOnly: true,  // javascript не получит токен
			Secure:   false, // при HTTPS true
			SameSite: http.SameSiteLaxMode,
			Path:     "/",
		})

		http.SetCookie(w, &http.Cookie{
			Name:     "csrf_token",
			Value:    csrfToken,
			Expires:  time.Now().Add(30 * time.Minute),
			HttpOnly: false,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
			Path:     "/",
		})

		if err := db.RedisDB.Set(ctx, "session:"+sessionToken, login, 30*time.Minute).Err(); err != nil {
			http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
			return
		} // Err() возвращает только ошибку без результата

		if err := db.RedisDB.Set(ctx, "csrf:"+csrfToken, login, 30*time.Minute).Err(); err != nil {
			http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/profile", http.StatusSeeOther)
	}

	tmpl.Execute(w, nil)

}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if _, err := Authorise(w, r); err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Now().Add(-(30 * time.Minute)),
		HttpOnly: true,  // javascript не получит токен
		Secure:   false, // при HTTPS true
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    "",
		Expires:  time.Now().Add(-(30 * time.Minute)),
		HttpOnly: false,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})

	sessionToken, _ := r.Cookie("session_token")
	csrfToken, _ := r.Cookie("csrf_token")

	ctx := context.Background()

	db.RedisDB.Set(ctx, "session:"+sessionToken.Value, "", 30*time.Minute)
	db.RedisDB.Set(ctx, "csrf:"+csrfToken.Value, "", 30*time.Minute)

	http.Redirect(w, r, "/profile", http.StatusSeeOther)

}
