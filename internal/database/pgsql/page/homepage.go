package page

import (
	"LoveMusic/internal/database/pgsql/factory"
	"fmt"
	"log"
	"strconv"
)

type HomepageRepository interface {
	GetUserID(login string) int
	SelectUser(user_id int) string
	GetTopTracksUser(user_id int) []TopTracksUser
	CheckUserID(response string, args ...any) string
}

type TopTracksUser struct {
	Place      string
	NameMusic  string
	NameArtist string
}

type homepageRepository struct {
	*factory.Storage
}

func NewHomepageRepository(db *factory.Storage) HomepageRepository {
	return &homepageRepository{Storage: db}
}

func (s *homepageRepository) GetUserID(login string) int {

	row := s.DB.QueryRow("SELECT user_id FROM users WHERE login=$1", login)

	user_id := struct{ user_id string }{}

	err := row.Scan(&user_id.user_id)
	if err != nil {
		log.Println(err)
	}

	t, _ := strconv.Atoi(user_id.user_id)
	return t
}

func (s *homepageRepository) SelectUser(user_id int) string {
	row := s.DB.QueryRow("SELECT username FROM users WHERE user_id = $1", user_id)

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

func (s *homepageRepository) GetTopTracksUser(user_id int) []TopTracksUser {

	table := []TopTracksUser{}

	rows, err := s.DB.Query("SELECT place, name_music, name_artist FROM top_tracks_user JOIN tracks USING(track_id) WHERE user_id=$1 ORDER BY place", user_id)
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
