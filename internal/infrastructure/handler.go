package infrastructure

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/IKolyas/otus-highload/internal/application"
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

var userUC = application.UserUseCase{}

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
		Login    string `validate:"required"`
		Password string `validate:"required"`
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
	_, err = userUC.Login(u.Login, u.Password, &UserRepository{Connection: PgsqlConnection.Connection})
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
	u := new(struct {
		Login    string `validate:"required"`
		Password string `validate:"required"`
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
	err = userUC.Register(u.Login, u.Password, &UserRepository{Connection: PgsqlConnection.Connection})
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{"status": "success"})
}

func GetUserHanlder(c *fiber.Ctx) error {

	auth := c.Locals("user").(*jwt.Token)
	_ = auth.Claims.(jwt.MapClaims)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	user, err := userUC.GetByID(id, &UserRepository{Connection: PgsqlConnection.Connection})
	if err != nil {
		return err
	}

	return c.JSON(user)
}
