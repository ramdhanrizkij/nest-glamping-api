package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/ramdhanrizkij/nest-glamping-api/internal/shared/response"
	"github.com/ramdhanrizkij/nest-glamping-api/pkg/jwt"
)

func Auth(jwtSecret string) fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return response.Error(c, fiber.NewError(fiber.StatusUnauthorized, "missing authorization header"))
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			return response.Error(c, fiber.NewError(fiber.StatusUnauthorized, "invalid authorization format"))
		}

		claims, err := jwt.ValidateToken(parts[1], jwtSecret)
		if err != nil {
			return response.Error(c, fiber.NewError(fiber.StatusUnauthorized, "invalid or expired token"))
		}

		c.Locals("userID", claims.UserID)
		c.Locals("roleID", claims.RoleID)
		return c.Next()
	}
}
