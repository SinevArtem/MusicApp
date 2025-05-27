package auth

import (
	ct "LoveMusic/internal/create_templates"
	generatetoken "LoveMusic/internal/handlers/auth/generate_token"

	"context"
	"html/template"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type logregSaver interface {
	CheckLoginDatabase(login string) string
	InsertRegisterValue(username, login, password_hash string) error
	SelectLoginOrPasswordOnDatabase(login string) (string, string)
}

func RegisterHandler(log *slog.Logger, logregSaver logregSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

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

			if logregSaver.CheckLoginDatabase(login) == login {
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
			logregSaver.InsertRegisterValue(username, login, string(password_hash))

			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return

		}

		tmpl.Execute(w, nil)
	}
	// if r.URL.Path != "/register" {
	// 	http.NotFound(w, r)
	// 	return
	// }

}

func LoginPageHandler(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, _ := template.ParseFiles("static/templates/login.html")

		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Ошибка обработки формы", http.StatusBadRequest)
			return
		}

		tmpl.Execute(w, nil)
	}
}

func LoginHandler(log *slog.Logger, logregSaver logregSaver, redisClient *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

			db_login, db_password := logregSaver.SelectLoginOrPasswordOnDatabase(login)
			if db_login != login {
				tmpl.Execute(w, ct.GetExeptionOnRegister("Неправильный логин или пароль"))
				return
			}

			err = bcrypt.CompareHashAndPassword([]byte(db_password), []byte(password))
			if err != nil {
				tmpl.Execute(w, ct.GetExeptionOnRegister("Неправильный логин или пароль"))
				return
			}

			sessionToken := generatetoken.GenerateToken(32)
			csrfToken := generatetoken.GenerateToken(32)

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			for {
				exists, err := redisClient.Exists(ctx, "session:"+sessionToken).Result() // Result возвращает результат и ошибку, если 1 - такой ключ есть ,0 - нет
				if err != nil || exists == 0 {
					break
				}
				sessionToken = generatetoken.GenerateToken(32)
			}

			for {
				exists, err := redisClient.Exists(ctx, "csrf:"+csrfToken).Result()
				if err != nil || exists == 0 {
					break
				}
				csrfToken = generatetoken.GenerateToken(32)
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

			if err := redisClient.Set(ctx, "session:"+sessionToken, login, 30*time.Minute).Err(); err != nil {
				http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
				return
			} // Err() возвращает только ошибку без результата

			if err := redisClient.Set(ctx, "csrf:"+csrfToken, login, 30*time.Minute).Err(); err != nil {
				http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, "/profile", http.StatusSeeOther)
		}

		tmpl.Execute(w, nil)
	}

}

func LogoutHandler(log *slog.Logger, redisClient *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, err := Authorise(w, r, redisClient); err != nil {
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

		redisClient.Set(ctx, "session:"+sessionToken.Value, "", 30*time.Minute)
		redisClient.Set(ctx, "csrf:"+csrfToken.Value, "", 30*time.Minute)

		http.Redirect(w, r, "/profile", http.StatusSeeOther)
	}

}
