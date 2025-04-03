package main

import (
	db "LoveMusic/internal/database"
	h "LoveMusic/internal/handlers"
	"net/http"
)

func main() {
	fileserver := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fileserver))

	http.HandleFunc("/", h.LoadProfile)
	http.HandleFunc("/register", h.RegisterHandler)
	http.HandleFunc("/login", h.LoginHandler)

	//http.HandleFunc("/profile", LoadProfile)

	db.OpenDatabase()

	if err := http.ListenAndServe(":8082", nil); err != nil {
		panic(err)
	}
}
