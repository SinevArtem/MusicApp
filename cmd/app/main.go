package main

import (
	db "LoveMusic/internal/database"
	h "LoveMusic/internal/handlers"
	"database/sql"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

var DB *sql.DB

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("Env not found")
	}

	if err := db.InitDatabase(); err != nil {
		panic(err)
	}
	defer db.Close()

	fileserver := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fileserver))

	http.HandleFunc("/", h.LoadProfile)
	http.HandleFunc("/register", h.RegisterHandler)
	http.HandleFunc("/login", h.LoginHandler)
	http.HandleFunc("/protected", h.Protected)

	//http.HandleFunc("/profile", LoadProfile)

	db.ProfileDatabase()

	if err := http.ListenAndServe(":8082", nil); err != nil {
		panic(err)
	}
}
