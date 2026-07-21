package response

import (
	"github.com/gofiber/fiber/v3"
	"github.com/ramdhanrizkij/nest-glamping-api/internal/shared/errors"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type PaginatedResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

type Meta struct {
	Page       int   `json:"page"`
	PerPage    int   `json:"per_page"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

func Success(c fiber.Ctx, code int, message string, data interface{}) error {
	return c.Status(code).JSON(Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func Created(c fiber.Ctx, message string, data interface{}) error {
	return Success(c, fiber.StatusCreated, message, data)
}

func OK(c fiber.Ctx, message string, data interface{}) error {
	return Success(c, fiber.StatusOK, message, data)
}

func Error(c fiber.Ctx, err error) error {
	if appErr, ok := err.(*errors.AppError); ok {
		return c.Status(appErr.Code).JSON(Response{
			Success: false,
			Error:   appErr.Message,
		})
	}
	return c.Status(fiber.StatusInternalServerError).JSON(Response{
		Success: false,
		Error:   "internal server error",
	})
}

func Paginated(c fiber.Ctx, message string, data interface{}, meta *Meta) error {
	return c.Status(fiber.StatusOK).JSON(PaginatedResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}
