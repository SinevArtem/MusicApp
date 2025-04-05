package database

import (
	"database/sql"
	"fmt"
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

	for _, u := range Users {
		fmt.Println(u.User_id, u.Name_user, u.Surname_user, u.City_user)
	}

	for _, t := range Tracks {
		fmt.Println(t.Track_id, t.Name_music, t.Name_artist)
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
