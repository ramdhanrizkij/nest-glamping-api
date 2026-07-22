package dto

type CreateTentRequest struct {
	TentTypeID string `json:"tent_type_id" validate:"required"`
	Code       string `json:"name_or_number" validate:"required,max=100"`
	Status     string `json:"status" validate:"omitempty,oneof=available maintenance"`
}

type UpdateTentRequest struct {
	Code   string `json:"name_or_number" validate:"omitempty,max=100"`
	Status string `json:"status" validate:"omitempty,oneof=available maintenance"`
}

type TentResponse struct {
	ID         string `json:"id"`
	TentTypeID string `json:"tent_type_id"`
	Code       string `json:"name_or_number"`
	Status     string `json:"status"`
}

type AvailableTentResponse struct {
	ID            string  `json:"id"`
	Code          string  `json:"name_or_number"`
	PricePerNight float64 `json:"price_per_night"`
}
