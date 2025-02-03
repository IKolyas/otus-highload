package database

import (
	"database/sql"
	"fmt"

	"github.com/IKolyas/otus-highload/config"
	_ "github.com/lib/pq"
)

type Pgsql struct {
	Connection *sql.DB
}

var PgConnection Pgsql = Pgsql{
	Connection: nil,
}

func (p *Pgsql) NewConnection(config config.PgsqlConfig) (bool, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", config.Host, config.Port, config.User, config.Password, config.Dbname)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return false, err
	}
	if err = db.Ping(); err != nil {
		return false, err
	}

	p.Connection = db

	return true, nil
}
