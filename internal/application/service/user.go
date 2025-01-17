package service

import (
	"errors"

	"github.com/IKolyas/otus-highload/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type User struct{}

func (u *User) Register(user *domain.User, r domain.Repository[domain.User]) error {
	passwd, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		return errors.New("Generate password error")
	}

	user.Password = string(passwd)

	if err := r.Create(user); err != nil {
		return err
	}

	return nil
}

func (u *User) Login(login string, password string, r domain.UserRepository) (*domain.User, error) {

	user, err := r.GetAuthData(login)
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

func (u *User) GetByID(id int, r domain.Repository[domain.User]) (map[string]interface{}, error) {

	user, err := r.GetByID(id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	return user.RequestData(), nil
}
