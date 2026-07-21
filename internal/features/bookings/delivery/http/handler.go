package http

import (
	"github.com/gofiber/fiber/v3"
	"github.com/ramdhanrizkij/nest-glamping-api/internal/features/bookings/domain"
	"github.com/ramdhanrizkij/nest-glamping-api/internal/features/bookings/dto"
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

func (h *Handler) CreateBooking(c fiber.Ctx) error {
	userID := auth.GetUserID(c)

	var req dto.CreateBookingRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.Error(c, fiber.NewError(fiber.StatusBadRequest, "invalid request body"))
	}

	if errs := validator.Validate(req); errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"errors":  errs,
		})
	}

	result, err := h.service.CreateBooking(userID, req)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Created(c, "booking created", result)
}

func (h *Handler) ListMyBookings(c fiber.Ctx) error {
	userID := auth.GetUserID(c)

	result, err := h.service.ListMyBookings(userID)
	if err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, "bookings retrieved", result)
}

func (h *Handler) GetBookingDetail(c fiber.Ctx) error {
	userID := auth.GetUserID(c)
	bookingID := c.Params("id")

	result, err := h.service.GetBookingDetail(bookingID, userID, false)
	if err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, "booking retrieved", result)
}

func (h *Handler) CancelBooking(c fiber.Ctx) error {
	userID := auth.GetUserID(c)
	bookingID := c.Params("id")

	if err := h.service.CancelBooking(bookingID, userID); err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, "booking cancelled", nil)
}

func (h *Handler) ListAllBookings(c fiber.Ctx) error {
	result, err := h.service.ListAllBookings()
	if err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, "bookings retrieved", result)
}

func (h *Handler) GetBookingDetailAdmin(c fiber.Ctx) error {
	bookingID := c.Params("id")

	result, err := h.service.GetBookingDetail(bookingID, "", true)
	if err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, "booking retrieved", result)
}

func (h *Handler) ConfirmBooking(c fiber.Ctx) error {
	bookingID := c.Params("id")

	if err := h.service.ConfirmBooking(bookingID); err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, "booking confirmed", nil)
}
