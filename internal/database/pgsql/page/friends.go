package page

import "fmt"

func (s *homepageRepository) CheckUserID(response string, args ...any) string {
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
