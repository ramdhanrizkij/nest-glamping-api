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

// Create godoc
// @Summary      Create tent unit
// @Description  Create a new tent unit (admin only)
// @Tags         Tents
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body  body      dto.CreateTentRequest  true  "Tent unit payload"
// @Success      201   {object}  response.Response{data=dto.TentResponse}
// @Failure      400   {object}  response.Response
// @Failure      401   {object}  response.Response
// @Failure      403   {object}  response.Response
// @Router       /tents/units [post]
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

// List godoc
// @Summary      List tent units
// @Description  Get a list of all tent units (admin only)
// @Tags         Tents
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  response.Response{data=[]dto.TentResponse}
// @Failure      401  {object}  response.Response
// @Failure      403  {object}  response.Response
// @Router       /tents/units [get]
func (h *Handler) List(c fiber.Ctx) error {
	result, err := h.service.List()
	if err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, "tents retrieved", result)
}

// ListByTentTypeID godoc
// @Summary      List tent units by type
// @Description  Get tent units filtered by tent type ID (admin only)
// @Tags         Tents
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      string  true  "Tent Type ID"
// @Success      200  {object}  response.Response{data=[]dto.TentResponse}
// @Failure      401  {object}  response.Response
// @Failure      403  {object}  response.Response
// @Router       /tents/tent-types/{id}/units [get]
func (h *Handler) ListByTentTypeID(c fiber.Ctx) error {
	tentTypeID := c.Params("id")

	result, err := h.service.ListByTentTypeID(tentTypeID)
	if err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, "tents retrieved", result)
}

// FindByID godoc
// @Summary      Get tent unit by ID
// @Description  Get a tent unit by its ID (admin only)
// @Tags         Tents
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      string  true  "Tent Unit ID"
// @Success      200  {object}  response.Response{data=dto.TentResponse}
// @Failure      401  {object}  response.Response
// @Failure      403  {object}  response.Response
// @Failure      404  {object}  response.Response
// @Router       /tents/units/{id} [get]
func (h *Handler) FindByID(c fiber.Ctx) error {
	id := c.Params("id")

	result, err := h.service.FindByID(id)
	if err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, "tent retrieved", result)
}

// Update godoc
// @Summary      Update tent unit
// @Description  Update a tent unit (admin only)
// @Tags         Tents
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id    path      string                true  "Tent Unit ID"
// @Param        body  body      dto.UpdateTentRequest  true  "Fields to update"
// @Success      200   {object}  response.Response{data=dto.TentResponse}
// @Failure      400   {object}  response.Response
// @Failure      401   {object}  response.Response
// @Failure      403   {object}  response.Response
// @Failure      404   {object}  response.Response
// @Router       /tents/units/{id} [put]
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

// Delete godoc
// @Summary      Delete tent unit
// @Description  Delete a tent unit (admin only)
// @Tags         Tents
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      string  true  "Tent Unit ID"
// @Success      200  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Failure      403  {object}  response.Response
// @Failure      404  {object}  response.Response
// @Router       /tents/units/{id} [delete]
func (h *Handler) Delete(c fiber.Ctx) error {
	id := c.Params("id")

	if err := h.service.Delete(id); err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, "tent deleted", nil)
}

// CheckAvailability godoc
// @Summary      Check tent availability
// @Description  Check available tent units for a given date range
// @Tags         Tents
// @Produce      json
// @Param        id        path     string  true   "Tent Type ID"
// @Param        check_in  query    string  true   "Check-in date (YYYY-MM-DD)"
// @Param        check_out query    string  true   "Check-out date (YYYY-MM-DD)"
// @Success      200       {object} response.Response{data=[]dto.AvailableTentResponse}
// @Failure      400       {object} response.Response
// @Failure      404       {object} response.Response
// @Router       /tents/tent-types/{id}/availability [get]
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
