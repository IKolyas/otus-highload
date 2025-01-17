package repository

import (
	"database/sql"
	"errors"

	"github.com/IKolyas/otus-highload/internal/domain"
)

type UserRepository struct {
	Connection *sql.DB
}

func (r *UserRepository) GetAuthData(login string) (*domain.User, error) {
	if r.Connection == nil {
		return nil, errors.New("database connection is nil")
	}

	user := domain.User{}

	row := "SELECT id, login, password FROM users WHERE login = $1"

	err := r.Connection.QueryRow(row, login).Scan(&user.ID, &user.Login, &user.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetByID(id int) (*domain.User, error) {

	if r.Connection == nil {
		return nil, errors.New("database connection is nil")
	}

	user := domain.User{}

	row := "SELECT id, login, first_name, second_name, birthdate, biography, city FROM users WHERE id = $1"

	err := r.Connection.QueryRow(row, id).Scan(&user.ID, &user.Login, &user.FirstName, &user.SecondName, &user.Birthdate, &user.Biography, &user.City)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetBy(field string, value interface{}) (*domain.User, error) {

	if r.Connection == nil {
		return nil, errors.New("database connection is nil")
	}

	user := domain.User{}

	row := "SELECT id, login, first_name, second_name, birthdate, biography, city FROM users WHERE " + field + " = $1"

	err := r.Connection.QueryRow(row, value).Scan(&user.ID, &user.Login, &user.FirstName, &user.SecondName, &user.Birthdate, &user.Biography, &user.City)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) Create(user *domain.User) error {

	if r.Connection == nil {
		return errors.New("DB connection error")
	}

	row := "INSERT INTO users (login, password, first_name, second_name, birthdate, biography, city) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id"

	lastInsert := 0
	err := r.Connection.QueryRow(row, user.Login, user.Password, user.FirstName, user.SecondName, user.Birthdate, user.Biography, user.City).Scan(&lastInsert)
	if err != nil {
		return errors.New("DB error: " + err.Error())
	}

	user.ID = int(lastInsert)

	return nil
}
