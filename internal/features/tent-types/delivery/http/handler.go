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

func (h *Handler) List(c fiber.Ctx) error {
	result, err := h.service.List()
	if err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, "tent types retrieved", result)
}

func (h *Handler) FindByID(c fiber.Ctx) error {
	id := c.Params("id")

	result, err := h.service.FindByID(id)
	if err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, "tent type retrieved", result)
}

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

func (h *Handler) Delete(c fiber.Ctx) error {
	id := c.Params("id")

	if err := h.service.Delete(id); err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, "tent type deleted", nil)
}

// --- Images ---

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

func (h *Handler) DeleteImage(c fiber.Ctx) error {
	imageID := c.Params("imageId")

	if err := h.service.DeleteImage(imageID); err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, "image deleted", nil)
}

func (h *Handler) SetPrimaryImage(c fiber.Ctx) error {
	tentTypeID := c.Params("id")
	imageID := c.Params("imageId")

	if err := h.service.SetPrimaryImage(tentTypeID, imageID); err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, "primary image set", nil)
}

// --- Rates ---

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

func (h *Handler) ListRates(c fiber.Ctx) error {
	tentTypeID := c.Params("id")

	result, err := h.service.ListRates(tentTypeID)
	if err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, "rates retrieved", result)
}

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

func (h *Handler) DeleteRate(c fiber.Ctx) error {
	rateID := c.Params("rateId")

	if err := h.service.DeleteRate(rateID); err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, "rate deleted", nil)
}
