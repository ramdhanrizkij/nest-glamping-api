package users

import (
	"github.com/gofiber/fiber/v3"
	usersHttp "github.com/ramdhanrizkij/nest-glamping-api/internal/features/users/delivery/http"
	"github.com/ramdhanrizkij/nest-glamping-api/internal/features/users/domain"
)

type Module struct {
	service domain.Service
}

func NewModule(service domain.Service) *Module {
	return &Module{service: service}
}

func (m *Module) RegisterRoutes(router fiber.Router, jwtSecret string) {
	handler := usersHttp.NewHandler(m.service)
	usersHttp.RegisterRoutes(router, handler, jwtSecret)
}
