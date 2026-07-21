package http

import (
	"github.com/gofiber/fiber/v3"
	"github.com/ramdhanrizkij/nest-glamping-api/internal/features/auth/domain"
	"github.com/ramdhanrizkij/nest-glamping-api/internal/features/auth/dto"
	"github.com/ramdhanrizkij/nest-glamping-api/internal/shared/response"
	"github.com/ramdhanrizkij/nest-glamping-api/internal/shared/validator"
)

type Handler struct {
	service domain.Service
}

func NewHandler(service domain.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(c fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.Error(c, fiber.NewError(fiber.StatusBadRequest, "invalid request body"))
	}

	if errs := validator.Validate(req); errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"errors":  errs,
		})
	}

	result, err := h.service.Register(req)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Created(c, "registration successful", result)
}

func (h *Handler) Login(c fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.Error(c, fiber.NewError(fiber.StatusBadRequest, "invalid request body"))
	}

	if errs := validator.Validate(req); errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"errors":  errs,
		})
	}

	result, err := h.service.Login(req)
	if err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, "login successful", result)
}

func (h *Handler) RefreshToken(c fiber.Ctx) error {
	var req struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}
	if err := c.Bind().JSON(&req); err != nil {
		return response.Error(c, fiber.NewError(fiber.StatusBadRequest, "invalid request body"))
	}

	if errs := validator.Validate(req); errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"errors":  errs,
		})
	}

	result, err := h.service.RefreshToken(req.RefreshToken)
	if err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, "token refreshed", result)
}

func (h *Handler) Logout(c fiber.Ctx) error {
	var req struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}
	if err := c.Bind().JSON(&req); err != nil {
		return response.Error(c, fiber.NewError(fiber.StatusBadRequest, "invalid request body"))
	}

	if err := h.service.Logout(req.RefreshToken); err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, "logged out successfully", nil)
}
