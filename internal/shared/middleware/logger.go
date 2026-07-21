package middleware

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
)

func Logger() fiber.Handler {
	return func(c fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		log.Info(c.Method(), " ", c.Path(), " ", c.Response().StatusCode(), " ", time.Since(start))
		return err
	}
}
