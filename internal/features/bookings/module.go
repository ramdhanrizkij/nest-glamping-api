package bookings

import (
	"github.com/gofiber/fiber/v3"
	bookingsHttp "github.com/ramdhanrizkij/nest-glamping-api/internal/features/bookings/delivery/http"
	"github.com/ramdhanrizkij/nest-glamping-api/internal/features/bookings/domain"
)

type Module struct {
	service domain.Service
}

func NewModule(service domain.Service) *Module {
	return &Module{service: service}
}

func (m *Module) RegisterRoutes(router fiber.Router, jwtSecret string) {
	handler := bookingsHttp.NewHandler(m.service)
	bookingsHttp.RegisterRoutes(router, handler, jwtSecret)
}
