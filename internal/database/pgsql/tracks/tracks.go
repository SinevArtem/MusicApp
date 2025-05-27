package tracks

import (
	"LoveMusic/internal/database/pgsql/factory"
	"fmt"
	"log"
)

type TracksRepository interface {
	GetTracks(track, artist string) []Tracks
	AddTrack(name_music, name_artist string)
	CheckTrackAndArtist(args ...any) (string, string)
	Insert(response string, args ...any) error
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
