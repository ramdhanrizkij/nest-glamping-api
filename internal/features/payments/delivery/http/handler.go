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

// InitiatePayment godoc
// @Summary      Initiate payment
// @Description  Initiate a payment for a booking
// @Tags         Payments
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id    path      string           true  "Booking ID"
// @Param        body  body      dto.PayRequest    true  "Payment method"
// @Success      200   {object}  response.Response{data=dto.PaymentResponse}
// @Failure      400   {object}  response.Response
// @Failure      401   {object}  response.Response
// @Failure      404   {object}  response.Response
// @Router       /bookings/{id}/pay [post]
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

// GetPaymentStatus godoc
// @Summary      Get payment status
// @Description  Get the payment status for a booking
// @Tags         Payments
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      string  true  "Booking ID"
// @Success      200  {object}  response.Response{data=dto.PaymentResponse}
// @Failure      401  {object}  response.Response
// @Failure      404  {object}  response.Response
// @Router       /bookings/{id}/payment [get]
func (h *Handler) GetPaymentStatus(c fiber.Ctx) error {
	userID := auth.GetUserID(c)
	bookingID := c.Params("id")

	result, err := h.service.GetPaymentByBookingID(userID, bookingID)
	if err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, "payment retrieved", result)
}

// HandleCallback godoc
// @Summary      Handle payment callback
// @Description  Process a payment gateway webhook callback
// @Tags         Payments
// @Accept       json
// @Produce      json
// @Param        id    path      string                true  "Payment ID"
// @Param        body  body      dto.CallbackRequest    true  "Callback payload"
// @Success      200   {object}  response.Response
// @Failure      400   {object}  response.Response
// @Router       /payments/{id}/callback [post]
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
