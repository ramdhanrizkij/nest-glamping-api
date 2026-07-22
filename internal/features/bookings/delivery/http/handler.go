package http

import (
	"strconv"

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

// CreateBooking godoc
// @Summary      Create booking
// @Description  Create a new tent booking
// @Tags         Bookings
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body  body      dto.CreateBookingRequest  true  "Booking payload"
// @Success      201   {object}  response.Response{data=dto.BookingResponse}
// @Failure      400   {object}  response.Response
// @Failure      401   {object}  response.Response
// @Router       /bookings [post]
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

// ListMyBookings godoc
// @Summary      List my bookings
// @Description  Get a paginated list of the authenticated user's bookings
// @Tags         Bookings
// @Security     BearerAuth
// @Produce      json
// @Param        page     query    int  false  "Page number"       default(1)
// @Param        per_page query    int  false  "Items per page"    default(10)
// @Success      200      {object} response.Response{data=dto.BookingListResponse}
// @Failure      401      {object} response.Response
// @Router       /bookings [get]
func (h *Handler) ListMyBookings(c fiber.Ctx) error {
	userID := auth.GetUserID(c)
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "10"))

	result, err := h.service.ListMyBookings(userID, page, perPage)
	if err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, "bookings retrieved", result)
}

// GetBookingDetail godoc
// @Summary      Get booking detail
// @Description  Get the detail of the authenticated user's booking
// @Tags         Bookings
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      string  true  "Booking ID"
// @Success      200  {object}  response.Response{data=dto.BookingDetailResponse}
// @Failure      401  {object}  response.Response
// @Failure      404  {object}  response.Response
// @Router       /bookings/{id} [get]
func (h *Handler) GetBookingDetail(c fiber.Ctx) error {
	userID := auth.GetUserID(c)
	bookingID := c.Params("id")

	result, err := h.service.GetBookingDetail(bookingID, userID, false)
	if err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, "booking retrieved", result)
}

// CancelBooking godoc
// @Summary      Cancel booking
// @Description  Cancel a pending booking (customer only)
// @Tags         Bookings
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      string  true  "Booking ID"
// @Success      200  {object}  response.Response
// @Failure      400  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Failure      404  {object}  response.Response
// @Router       /bookings/{id}/cancel [patch]
func (h *Handler) CancelBooking(c fiber.Ctx) error {
	userID := auth.GetUserID(c)
	bookingID := c.Params("id")

	if err := h.service.CancelBooking(bookingID, userID); err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, "booking cancelled", nil)
}

// ListAllBookings godoc
// @Summary      List all bookings
// @Description  Get a paginated list of all bookings (admin only)
// @Tags         Admin - Bookings
// @Security     BearerAuth
// @Produce      json
// @Param        page     query    int  false  "Page number"       default(1)
// @Param        per_page query    int  false  "Items per page"    default(10)
// @Success      200      {object} response.Response{data=dto.BookingListResponse}
// @Failure      401      {object} response.Response
// @Failure      403      {object} response.Response
// @Router       /admin/bookings [get]
func (h *Handler) ListAllBookings(c fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "10"))

	result, err := h.service.ListAllBookings(page, perPage)
	if err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, "bookings retrieved", result)
}

// GetBookingDetailAdmin godoc
// @Summary      Get booking detail (admin)
// @Description  Get the detail of any booking (admin only)
// @Tags         Admin - Bookings
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      string  true  "Booking ID"
// @Success      200  {object}  response.Response{data=dto.BookingDetailResponse}
// @Failure      401  {object}  response.Response
// @Failure      403  {object}  response.Response
// @Failure      404  {object}  response.Response
// @Router       /admin/bookings/{id} [get]
func (h *Handler) GetBookingDetailAdmin(c fiber.Ctx) error {
	bookingID := c.Params("id")

	result, err := h.service.GetBookingDetail(bookingID, "", true)
	if err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, "booking retrieved", result)
}

// ConfirmBooking godoc
// @Summary      Confirm booking
// @Description  Confirm a pending booking (admin only)
// @Tags         Admin - Bookings
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      string  true  "Booking ID"
// @Success      200  {object}  response.Response
// @Failure      400  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Failure      403  {object}  response.Response
// @Failure      404  {object}  response.Response
// @Router       /admin/bookings/{id}/confirm [patch]
func (h *Handler) ConfirmBooking(c fiber.Ctx) error {
	bookingID := c.Params("id")

	if err := h.service.ConfirmBooking(bookingID); err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, "booking confirmed", nil)
}
