package handlers

import (
	ct "LoveMusic/internal/create_templates"
	db "LoveMusic/internal/database"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func LoadProfile(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := Authorise(w, r); err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	v := ct.GetChartUser()
	tmpl, _ := template.ParseFiles("static/templates/profile.html")
	tmpl.Execute(w, v)

	//http.ServeFile(w, r, "static/templates/profile.html")

}

func UserFriends(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := Authorise(w, r); err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	v := ct.GetChartUser()
	tmpl, _ := template.ParseFiles("static/templates/friends.html")
	tmpl.Execute(w, v)
}

func UserProfileHandler(w http.ResponseWriter, r *http.Request) {
	if err := Authorise(w, r); err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	myurl := strings.Split(r.URL.Path, "/")
	if len(myurl) < 3 || myurl[2] == "" {
		fmt.Println("f")
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	user_id, _ := strconv.Atoi(myurl[2])

	username := db.SelectUser(user_id)
	if username == "" {
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}

	viewprofile := struct {
		Username string
		User_ID  int
	}{
		Username: username,
		User_ID:  user_id,
	}

	topTracks := db.GetTopTracksUser(user_id)

	data := struct {
		Profile struct {
			Username string
			User_ID  int
		}
		TopTracks []db.TopTracksUser
	}{
		Profile:   viewprofile,
		TopTracks: topTracks,
	}

	tmpl, err := template.ParseFiles("static/templates/user_profile.html")
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(data)

	if err := tmpl.Execute(w, data); err != nil {
		log.Println(err)
		return
	}
	fmt.Println(username)

}
