package http

import (
	"github.com/gofiber/fiber/v3"
	"github.com/ramdhanrizkij/nest-glamping-api/internal/features/payments/domain"
	"github.com/ramdhanrizkij/nest-glamping-api/internal/features/payments/dto"
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

func (h *Handler) InitiatePayment(c fiber.Ctx) error {
	userID := auth.GetUserID(c)
	bookingID := c.Params("id")

	var req dto.PayRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.Error(c, fiber.NewError(fiber.StatusBadRequest, "invalid request body"))
	}

	if errs := validator.Validate(req); errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"errors":  errs,
		})
	}

	result, err := h.service.InitiatePayment(userID, bookingID, req)
	if err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, "payment initiated", result)
}

func (h *Handler) GetPaymentStatus(c fiber.Ctx) error {
	userID := auth.GetUserID(c)
	bookingID := c.Params("id")

	result, err := h.service.GetPaymentByBookingID(userID, bookingID)
	if err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, "payment retrieved", result)
}

func (h *Handler) HandleCallback(c fiber.Ctx) error {
	var req dto.CallbackRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.Error(c, fiber.NewError(fiber.StatusBadRequest, "invalid request body"))
	}

	if err := h.service.HandleCallback(req); err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, "callback processed", nil)
}
