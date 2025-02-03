package repository

import (
	"database/sql"
	"errors"
	"strconv"
	"strings"

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

	row := "SELECT id, login, first_name, second_name, gender, birthdate, biography, city FROM users WHERE id = $1"

	err := r.Connection.QueryRow(row, id).Scan(&user.ID, &user.Login, &user.FirstName, &user.SecondName, &user.Gender, &user.Birthdate, &user.Biography, &user.City)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) Find(fields map[string]string) ([]domain.User, error) {
	if r.Connection == nil {
		return nil, errors.New("database connection is nil")
	}

	// Строим динамический запрос
	query := "SELECT id, login, first_name, second_name, gender, birthdate, biography, city FROM users WHERE "
	var conditions []string
	var args []interface{}

	count := 1
	for field, value := range fields {
		conditions = append(conditions, field+" LIKE $"+strconv.Itoa(count))
		args = append(args, value+"%")
		count++
	}

	query += strings.Join(conditions, " AND ")

	// Добавляем сортировку
	query += " ORDER BY id"

	rows, err := r.Connection.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		user := domain.User{}
		if err := rows.Scan(&user.ID, &user.Login, &user.FirstName, &user.SecondName, &user.Gender, &user.Birthdate, &user.Biography, &user.City); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) Save(user *domain.User) (res int, err error) {
	if r.Connection == nil {
		return 0, errors.New("DB connection error")
	}

	row := "INSERT INTO users (login, password, first_name, second_name, gender, birthdate, biography, city) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id"

	lastInsert := 0
	err = r.Connection.QueryRow(row, user.Login, user.Password, user.FirstName, user.SecondName, &user.Gender, user.Birthdate, user.Biography, user.City).Scan(&lastInsert)
	if err != nil {
		return 0, errors.New("DB error: " + err.Error())
	}

	user.ID = int(lastInsert)

	return user.ID, nil
}
