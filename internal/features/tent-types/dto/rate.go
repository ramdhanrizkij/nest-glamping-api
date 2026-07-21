package dto

type CreateRateRequest struct {
	StartDate     string  `json:"start_date" validate:"required"`
	EndDate       string  `json:"end_date" validate:"required"`
	PricePerNight float64 `json:"price_per_night" validate:"required,gt=0"`
	Description   string  `json:"description" validate:"omitempty,max=255"`
	IsActive      *bool   `json:"is_active"`
}

type UpdateRateRequest struct {
	StartDate     string  `json:"start_date" validate:"omitempty"`
	EndDate       string  `json:"end_date" validate:"omitempty"`
	PricePerNight float64 `json:"price_per_night" validate:"omitempty,gt=0"`
	Description   string  `json:"description" validate:"omitempty,max=255"`
	IsActive      *bool   `json:"is_active"`
}
