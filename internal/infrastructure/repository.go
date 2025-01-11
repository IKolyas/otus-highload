package infrastructure

import (
	"database/sql"
	"errors"

	"github.com/IKolyas/otus-highload/internal/domain"
)

type UserRepository struct {
	Connection *sql.DB
}

func (r *UserRepository) GetByID(id int) (*domain.User, error) {

	if r.Connection == nil {
		return nil, errors.New("database connection is nil")
	}

	user := domain.User{}

	err := r.Connection.QueryRow("SELECT first_name, second_name, birthdate, biography, city FROM users WHERE id = $1", id).Scan(&user.FirstName, &user.SecondName, &user.Birthdate, &user.Biography, &user.City)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetBy(field string, value interface{}) (*domain.User, error) {

	if r.Connection == nil {
		return nil, errors.New("database connection is nil")
	}

	user := domain.User{}

	err := r.Connection.QueryRow("SELECT * FROM users WHERE "+field+" = $1", value).Scan(&user.ID, &user.Login, &user.Password, &user.FirstName, &user.SecondName, &user.Birthdate, &user.Biography, &user.City)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) Create(login string, password string) error {

	if r.Connection == nil {
		return errors.New("database connection is nil")
	}

	_, err := r.Connection.Exec("INSERT INTO users (login, password) VALUES ($1, $2)", login, password)
	if err != nil {
		return err
	}
	return nil
}
