package application

import (
	"errors"

	"github.com/IKolyas/otus-highload/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type Token string

type UserRepository interface {
	GetByID(id int) (*domain.User, error)
	GetBy(field string, value interface{}) (*domain.User, error)
	Create(login string, password string) error
}

type UserUseCase struct {
}

func (u *UserUseCase) Register(login string, password string, r UserRepository) error {
	passwd, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	return r.Create(login, string(passwd))
}

func (u *UserUseCase) Login(login string, password string, r UserRepository) (*domain.User, error) {

	user, err := r.GetBy("login", login)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("incorrect password")
	}

	return user, nil
}

func (u *UserUseCase) GetByID(id int, r UserRepository) (map[string]interface{}, error) {

	user, err := r.GetByID(id)
	if err != nil {
		res := map[string]interface{}{}
		return res, err
	}

	res := map[string]interface{}{
		"FirstName":  user.FirstName,
		"SecondName": user.SecondName,
		"Birthdate":  user.Birthdate,
		"Biography":  user.Biography,
		"City":       user.City}

	return res, nil
}
