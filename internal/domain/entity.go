package domain

type User struct {
	ID         int    `json:"id"`
	Login      string `json:"login"`
	Password   string `json:"password"`
	FirstName  string `json:"firstName"`
	SecondName string `json:"secondName"`
	Birthdate  string `json:"birthdate"`
	Biography  string `json:"biography"`
	City       string `json:"city"`
}

func (u *User) RequestData() map[string]interface{} {
	return map[string]interface{}{
		"id":         u.ID,
		"login":      u.Login,
		"firstName":  u.FirstName,
		"secondName": u.SecondName,
		"birthdate":  u.Birthdate,
		"biography":  u.Biography,
		"city":       u.City,
	}
}
