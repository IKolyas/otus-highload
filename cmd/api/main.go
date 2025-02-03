package main

import (
	"log"

	"github.com/IKolyas/otus-highload/config"
	"github.com/IKolyas/otus-highload/internal/infrastructure"
	"github.com/IKolyas/otus-highload/internal/infrastructure/database"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load("/app/.env")
	if err != nil {
		log.Println("Error loading .env file")
	}

	appConfig := config.AppConfig{}
	appConfig.Load()

	pgsqlConfig := config.PgsqlConfig{}
	pgsqlConfig.Load()

	database.PgConnection.NewConnection(pgsqlConfig)

	if err != nil {
		panic(err)
	}

	router := infrastructure.Router()

	router.Listen(":" + appConfig.AppPort)

	log.Fatal(router)
}
