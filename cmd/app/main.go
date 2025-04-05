package main

import (
	db "LoveMusic/internal/database"
	h "LoveMusic/internal/handlers"
	"database/sql"
	"net/http"
)

var DB *sql.DB

func main() {

	if err := db.InitDatabase(); err != nil {
		panic(err)
	}
	defer db.Close()

	fileserver := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fileserver))

	http.HandleFunc("/", h.LoadProfile)
	http.HandleFunc("/register", h.RegisterHandler)
	http.HandleFunc("/login", h.LoginHandler)

	//http.HandleFunc("/profile", LoadProfile)

	db.ProfileDatabase()

	if err := http.ListenAndServe(":8082", nil); err != nil {
		panic(err)
	}
}
