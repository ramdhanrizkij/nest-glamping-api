package http

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/ramdhanrizkij/nest-glamping-api/internal/features/tents/domain"
	"github.com/ramdhanrizkij/nest-glamping-api/internal/features/tents/dto"
	"github.com/ramdhanrizkij/nest-glamping-api/internal/shared/response"
	"github.com/ramdhanrizkij/nest-glamping-api/internal/shared/validator"
)

type Handler struct {
	service domain.Service
}

func NewHandler(service domain.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Create(c fiber.Ctx) error {
	var req dto.CreateTentRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.Error(c, fiber.NewError(fiber.StatusBadRequest, "invalid request body"))
	}

	if errs := validator.Validate(req); errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"errors":  errs,
		})
	}

	result, err := h.service.Create(req)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Created(c, "tent created", result)
}

func (h *Handler) List(c fiber.Ctx) error {
	result, err := h.service.List()
	if err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, "tents retrieved", result)
}

func (h *Handler) ListByTentTypeID(c fiber.Ctx) error {
	tentTypeID := c.Params("id")

	result, err := h.service.ListByTentTypeID(tentTypeID)
	if err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, "tents retrieved", result)
}

func (h *Handler) FindByID(c fiber.Ctx) error {
	id := c.Params("id")

	result, err := h.service.FindByID(id)
	if err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, "tent retrieved", result)
}

func (h *Handler) Update(c fiber.Ctx) error {
	id := c.Params("id")

	var req dto.UpdateTentRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.Error(c, fiber.NewError(fiber.StatusBadRequest, "invalid request body"))
	}

	if errs := validator.Validate(req); errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"errors":  errs,
		})
	}

	result, err := h.service.Update(id, req)
	if err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, "tent updated", result)
}

func (h *Handler) Delete(c fiber.Ctx) error {
	id := c.Params("id")

	if err := h.service.Delete(id); err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, "tent deleted", nil)
}

func (h *Handler) CheckAvailability(c fiber.Ctx) error {
	tentTypeID := c.Params("id")
	checkInStr := c.Query("check_in")
	checkOutStr := c.Query("check_out")

	if checkInStr == "" || checkOutStr == "" {
		return response.Error(c, fiber.NewError(fiber.StatusBadRequest, "check_in and check_out are required"))
	}

	checkIn, err := time.Parse("2006-01-02", checkInStr)
	if err != nil {
		return response.Error(c, fiber.NewError(fiber.StatusBadRequest, "invalid check_in format, use YYYY-MM-DD"))
	}

	checkOut, err := time.Parse("2006-01-02", checkOutStr)
	if err != nil {
		return response.Error(c, fiber.NewError(fiber.StatusBadRequest, "invalid check_out format, use YYYY-MM-DD"))
	}

	if !checkOut.After(checkIn) {
		return response.Error(c, fiber.NewError(fiber.StatusBadRequest, "check_out must be after check_in"))
	}

	result, err := h.service.CheckAvailability(tentTypeID, checkIn, checkOut)
	if err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, "availability checked", result)
}
