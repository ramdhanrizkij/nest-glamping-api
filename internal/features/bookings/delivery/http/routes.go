package http

import (
	"github.com/gofiber/fiber/v3"
	"github.com/ramdhanrizkij/nest-glamping-api/internal/shared/middleware"
)

func RegisterRoutes(router fiber.Router, handler *Handler, jwtSecret string) {
	// Customer bookings
	bookings := router.Group("/bookings", middleware.Auth(jwtSecret))
	bookings.Post("", handler.CreateBooking)
	bookings.Get("", handler.ListMyBookings)
	bookings.Get("/:id", handler.GetBookingDetail)
	bookings.Patch("/:id/cancel", handler.CancelBooking)

	// Admin bookings
	admin := router.Group("/admin/bookings", middleware.Auth(jwtSecret), middleware.RoleAllowed("admin"))
	admin.Get("", handler.ListAllBookings)
	admin.Get("/:id", handler.GetBookingDetailAdmin)
	admin.Patch("/:id/confirm", handler.ConfirmBooking)
}
