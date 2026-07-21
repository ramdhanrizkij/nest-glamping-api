package auth

import (
	"github.com/gofiber/fiber/v3"
	authHttp "github.com/ramdhanrizkij/nest-glamping-api/internal/features/auth/delivery/http"
	"github.com/ramdhanrizkij/nest-glamping-api/internal/features/auth/domain"
)

type Module struct {
	service domain.Service
}

func NewModule(service domain.Service) *Module {
	return &Module{service: service}
}

func (m *Module) RegisterRoutes(router fiber.Router) {
	handler := authHttp.NewHandler(m.service)
	authHttp.RegisterRoutes(router, handler)
}
