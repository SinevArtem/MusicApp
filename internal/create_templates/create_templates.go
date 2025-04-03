package createtemplates

import (
	db "LoveMusic/internal/database"
	// 	"html/template"
)

type viewprofile struct {
	User_id      int
	Name_user    string
	Surname_user string
	City_user    string

	Track_id    int
	Name_music  string
	Name_artist string
}

type exeptionOnRegOrLog struct {
	Exeption string
}

func GetChartUser() *viewprofile {
	ViewProfile := viewprofile{}
	for _, u := range db.Users {
		if u.User_id == 1 {
			ViewProfile.User_id = u.User_id
			ViewProfile.Name_user = u.Name_user
			ViewProfile.Surname_user = u.Surname_user
			ViewProfile.City_user = u.City_user
		}
	}

	for _, t := range db.Tracks {
		if t.Track_id == 1 {
			ViewProfile.Track_id = t.Track_id
			ViewProfile.Name_music = t.Name_music
			ViewProfile.Name_artist = t.Name_artist
		}
	}

	return &ViewProfile
}

func GetExeptionOnRegister(s string) *exeptionOnRegOrLog {

	exeptionOnRegOrLog := &exeptionOnRegOrLog{Exeption: s}
	return exeptionOnRegOrLog
}
