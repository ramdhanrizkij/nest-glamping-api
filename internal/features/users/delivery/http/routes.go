package http

import (
	"github.com/gofiber/fiber/v3"
	"github.com/ramdhanrizkij/nest-glamping-api/internal/shared/middleware"
)

func RegisterRoutes(router fiber.Router, handler *Handler, jwtSecret string) {
	users := router.Group("/users", middleware.Auth(jwtSecret))
	users.Get("/profile", handler.GetProfile)
	users.Put("/profile", handler.UpdateProfile)
}
