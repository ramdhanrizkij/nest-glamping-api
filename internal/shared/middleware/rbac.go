package middleware

import (
	"sync"

	"github.com/gofiber/fiber/v3"
	"github.com/ramdhanrizkij/nest-glamping-api/internal/shared/response"
	"gorm.io/gorm"
)

var (
	permissionCache   = make(map[string][]string)
	permissionCacheMu sync.RWMutex
)

func PermissionRequired(db *gorm.DB, permissionName string) fiber.Handler {
	return func(c fiber.Ctx) error {
		roleID, ok := c.Locals("roleID").(string)
		if !ok {
			return response.Error(c, fiber.NewError(fiber.StatusForbidden, "role not found"))
		}

		permissions := getPermissions(db, roleID)

		for _, p := range permissions {
			if p == permissionName {
				return c.Next()
			}
		}

		return response.Error(c, fiber.NewError(fiber.StatusForbidden, "insufficient permissions"))
	}
}

func getPermissions(db *gorm.DB, roleID string) []string {
	permissionCacheMu.RLock()
	if perms, ok := permissionCache[roleID]; ok {
		permissionCacheMu.RUnlock()
		return perms
	}
	permissionCacheMu.RUnlock()

	var perms []string
	err := db.Raw(`
		SELECT p.name FROM permissions p
		JOIN role_permissions rp ON rp.permission_id = p.id
		WHERE rp.role_id = ?
	`, roleID).Scan(&perms).Error
	if err != nil {
		return nil
	}

	permissionCacheMu.Lock()
	permissionCache[roleID] = perms
	permissionCacheMu.Unlock()

	return perms
}

func ClearPermissionCache() {
	permissionCacheMu.Lock()
	defer permissionCacheMu.Unlock()
	permissionCache = make(map[string][]string)
}

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
