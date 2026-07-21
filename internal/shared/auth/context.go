package auth

import "github.com/gofiber/fiber/v3"

func GetUserID(c fiber.Ctx) string {
	return c.Locals("userID").(string)
}

func GetRoleID(c fiber.Ctx) string {
	return c.Locals("roleID").(string)
}
