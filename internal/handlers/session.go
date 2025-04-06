package handlers

import (
	db "LoveMusic/internal/database"
	"context"
	"errors"
	"net/http"
)

var ErrAuth = errors.New("Unauthorized")

func Authorise(r *http.Request) error {

	sessionToken, err := r.Cookie("session_token")
	if err != nil || sessionToken.Value == "" {
		return ErrAuth
	}

	ctx := context.Background()
	login, err := db.RedisDB.Get(ctx, "session:"+sessionToken.Value).Result()
	if err != nil {
		return ErrAuth
	}

	csrfToken, err := r.Cookie("csrf_token")
	if err != nil || csrfToken.Value == "" {
		return ErrAuth
	}

	storedLogin, err := db.RedisDB.Get(ctx, "csrf:"+csrfToken.Value).Result()
	if err != nil || storedLogin != login {
		return ErrAuth
	}

	return nil

}

// func CheckCSRV(r *http.Request, login string) bool {
// 	cookieToken, err := r.Cookie("csrf_token")
// 	if err != nil || cookieToken.Value == "" {
// 		return false
// 	}

// 	ctx := context.Background()
// 	storedLogin, err := db.RedisDB.Get(ctx, "csrf:"+cookieToken.Value).Result()
// 	if err != nil {
// 		return false
// 	}
// 	fmt.Printf("%s %s", login, storedLogin)
// 	return storedLogin == login
// }

// func CheckSession(r *http.Request, login string) bool {
// 	cookieToken, err := r.Cookie("session_token")
// 	if err != nil || cookieToken.Value == "" {
// 		return false
// 	}

// 	ctx := context.Background()
// 	storedLogin, err := db.RedisDB.Get(ctx, "session:"+cookieToken.Value).Result()
// 	if err != nil {
// 		return false
// 	}

// 	fmt.Printf("%s %s", login, storedLogin)
// 	return storedLogin == login
// }
