package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

type LoginAndPassword struct {
	Login    string
	Password string
}

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

	return table
}

func GetUserID(login string) int {

	row := DB.QueryRow("SELECT user_id FROM users WHERE login=$1", login)

	user_id := struct{ user_id string }{}

	err := row.Scan(&user_id.user_id)
	if err != nil {
		log.Println(err)
	}

	t, _ := strconv.Atoi(user_id.user_id)
	return t
}

func AddTrack(name_music, name_artist string) {
	_, err := DB.Exec("INSERT INTO tracks (name_music, name_artist) VALUES ($1,$2)", name_music, name_artist)
	if err != nil {
		log.Println(err)
		return
	}

}

func CheckTrackAndArtist(response string, args ...any) (string, string) {
	row := DB.QueryRow(response, args...)
	l := struct {
		track  string
		artist string
	}{}
	err := row.Scan(&l.track, &l.artist)
	if err != nil {
		fmt.Println("данные не были получены")

	}
	return l.track, l.artist
}

func CheckUserID(response string, args ...any) string {
	row := DB.QueryRow(response, args...)
	l := struct {
		user_id string
	}{}
	err := row.Scan(&l.user_id)
	if err != nil {
		fmt.Println("данные не были получены")

	}
	return l.user_id
}
