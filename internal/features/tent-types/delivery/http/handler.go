package http

import (
	"github.com/gofiber/fiber/v3"
	"github.com/ramdhanrizkij/nest-glamping-api/internal/features/tent-types/domain"
	"github.com/ramdhanrizkij/nest-glamping-api/internal/features/tent-types/dto"
	"github.com/ramdhanrizkij/nest-glamping-api/internal/shared/response"
	"github.com/ramdhanrizkij/nest-glamping-api/internal/shared/validator"
)

type Handler struct {
	service domain.Service
}

func NewHandler(service domain.Service) *Handler {
	return &Handler{service: service}
}

// --- TentType CRUD ---

// Create godoc
// @Summary      Create tent type
// @Description  Create a new tent type (admin only)
// @Tags         Tent Types
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body  body      dto.CreateTentTypeRequest  true  "Tent type payload"
// @Success      201   {object}  response.Response{data=dto.TentTypeResponse}
// @Failure      400   {object}  response.Response
// @Failure      401   {object}  response.Response
// @Failure      403   {object}  response.Response
// @Router       /tent-types [post]
func (h *Handler) Create(c fiber.Ctx) error {
	var req dto.CreateTentTypeRequest
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

	return response.Created(c, "tent type created", result)
}

// List godoc
// @Summary      List tent types
// @Description  Get a list of all tent types
// @Tags         Tent Types
// @Produce      json
// @Success      200  {object}  response.Response{data=[]dto.TentTypeResponse}
// @Failure      500  {object}  response.Response
// @Router       /tent-types [get]
func (h *Handler) List(c fiber.Ctx) error {
	result, err := h.service.List()
	if err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, "tent types retrieved", result)
}

// FindByID godoc
// @Summary      Get tent type detail
// @Description  Get a tent type with images, amenities, and rates
// @Tags         Tent Types
// @Produce      json
// @Param        id   path      string  true  "Tent Type ID"
// @Success      200  {object}  response.Response{data=dto.TentTypeDetailResponse}
// @Failure      404  {object}  response.Response
// @Router       /tent-types/{id} [get]
func (h *Handler) FindByID(c fiber.Ctx) error {
	id := c.Params("id")

	result, err := h.service.FindByID(id)
	if err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, "tent type retrieved", result)
}

// Update godoc
// @Summary      Update tent type
// @Description  Update a tent type (admin only)
// @Tags         Tent Types
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id    path      string                   true  "Tent Type ID"
// @Param        body  body      dto.UpdateTentTypeRequest  true  "Fields to update"
// @Success      200   {object}  response.Response{data=dto.TentTypeResponse}
// @Failure      400   {object}  response.Response
// @Failure      401   {object}  response.Response
// @Failure      403   {object}  response.Response
// @Failure      404   {object}  response.Response
// @Router       /tent-types/{id} [put]
func (h *Handler) Update(c fiber.Ctx) error {
	id := c.Params("id")

	var req dto.UpdateTentTypeRequest
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

	return response.OK(c, "tent type updated", result)
}

// Delete godoc
// @Summary      Delete tent type
// @Description  Delete a tent type (admin only)
// @Tags         Tent Types
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      string  true  "Tent Type ID"
// @Success      200  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Failure      403  {object}  response.Response
// @Failure      404  {object}  response.Response
// @Router       /tent-types/{id} [delete]
func (h *Handler) Delete(c fiber.Ctx) error {
	id := c.Params("id")

	if err := h.service.Delete(id); err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, "tent type deleted", nil)
}

// --- Images ---

// AddImage godoc
// @Summary      Add image to tent type
// @Description  Add an image to a tent type (admin only)
// @Tags         Tent Types - Images
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id    path      string              true  "Tent Type ID"
// @Param        body  body      dto.AddImageRequest  true  "Image payload"
// @Success      201   {object}  response.Response{data=dto.TentTypeImageResponse}
// @Failure      400   {object}  response.Response
// @Failure      401   {object}  response.Response
// @Failure      403   {object}  response.Response
// @Router       /tent-types/{id}/images [post]
func (h *Handler) AddImage(c fiber.Ctx) error {
	tentTypeID := c.Params("id")

	var req dto.AddImageRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.Error(c, fiber.NewError(fiber.StatusBadRequest, "invalid request body"))
	}

	if errs := validator.Validate(req); errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"errors":  errs,
		})
	}

	result, err := h.service.AddImage(tentTypeID, req)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Created(c, "image added", result)
}

// DeleteImage godoc
// @Summary      Delete tent type image
// @Description  Delete an image from a tent type (admin only)
// @Tags         Tent Types - Images
// @Security     BearerAuth
// @Produce      json
// @Param        id       path      string  true  "Tent Type ID"
// @Param        imageId  path      string  true  "Image ID"
// @Success      200      {object}  response.Response
// @Failure      401      {object}  response.Response
// @Failure      403      {object}  response.Response
// @Failure      404      {object}  response.Response
// @Router       /tent-types/{id}/images/{imageId} [delete]
func (h *Handler) DeleteImage(c fiber.Ctx) error {
	imageID := c.Params("imageId")

	if err := h.service.DeleteImage(imageID); err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, "image deleted", nil)
}

// SetPrimaryImage godoc
// @Summary      Set primary image
// @Description  Set an image as the primary image for a tent type (admin only)
// @Tags         Tent Types - Images
// @Security     BearerAuth
// @Produce      json
// @Param        id       path      string  true  "Tent Type ID"
// @Param        imageId  path      string  true  "Image ID"
// @Success      200      {object}  response.Response
// @Failure      401      {object}  response.Response
// @Failure      403      {object}  response.Response
// @Failure      404      {object}  response.Response
// @Router       /tent-types/{id}/images/{imageId}/primary [put]
func (h *Handler) SetPrimaryImage(c fiber.Ctx) error {
	tentTypeID := c.Params("id")
	imageID := c.Params("imageId")

	if err := h.service.SetPrimaryImage(tentTypeID, imageID); err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, "primary image set", nil)
}

// --- Rates ---

// CreateRate godoc
// @Summary      Create dynamic rate
// @Description  Create a dynamic rate for a tent type (admin only)
// @Tags         Tent Types - Rates
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id    path      string                  true  "Tent Type ID"
// @Param        body  body      dto.CreateRateRequest    true  "Rate payload"
// @Success      201   {object}  response.Response{data=dto.TentTypeRateResponse}
// @Failure      400   {object}  response.Response
// @Failure      401   {object}  response.Response
// @Failure      403   {object}  response.Response
// @Router       /tent-types/{id}/rates [post]
func (h *Handler) CreateRate(c fiber.Ctx) error {
	tentTypeID := c.Params("id")

	var req dto.CreateRateRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.Error(c, fiber.NewError(fiber.StatusBadRequest, "invalid request body"))
	}

	if errs := validator.Validate(req); errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"errors":  errs,
		})
	}

	result, err := h.service.CreateRate(tentTypeID, req)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Created(c, "rate created", result)
}

// ListRates godoc
// @Summary      List rates for tent type
// @Description  Get all dynamic rates for a tent type
// @Tags         Tent Types - Rates
// @Produce      json
// @Param        id   path      string  true  "Tent Type ID"
// @Success      200  {object}  response.Response{data=[]dto.TentTypeRateResponse}
// @Failure      404  {object}  response.Response
// @Router       /tent-types/{id}/rates [get]
func (h *Handler) ListRates(c fiber.Ctx) error {
	tentTypeID := c.Params("id")

	result, err := h.service.ListRates(tentTypeID)
	if err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, "rates retrieved", result)
}

// UpdateRate godoc
// @Summary      Update dynamic rate
// @Description  Update a dynamic rate (admin only)
// @Tags         Tent Types - Rates
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        rateId  path      string                 true  "Rate ID"
// @Param        body    body      dto.UpdateRateRequest   true  "Fields to update"
// @Success      200     {object}  response.Response{data=dto.TentTypeRateResponse}
// @Failure      400     {object}  response.Response
// @Failure      401     {object}  response.Response
// @Failure      403     {object}  response.Response
// @Failure      404     {object}  response.Response
// @Router       /tent-types/{id}/rates/{rateId} [put]
func (h *Handler) UpdateRate(c fiber.Ctx) error {
	rateID := c.Params("rateId")

	var req dto.UpdateRateRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.Error(c, fiber.NewError(fiber.StatusBadRequest, "invalid request body"))
	}

	if errs := validator.Validate(req); errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"errors":  errs,
		})
	}

	result, err := h.service.UpdateRate(rateID, req)
	if err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, "rate updated", result)
}

// DeleteRate godoc
// @Summary      Delete dynamic rate
// @Description  Delete a dynamic rate (admin only)
// @Tags         Tent Types - Rates
// @Security     BearerAuth
// @Produce      json
// @Param        rateId  path      string  true  "Rate ID"
// @Success      200     {object}  response.Response
// @Failure      401     {object}  response.Response
// @Failure      403     {object}  response.Response
// @Failure      404     {object}  response.Response
// @Router       /tent-types/{id}/rates/{rateId} [delete]
func (h *Handler) DeleteRate(c fiber.Ctx) error {
	rateID := c.Params("rateId")

	if err := h.service.DeleteRate(rateID); err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, "rate deleted", nil)
}
