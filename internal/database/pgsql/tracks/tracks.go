package tracks

import (
	"LoveMusic/internal/database/pgsql/factory"
	"fmt"
	"log"
	"strconv"
)

type TracksRepository interface {
	GetTracks(track, artist string) []Tracks
	AddTrack(name_music, name_artist string)
	CheckTrackAndArtist(args ...any) (string, string)
	Insert(response string, args ...any) error
	GetUserID(login string) int
	GetTrackID(name_music, name_artist string) int
	AddToCollection(user_id, track_id int) error
}

type Tracks struct {
	NameMusic  string
	NameArtist string
}

type tracksRepository struct {
	*factory.Storage
}

func NewTracksRepository(db *factory.Storage) TracksRepository {
	return &tracksRepository{Storage: db}
}

func (s *tracksRepository) GetTracks(track, artist string) []Tracks {

	table := []Tracks{}

	rows, err := s.DB.Query("SELECT name_music, name_artist FROM tracks WHERE name_music=$1 AND name_artist=$2", track, artist)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		temp := Tracks{}

		err := rows.Scan(&temp.NameMusic, &temp.NameArtist)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		table = append(table, temp)
	}

	return table
}

func (s *tracksRepository) AddTrack(name_music, name_artist string) {
	_, err := s.DB.Exec("INSERT INTO tracks (name_music, name_artist) VALUES ($1,$2)", name_music, name_artist)
	if err != nil {
		log.Println(err)
		return
	}

}

func (s *tracksRepository) CheckTrackAndArtist(args ...any) (string, string) {
	row := s.DB.QueryRow("SELECT name_music, name_artist FROM tracks WHERE name_music=$1 AND name_artist=$2", args...)
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

func (s *tracksRepository) CheckUserID(response string, args ...any) string {
	row := s.DB.QueryRow(response, args...)
	l := struct {
		user_id string
	}{}
	err := row.Scan(&l.user_id)
	if err != nil {
		fmt.Println("данные не были получены")

	}
	return l.user_id
}

func (s *tracksRepository) GetTrackID(name_music, name_artist string) int {
	row := s.DB.QueryRow("SELECT track_id FROM tracks WHERE name_music=$1 AND name_artist=$2", name_music, name_artist)
	track_id := struct{ track_id string }{}

	err := row.Scan(&track_id.track_id)
	if err != nil {
		log.Println(err)
	}

	t, _ := strconv.Atoi(track_id.track_id)
	return t
}

func (s *tracksRepository) AddToCollection(user_id, track_id int) error {
	_, err := s.DB.Exec("INSERT INTO top_tracks_user (user_id, track_id, place) VALUES ($1,$2, COALESCE((SELECT MAX(place) FROM top_tracks_user WHERE user_id=$1), 0) + 1)", user_id, track_id)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (s *tracksRepository) GetUserID(login string) int {

	row := s.DB.QueryRow("SELECT user_id FROM users WHERE login=$1", login)

	user_id := struct{ user_id string }{}

	err := row.Scan(&user_id.user_id)
	if err != nil {
		log.Println(err)
	}

	t, _ := strconv.Atoi(user_id.user_id)
	return t
}
