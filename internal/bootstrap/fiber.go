package bootstrap

import (
	"github.com/gofiber/fiber/v3"
)

func NewFiber() *fiber.App {
	return fiber.New(fiber.Config{
		AppName:      "Glamping API",
		ErrorHandler: customErrorHandler,
	})
}

func customErrorHandler(c fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}
	return c.Status(code).JSON(fiber.Map{
		"success": false,
		"error":   err.Error(),
	})
}
