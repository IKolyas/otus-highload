package domain

type User struct {
	ID         string  `json:"id"`
	Login      string  `json:"login"`
	Password   string  `json:"password"`
	FirstName  *string `json:"firstName"`
	SecondName *string `json:"secondName"`
	Birthdate  *string `json:"birthdate"`
	Biography  *string `json:"biography"`
	City       *string `json:"city"`
}
