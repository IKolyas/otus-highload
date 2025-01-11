package main

import (
	"log"

	"github.com/IKolyas/otus-highload/config"
	"github.com/IKolyas/otus-highload/internal/infrastructure"
)

func main() {

	appConfig := config.AppConfig{}
	appConfig.Load()

	pgsqlConfig := config.PgsqlConfig{}
	pgsqlConfig.Load()

	_, err := infrastructure.PgsqlConnection.NewConnection(pgsqlConfig)
	if err != nil {
		panic(err)
	}

	router := infrastructure.Router()

	router.Listen(":" + appConfig.AppPort)

	log.Fatal(router)
}
