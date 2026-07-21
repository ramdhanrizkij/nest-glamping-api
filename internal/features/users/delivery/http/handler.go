package http

import (
	"github.com/gofiber/fiber/v3"
	"github.com/ramdhanrizkij/nest-glamping-api/internal/features/users/domain"
	"github.com/ramdhanrizkij/nest-glamping-api/internal/features/users/dto"
	"github.com/ramdhanrizkij/nest-glamping-api/internal/shared/auth"
	"github.com/ramdhanrizkij/nest-glamping-api/internal/shared/response"
	"github.com/ramdhanrizkij/nest-glamping-api/internal/shared/validator"
)

type Handler struct {
	service domain.Service
}

func NewHandler(service domain.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetProfile(c fiber.Ctx) error {
	userID := auth.GetUserID(c)

	result, err := h.service.GetProfile(userID)
	if err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, "profile retrieved", result)
}

func (h *Handler) UpdateProfile(c fiber.Ctx) error {
	userID := auth.GetUserID(c)

	var req dto.UpdateProfileRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.Error(c, fiber.NewError(fiber.StatusBadRequest, "invalid request body"))
	}

	if errs := validator.Validate(req); errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"errors":  errs,
		})
	}

	result, err := h.service.UpdateProfile(userID, req)
	if err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, "profile updated", result)
}
