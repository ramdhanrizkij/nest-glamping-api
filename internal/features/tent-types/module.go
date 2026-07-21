package tenttypes

import (
	"github.com/gofiber/fiber/v3"
	ttentTypesHttp "github.com/ramdhanrizkij/nest-glamping-api/internal/features/tent-types/delivery/http"
	"github.com/ramdhanrizkij/nest-glamping-api/internal/features/tent-types/domain"
)

type Module struct {
	service domain.Service
}

func NewModule(service domain.Service) *Module {
	return &Module{service: service}
}

func (m *Module) RegisterRoutes(router fiber.Router, jwtSecret string) {
	handler := ttentTypesHttp.NewHandler(m.service)
	ttentTypesHttp.RegisterRoutes(router, handler, jwtSecret)
}
