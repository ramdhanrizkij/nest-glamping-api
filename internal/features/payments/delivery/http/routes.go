package http

import (
	"github.com/gofiber/fiber/v3"
	"github.com/ramdhanrizkij/nest-glamping-api/internal/shared/middleware"
)

func RegisterRoutes(router fiber.Router, handler *Handler, jwtSecret string) {
	payments := router.Group("/bookings/:id", middleware.Auth(jwtSecret))
	payments.Post("/pay", handler.InitiatePayment)
	payments.Get("/payment", handler.GetPaymentStatus)

	router.Post("/payments/:id/callback", handler.HandleCallback)
}
