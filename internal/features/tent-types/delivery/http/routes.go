package http

import (
	"github.com/gofiber/fiber/v3"
	"github.com/ramdhanrizkij/nest-glamping-api/internal/shared/middleware"
)

func RegisterRoutes(router fiber.Router, handler *Handler, jwtSecret string) {
	tt := router.Group("/tent-types")
	tt.Get("", handler.List)
	tt.Get("/:id", handler.FindByID)

	// Images (nested)
	tt.Post("/:id/images", middleware.Auth(jwtSecret), middleware.RoleAllowed("admin"), handler.AddImage)
	tt.Delete("/:id/images/:imageId", middleware.Auth(jwtSecret), middleware.RoleAllowed("admin"), handler.DeleteImage)
	tt.Put("/:id/images/:imageId/primary", middleware.Auth(jwtSecret), middleware.RoleAllowed("admin"), handler.SetPrimaryImage)

	// Rates (nested)
	tt.Get("/:id/rates", handler.ListRates)
	tt.Post("/:id/rates", middleware.Auth(jwtSecret), middleware.RoleAllowed("admin"), handler.CreateRate)
	tt.Put("/:id/rates/:rateId", middleware.Auth(jwtSecret), middleware.RoleAllowed("admin"), handler.UpdateRate)
	tt.Delete("/:id/rates/:rateId", middleware.Auth(jwtSecret), middleware.RoleAllowed("admin"), handler.DeleteRate)

	// Admin CRUD
	admin := tt.Group("", middleware.Auth(jwtSecret), middleware.RoleAllowed("admin"))
	admin.Post("", handler.Create)
	admin.Put("/:id", handler.Update)
	admin.Delete("/:id", handler.Delete)
}
