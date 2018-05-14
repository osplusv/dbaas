package handlers

import (
	"github.com/labstack/echo"
	validator "gopkg.in/go-playground/validator.v9"
)

var validate *validator.Validate

type Response struct {
	Data   interface{} `json:"data,omitempty"`
	Status int         `json:"status"`
	Error  *Error      `json:"error,omitempty"`
}

type Error struct {
	Message string `json:"message,omitempty"`
}

func bindPayload(payload interface{}, c echo.Context) error {
	if err := c.Bind(payload); err != nil {
		return err
	}
	return doValidation(payload)
}

func doValidation(payload interface{}) error {
	validate = validator.New()
	return validate.Struct(payload)
}

func respondWithError(httpStatus int, message string, c echo.Context) error {
	return c.JSON(httpStatus, Response{Status: httpStatus, Error: &Error{Message: message}})
}

func respondWithPayload(httpStatus int, payload interface{}, c echo.Context) error {
	return c.JSON(httpStatus, Response{Data: payload, Status: httpStatus, Error: nil})
}

func responseNoContent(httpStatus int, c echo.Context) error {
	return c.NoContent(httpStatus)
}
