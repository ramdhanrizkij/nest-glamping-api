package bootstrap

import (
	"github.com/gofiber/fiber/v3"
	"github.com/ramdhanrizkij/nest-glamping-api/internal/shared/middleware"
	"github.com/ramdhanrizkij/nest-glamping-api/internal/shared/swagger"
)

func SetupRoutes(app *fiber.App, deps *Dependencies) {
	app.Use(middleware.CORS())
	app.Use(middleware.Logger())

	// Swagger UI
	app.Get("/swagger/*", swagger.Handler())
	app.Get("/swagger", swagger.SwaggerUI("Glamping API", "/swagger/doc.json"))

	api := app.Group("/api/v1")

	api.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// Feature modules
	deps.AuthModule.RegisterRoutes(api)
	deps.UserModule.RegisterRoutes(api, deps.JWTSecret)
	deps.AmenityModule.RegisterRoutes(api, deps.JWTSecret)
	deps.TentTypeModule.RegisterRoutes(api, deps.JWTSecret)
	deps.TentModule.RegisterRoutes(api, deps.JWTSecret)
	deps.BookingModule.RegisterRoutes(api, deps.JWTSecret)
	deps.PaymentModule.RegisterRoutes(api, deps.JWTSecret)
}
