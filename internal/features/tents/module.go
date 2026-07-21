package tents

import (
	"github.com/gofiber/fiber/v3"
	tentsHttp "github.com/ramdhanrizkij/nest-glamping-api/internal/features/tents/delivery/http"
	"github.com/ramdhanrizkij/nest-glamping-api/internal/features/tents/domain"
)

type Module struct {
	service domain.Service
}

func NewModule(service domain.Service) *Module {
	return &Module{service: service}
}

func (m *Module) RegisterRoutes(router fiber.Router, jwtSecret string) {
	handler := tentsHttp.NewHandler(m.service)
	tentsHttp.RegisterRoutes(router, handler, jwtSecret)
}
