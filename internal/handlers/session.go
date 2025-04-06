package handlers

import (
	db "LoveMusic/internal/database"
	"context"
	"errors"
	"fmt"
	"net/http"
)

var AuthError = errors.New("Unauthorized")

func Authorise(r *http.Request) error {

	login := r.FormValue("login")
	if login == "" {
		return AuthError
	}

	if !CheckCSRV(r, login) || !CheckSession(r, login) {
		return AuthError
	}

	return nil

}

func CheckCSRV(r *http.Request, login string) bool {
	cookieToken, err := r.Cookie("csrf_token")
	if err != nil || cookieToken.Value == "" {
		return false
	}

	ctx := context.Background()
	storedLogin, err := db.RedisDB.Get(ctx, "csrf:"+cookieToken.Value).Result()
	if err != nil {
		return false
	}
	fmt.Printf("%s %s", login, storedLogin)
	return storedLogin == login
}

func CheckSession(r *http.Request, login string) bool {
	cookieToken, err := r.Cookie("session_token")
	if err != nil || cookieToken.Value == "" {
		return false
	}

	ctx := context.Background()
	storedLogin, err := db.RedisDB.Get(ctx, "session:"+cookieToken.Value).Result()
	if err != nil {
		return false
	}

	fmt.Printf("%s %s", login, storedLogin)
	return storedLogin == login
}
