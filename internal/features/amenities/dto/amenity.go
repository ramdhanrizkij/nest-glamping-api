package dto

type CreateAmenityRequest struct {
	Name        string `json:"name" validate:"required,max=100"`
	IconURL     string `json:"icon_url" validate:"omitempty,max=255"`
	Description string `json:"description" validate:"omitempty"`
}

type UpdateAmenityRequest struct {
	Name        string `json:"name" validate:"omitempty,max=100"`
	IconURL     string `json:"icon_url" validate:"omitempty,max=255"`
	Description string `json:"description" validate:"omitempty"`
}

type AmenityResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	IconURL     string `json:"icon_url,omitempty"`
	Description string `json:"description,omitempty"`
}
