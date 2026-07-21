package amenities

import (
	"github.com/gofiber/fiber/v3"
	amenitiesHttp "github.com/ramdhanrizkij/nest-glamping-api/internal/features/amenities/delivery/http"
	"github.com/ramdhanrizkij/nest-glamping-api/internal/features/amenities/domain"
)

type Module struct {
	service domain.Service
}

func NewModule(service domain.Service) *Module {
	return &Module{service: service}
}

func (m *Module) RegisterRoutes(router fiber.Router, jwtSecret string) {
	handler := amenitiesHttp.NewHandler(m.service)
	amenitiesHttp.RegisterRoutes(router, handler, jwtSecret)
}
