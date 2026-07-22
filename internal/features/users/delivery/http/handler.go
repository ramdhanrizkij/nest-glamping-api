package http

import (
	"github.com/gofiber/fiber/v3"
	"github.com/ramdhanrizkij/nest-glamping-api/internal/features/users/domain"
	"github.com/ramdhanrizkij/nest-glamping-api/internal/features/users/dto"
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

// GetProfile godoc
// @Summary      Get user profile
// @Description  Get the authenticated user's profile
// @Tags         Users
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  response.Response{data=dto.UserResponse}
// @Failure      401  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /users/profile [get]
func (h *Handler) GetProfile(c fiber.Ctx) error {
	userID := auth.GetUserID(c)

	result, err := h.service.GetProfile(userID)
	if err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, "profile retrieved", result)
}

// UpdateProfile godoc
// @Summary      Update user profile
// @Description  Update the authenticated user's name or phone number
// @Tags         Users
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body  body      dto.UpdateProfileRequest  true  "Fields to update"
// @Success      200   {object}  response.Response{data=dto.UserResponse}
// @Failure      400   {object}  response.Response
// @Failure      401   {object}  response.Response
// @Router       /users/profile [put]
func (h *Handler) UpdateProfile(c fiber.Ctx) error {
	userID := auth.GetUserID(c)

	var req dto.UpdateProfileRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.Error(c, fiber.NewError(fiber.StatusBadRequest, "invalid request body"))
	}

	if errs := validator.Validate(req); errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"errors":  errs,
		})
	}

	result, err := h.service.UpdateProfile(userID, req)
	if err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, "profile updated", result)
}

// ListAll godoc
// @Summary      List all users
// @Description  Get a list of all users (admin only)
// @Tags         Admin - Users
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  response.Response{data=[]dto.UserResponse}
// @Failure      401  {object}  response.Response
// @Failure      403  {object}  response.Response
// @Router       /admin/users [get]
func (h *Handler) ListAll(c fiber.Ctx) error {
	result, err := h.service.ListAll()
	if err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, "users retrieved", result)
}

// GetUserByID godoc
// @Summary      Get user by ID
// @Description  Get a user by their ID (admin only)
// @Tags         Admin - Users
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  response.Response{data=dto.UserResponse}
// @Failure      401  {object}  response.Response
// @Failure      403  {object}  response.Response
// @Failure      404  {object}  response.Response
// @Router       /admin/users/{id} [get]
func (h *Handler) GetUserByID(c fiber.Ctx) error {
	id := c.Params("id")

	result, err := h.service.GetUserByID(id)
	if err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, "user retrieved", result)
}

// UpdateUser godoc
// @Summary      Update user
// @Description  Update a user's profile (admin only)
// @Tags         Admin - Users
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id    path      string                  true  "User ID"
// @Param        body  body      dto.UpdateUserRequest    true  "Fields to update"
// @Success      200   {object}  response.Response{data=dto.UserResponse}
// @Failure      400   {object}  response.Response
// @Failure      401   {object}  response.Response
// @Failure      403   {object}  response.Response
// @Failure      404   {object}  response.Response
// @Router       /admin/users/{id} [put]
func (h *Handler) UpdateUser(c fiber.Ctx) error {
	id := c.Params("id")

	var req dto.UpdateUserRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.Error(c, fiber.NewError(fiber.StatusBadRequest, "invalid request body"))
	}

	if errs := validator.Validate(req); errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"errors":  errs,
		})
	}

	result, err := h.service.UpdateUser(id, req)
	if err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, "user updated", result)
}

// DeleteUser godoc
// @Summary      Delete user
// @Description  Soft-delete a user (admin only)
// @Tags         Admin - Users
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Failure      403  {object}  response.Response
// @Failure      404  {object}  response.Response
// @Router       /admin/users/{id} [delete]
func (h *Handler) DeleteUser(c fiber.Ctx) error {
	id := c.Params("id")

	if err := h.service.DeleteUser(id); err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, "user deleted", nil)
}
