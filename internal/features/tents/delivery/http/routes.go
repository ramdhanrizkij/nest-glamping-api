package http

import (
	"github.com/gofiber/fiber/v3"
	"github.com/ramdhanrizkij/nest-glamping-api/internal/shared/middleware"
)

func RegisterRoutes(router fiber.Router, handler *Handler, jwtSecret string) {
	tents := router.Group("/tents")

	// Public
	tents.Get("/tent-types/:id/availability", handler.CheckAvailability)

	// Admin
	admin := tents.Group("", middleware.Auth(jwtSecret), middleware.RoleAllowed("admin"))
	admin.Get("/units", handler.List)
	admin.Post("/units", handler.Create)
	admin.Get("/units/:id", handler.FindByID)
	admin.Put("/units/:id", handler.Update)
	admin.Delete("/units/:id", handler.Delete)
	admin.Get("/tent-types/:id/units", handler.ListByTentTypeID)
}
