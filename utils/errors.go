package utils

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type HttpError struct {
	Code    int
	Message string
}

func (e *HttpError) Error() string {
	return e.Message
}

func BadRequest(messages ...string) *HttpError {
	message := "Bad Request"
	if len(messages) > 0 {
		message = messages[0]
	}
	return &HttpError{
		Code:    400,
		Message: message,
	}
}

func Forbidden(messages ...string) *HttpError {
	message := "您没有权限进行此操作"
	if len(messages) > 0 {
		message = messages[0]
	}
	return &HttpError{
		Code:    403,
		Message: message,
	}
}

func MyErrorHandler(ctx *fiber.Ctx, err error) error {
	if err == nil {
		return nil
	}

	code := 500
	message := err.Error()

	if e, ok := err.(*HttpError); ok {
		code = e.Code
	} else if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		code = 404
	}

	return ctx.Status(code).JSON(fiber.Map{"message": message})
}
