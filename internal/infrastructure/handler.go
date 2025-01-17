package infrastructure

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/IKolyas/otus-highload/internal/application/service"
	"github.com/IKolyas/otus-highload/internal/domain"
	"github.com/IKolyas/otus-highload/internal/infrastructure/repository"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type (
	ErrorResponse struct {
		Error       bool
		FailedField string
		Tag         string
		Value       interface{}
	}

	XValidator struct {
		validator *validator.Validate
	}

	GlobalErrorHandlerResp struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}
)

var xv = &XValidator{
	validator: validator.New(),
}

func (v XValidator) Validate(data interface{}) error {
	validationErrors := []ErrorResponse{}

	errs := xv.validator.Struct(data)
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			var elem ErrorResponse

			elem.FailedField = err.Field()
			elem.Tag = err.Tag()
			elem.Value = err.Value()
			elem.Error = true

			validationErrors = append(validationErrors, elem)
		}
	}

	if len(validationErrors) > 0 {
		errMsgs := make([]string, 0)

		for _, err := range validationErrors {
			errMsgs = append(errMsgs, fmt.Sprintf(
				"[%s]: '%v' | need to be realized '%s'",
				err.FailedField,
				err.Value,
				err.Tag,
			))
		}

		return &fiber.Error{
			Code:    fiber.ErrBadRequest.Code,
			Message: strings.Join(errMsgs, " and "),
		}
	}

	return nil
}

func LoginHandler(c *fiber.Ctx) error {
	u := new(struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	})
	if err := c.BodyParser(u); err != nil {
		return err
	}

	user := struct {
		Login    string `validate:"required"`
		Password string `validate:"required"`
	}{Login: u.Login, Password: u.Password}

	err := xv.Validate(user)
	if err != nil {
		return err
	}
	uService := service.User{}
	r := repository.UserRepository{Connection: PgsqlConnection.Connection}
	_, err = uService.Login(u.Login, u.Password, &r)
	if err != nil {
		return err
	}

	claims := jwt.MapClaims{
		"name":  u.Login,
		"admin": true,
		"exp":   time.Now().Add(time.Hour * 1).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte("секрет"))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"token": t})
}

func RegisterHandler(c *fiber.Ctx) error {
	body := new(struct {
		Login      string `json:"login"`
		Password   string `json:"password"`
		FirstName  string `json:"firstName"`
		SecondName string `json:"secondName"`
		Birthdate  string `json:"birthdate"`
		Biography  string `json:"biography"`
		City       string `json:"city"`
	})

	if err := c.BodyParser(body); err != nil {
		c.SendStatus(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"status": "error", "message": err})
	}

	request := struct {
		Login      string `validate:"required,min=6"`
		Password   string `validate:"required,min=8"`
		FirstName  string `validate:"required,min=2"`
		SecondName string `validate:"required,min=2"`
		Birthdate  string `validate:"required"`
		Biography  string `validate:"required"`
		City       string `validate:"required,min=2"`
	}{Login: body.Login, Password: body.Password, FirstName: body.FirstName, SecondName: body.SecondName, Birthdate: body.Birthdate, Biography: body.Biography, City: body.City}

	if err := xv.Validate(request); err != nil {
		c.SendStatus(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"status": "error", "message": err})
	}

	if _, err := time.Parse("01.02.2006", request.Birthdate); err != nil {
		c.SendStatus(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"status": "error", "message": errors.New("Invalid Birthdate field")})
	}

	uService := service.User{}
	user := domain.User{
		Login:      request.Login,
		Password:   request.Password,
		FirstName:  request.FirstName,
		SecondName: request.SecondName,
		Birthdate:  request.Birthdate,
		Biography:  request.Biography,
		City:       request.City,
	}

	conn := PgsqlConnection.Connection

	if err := uService.Register(&user, &repository.UserRepository{Connection: conn}); err != nil {
		c.SendStatus(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"status": "error", "message": err})
	}

	return c.JSON(fiber.Map{"status": "success"})
}

func GetUserHanlder(c *fiber.Ctx) error {

	auth := c.Locals("user").(*jwt.Token)
	_ = auth.Claims.(jwt.MapClaims)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		c.SendStatus(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"status": "error", "message": err})
	}

	uService := service.User{}

	profile, err := uService.GetByID(id, &repository.UserRepository{Connection: PgsqlConnection.Connection})

	return c.JSON(fiber.Map{"status": "success", "data": profile})
}
