package http

import (
	"github.com/gofiber/fiber/v3"
	"github.com/ramdhanrizkij/nest-glamping-api/internal/shared/middleware"
)

func RegisterRoutes(router fiber.Router, handler *Handler, jwtSecret string) {
	amenities := router.Group("/amenities")
	amenities.Get("", handler.List)
	amenities.Get("/:id", handler.FindByID)

	admin := amenities.Group("", middleware.Auth(jwtSecret), middleware.RoleAllowed("admin"))
	admin.Post("", handler.Create)
	admin.Put("/:id", handler.Update)
	admin.Delete("/:id", handler.Delete)
}
