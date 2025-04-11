package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type User struct {
	User_id      int
	Name_user    string
	Surname_user string
	City_user    string
}

type Track struct {
	Track_id    int
	Name_music  string
	Name_artist string
}

type LoginAndPassword struct {
	Login    string
	Password string
}

var Users = []User{}
var Tracks = []Track{}

var DB *sql.DB

func InitDatabase() error {

	s, _ := os.LookupEnv("INIT_DB")

	connStr := s
	var err error

	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	DB.SetMaxOpenConns(25)    // max соединений
	DB.SetConnMaxIdleTime(25) // max бездействующих соединений
	DB.SetConnMaxLifetime(3 * time.Minute)

	return DB.Ping()
}

func Close() error {
	return DB.Close()
}

func ProfileDatabase() {

	row, err := DB.Query("SELECT * FROM users")
	if err != nil {
		panic(err)
	}
	defer row.Close()

	for row.Next() {
		u := User{}
		err := row.Scan(&u.User_id, &u.Name_user, &u.Surname_user, &u.City_user)
		if err != nil {
			fmt.Println(err)
			continue
		}
		Users = append(Users, u)
	}

	row, err = DB.Query("SELECT * FROM tracks")
	if err != nil {
		panic(err)
	}
	defer row.Close()

	for row.Next() {
		t := Track{}
		err := row.Scan(&t.Track_id, &t.Name_music, &t.Name_artist)
		if err != nil {
			fmt.Println(err)
			continue
		}
		Tracks = append(Tracks, t)
	}

}

func InsertResponseDatabase(response string, args ...any) {

	DB.Exec(response, args...)

}

func CheckLoginDatabase(response string, args ...any) string {
	row := DB.QueryRow(response, args...)
	l := struct{ login string }{}
	err := row.Scan(&l.login)
	if err != nil {
		fmt.Println("данные не были получены")
		return ""
	}
	return l.login
}

func SelectLoginOrPasswordOnDatabase(login string) *LoginAndPassword {

	row := DB.QueryRow("SELECT login,password FROM users WHERE login=$1", login)

	lp := &LoginAndPassword{}

	err := row.Scan(&lp.Login, &lp.Password)
	if err != nil {
		fmt.Println("данные не были получены")
		return nil
	}

	return lp
}

func SelectUser(user_id int) string {
	row := DB.QueryRow("SELECT username FROM users WHERE user_id = $1", user_id)

	myuser := struct {
		Username string
	}{}

	err := row.Scan(&myuser.Username)
	if err != nil {
		fmt.Println("Данные не были получены")
		return ""
	}

	return myuser.Username

}

type TopTracksUser struct {
	Place      string
	NameMusic  string
	NameArtist string
}

func GetTopTracksUser(user_id int) []TopTracksUser {

	table := []TopTracksUser{}

	rows, err := DB.Query("SELECT place, name_music, name_artist FROM top_tracks_user JOIN tracks USING(track_id) WHERE user_id=$1 ORDER BY place", user_id)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		temp := TopTracksUser{}

		err := rows.Scan(&temp.Place, &temp.NameMusic, &temp.NameArtist)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		table = append(table, temp)
	}
	fmt.Println(table)
	return table
}
