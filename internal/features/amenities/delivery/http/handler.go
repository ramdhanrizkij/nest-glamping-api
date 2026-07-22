package http

import (
	"github.com/gofiber/fiber/v3"
	"github.com/ramdhanrizkij/nest-glamping-api/internal/features/amenities/domain"
	"github.com/ramdhanrizkij/nest-glamping-api/internal/features/amenities/dto"
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
// @Summary      Create amenity
// @Description  Create a new amenity (admin only)
// @Tags         Amenities
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body  body      dto.CreateAmenityRequest  true  "Amenity payload"
// @Success      201   {object}  response.Response{data=dto.AmenityResponse}
// @Failure      400   {object}  response.Response
// @Failure      401   {object}  response.Response
// @Failure      403   {object}  response.Response
// @Router       /amenities [post]
func (h *Handler) Create(c fiber.Ctx) error {
	var req dto.CreateAmenityRequest
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

	return response.Created(c, "amenity created", result)
}

// List godoc
// @Summary      List amenities
// @Description  Get a list of all amenities
// @Tags         Amenities
// @Produce      json
// @Success      200  {object}  response.Response{data=[]dto.AmenityResponse}
// @Failure      500  {object}  response.Response
// @Router       /amenities [get]
func (h *Handler) List(c fiber.Ctx) error {
	result, err := h.service.List()
	if err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, "amenities retrieved", result)
}

// FindByID godoc
// @Summary      Get amenity by ID
// @Description  Get an amenity by its ID
// @Tags         Amenities
// @Produce      json
// @Param        id   path      string  true  "Amenity ID"
// @Success      200  {object}  response.Response{data=dto.AmenityResponse}
// @Failure      404  {object}  response.Response
// @Router       /amenities/{id} [get]
func (h *Handler) FindByID(c fiber.Ctx) error {
	id := c.Params("id")

	result, err := h.service.FindByID(id)
	if err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, "amenity retrieved", result)
}

// Update godoc
// @Summary      Update amenity
// @Description  Update an amenity (admin only)
// @Tags         Amenities
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id    path      string                   true  "Amenity ID"
// @Param        body  body      dto.UpdateAmenityRequest  true  "Fields to update"
// @Success      200   {object}  response.Response{data=dto.AmenityResponse}
// @Failure      400   {object}  response.Response
// @Failure      401   {object}  response.Response
// @Failure      403   {object}  response.Response
// @Failure      404   {object}  response.Response
// @Router       /amenities/{id} [put]
func (h *Handler) Update(c fiber.Ctx) error {
	id := c.Params("id")

	var req dto.UpdateAmenityRequest
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

	return response.OK(c, "amenity updated", result)
}

// Delete godoc
// @Summary      Delete amenity
// @Description  Delete an amenity (admin only)
// @Tags         Amenities
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      string  true  "Amenity ID"
// @Success      200  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Failure      403  {object}  response.Response
// @Failure      404  {object}  response.Response
// @Router       /amenities/{id} [delete]
func (h *Handler) Delete(c fiber.Ctx) error {
	id := c.Params("id")

	if err := h.service.Delete(id); err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, "amenity deleted", nil)
}
