package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type (
	AppConfig struct {
		AppPort string
	}

	PgsqlConfig struct {
		Host     string
		Port     string
		User     string
		Password string
		Dbname   string
	}
)

func init() {
	err := godotenv.Load("/app/.env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}
}

func (c *AppConfig) Load() *AppConfig {
	c.AppPort = os.Getenv("APP_PORT")
	return c
}

func (c *PgsqlConfig) Load() *PgsqlConfig {
	c.Host = os.Getenv("PGSQL_HOST")
	c.Port = os.Getenv("PGSQL_PORT")
	c.User = os.Getenv("PGSQL_USER")
	c.Password = os.Getenv("PGSQL_PASSWORD")
	c.Dbname = os.Getenv("PGSQL_DB")
	return c

}
