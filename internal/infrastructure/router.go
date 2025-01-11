package infrastructure

import (
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
)

func Router() *fiber.App {

	router := fiber.New()
	router.Post("/login", LoginHandler)
	router.Post("/register", RegisterHandler)

	router.Get("/metrics", monitor.New(monitor.Config{Title: "MyService Metrics Page"}))
	router.Use(logger.New())
	router.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte("секрет")},
	}))

	router.Get("/users/:id", GetUserHanlder)

	return router
}
