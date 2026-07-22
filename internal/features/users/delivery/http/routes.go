package http

import (
	"github.com/gofiber/fiber/v3"
	"github.com/ramdhanrizkij/nest-glamping-api/internal/shared/middleware"
)

func RegisterRoutes(router fiber.Router, handler *Handler, jwtSecret string) {
	users := router.Group("/users", middleware.Auth(jwtSecret))
	users.Get("/profile", handler.GetProfile)
	users.Put("/profile", handler.UpdateProfile)

	admin := router.Group("/admin/users", middleware.Auth(jwtSecret), middleware.RoleAllowed("admin"))
	admin.Get("", handler.ListAll)
	admin.Get("/:id", handler.GetUserByID)
	admin.Put("/:id", handler.UpdateUser)
	admin.Delete("/:id", handler.DeleteUser)
}
