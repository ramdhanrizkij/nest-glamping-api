package payments

import (
	"github.com/gofiber/fiber/v3"
	paymentsHttp "github.com/ramdhanrizkij/nest-glamping-api/internal/features/payments/delivery/http"
	"github.com/ramdhanrizkij/nest-glamping-api/internal/features/payments/domain"
)

type Module struct {
	service domain.Service
}

func NewModule(service domain.Service) *Module {
	return &Module{service: service}
}

func (m *Module) RegisterRoutes(router fiber.Router, jwtSecret string) {
	handler := paymentsHttp.NewHandler(m.service)
	paymentsHttp.RegisterRoutes(router, handler, jwtSecret)
}
