package domain

type Repository[T any] interface {
	GetByID(id int) (*T, error)
	GetBy(field string, value interface{}) (*T, error)
	Create(*T) error
}

type UserRepository interface {
	Repository[User]
	GetAuthData(login string) (*User, error)
}
