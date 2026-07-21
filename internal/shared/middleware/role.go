package middleware

import (
	"github.com/gofiber/fiber/v3"
	"github.com/ramdhanrizkij/nest-glamping-api/internal/shared/response"
)

func RoleAllowed(roleIDs ...string) fiber.Handler {
	return func(c fiber.Ctx) error {
		userRoleID, ok := c.Locals("roleID").(string)
		if !ok {
			return response.Error(c, fiber.NewError(fiber.StatusForbidden, "role not found"))
		}

		for _, allowed := range roleIDs {
			if userRoleID == allowed {
				return c.Next()
			}
		}

		return response.Error(c, fiber.NewError(fiber.StatusForbidden, "insufficient permissions"))
	}
}
