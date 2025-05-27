package auth

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

var ErrAuth = errors.New("Unauthorized")

func Authorise(w http.ResponseWriter, r *http.Request, redisClient *redis.Client) (string, error) {

	sessionToken, err := r.Cookie("session_token")
	if err != nil || sessionToken.Value == "" {
		return "", ErrAuth
	}

	ctx := context.Background()

	login, err := redisClient.Get(ctx, "session:"+sessionToken.Value).Result()
	if err != nil {
		return "", ErrAuth
	}

	csrfToken, err := r.Cookie("csrf_token")
	if err != nil || csrfToken.Value == "" {
		return "", ErrAuth
	}

	storedLogin, err := redisClient.Get(ctx, "csrf:"+csrfToken.Value).Result()
	if err != nil || storedLogin != login {
		return "", ErrAuth
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken.Value,
		Expires:  time.Now().Add(30 * time.Minute),
		HttpOnly: true,  // javascript не получит токен
		Secure:   false, // при HTTPS true
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    csrfToken.Value,
		Expires:  time.Now().Add(30 * time.Minute),
		HttpOnly: false,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})

	if err := redisClient.Set(ctx, "session:"+sessionToken.Value, login, 30*time.Minute).Err(); err != nil {
		log.Println("Не удалось продлить сессию:", err)
	}

	if err := redisClient.Set(ctx, "csrf:"+csrfToken.Value, login, 30*time.Minute).Err(); err != nil {
		log.Println("Не удалось продлить сессию:", err)
	}

	return login, nil

}
