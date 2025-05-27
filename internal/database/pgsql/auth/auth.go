package auth

import (
	"LoveMusic/internal/database/pgsql/factory"
	"fmt"
)

type AuthRepository interface {
	InsertRegisterValue(username, login, password_hash string) error
	CheckLoginDatabase(login string) string
	SelectLoginOrPasswordOnDatabase(login string) (string, string)
}

type authRepository struct {
	*factory.Storage
}

func NewAuthRepository(db *factory.Storage) AuthRepository {
	return &authRepository{Storage: db}
}

func (s *authRepository) InsertRegisterValue(username, login, password_hash string) error {
	if err := s.Insert("INSERT INTO users (username, login, password) VALUES ($1, $2, $3);", username, login, password_hash); err != nil {
		return fmt.Errorf("error insert register value: %w", err)
	}
	return nil
}

func (s *authRepository) CheckLoginDatabase(login string) string {
	row := s.DB.QueryRow("SELECT login FROM users WHERE login=$1", login)
	l := struct{ login string }{}
	err := row.Scan(&l.login)
	if err != nil {
		fmt.Println("данные не были получены")
		return ""
	}
	return l.login
}

func (s *authRepository) SelectLoginOrPasswordOnDatabase(login string) (string, string) {

	row := s.DB.QueryRow("SELECT login,password FROM users WHERE login=$1", login)

	lp := struct {
		Login    string
		Password string
	}{}

	err := row.Scan(&lp.Login, &lp.Password)
	if err != nil {
		fmt.Println("данные не были получены")
		return "", ""
	}

	return lp.Login, lp.Password
}
