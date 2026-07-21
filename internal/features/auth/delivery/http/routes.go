package http

import "github.com/gofiber/fiber/v3"

func RegisterRoutes(router fiber.Router, handler *Handler) {
	auth := router.Group("/auth")
	auth.Post("/register", handler.Register)
	auth.Post("/login", handler.Login)
	auth.Post("/refresh", handler.RefreshToken)
	auth.Post("/logout", handler.Logout)
}
