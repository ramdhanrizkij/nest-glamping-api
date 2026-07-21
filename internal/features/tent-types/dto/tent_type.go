package dto

type CreateTentTypeRequest struct {
	Name        string  `json:"name" validate:"required,max=255"`
	Description string  `json:"description" validate:"omitempty"`
	Capacity    int     `json:"capacity" validate:"required,min=1"`
	BasePrice   float64 `json:"base_price" validate:"required,gt=0"`
	AmenityIDs  []string `json:"amenity_ids" validate:"omitempty"`
}

type UpdateTentTypeRequest struct {
	Name        string  `json:"name" validate:"omitempty,max=255"`
	Description string  `json:"description" validate:"omitempty"`
	Capacity    int     `json:"capacity" validate:"omitempty,min=1"`
	BasePrice   float64 `json:"base_price" validate:"omitempty,gt=0"`
	AmenityIDs  []string `json:"amenity_ids" validate:"omitempty"`
}

type TentTypeResponse struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description,omitempty"`
	Capacity    int     `json:"capacity"`
	BasePrice   float64 `json:"base_price"`
}

type TentTypeDetailResponse struct {
	TentTypeResponse
	Images  []TentTypeImageResponse  `json:"images"`
	Amenities []TentTypeAmenityResponse `json:"amenities"`
	Rates    []TentTypeRateResponse    `json:"rates"`
}

type TentTypeImageResponse struct {
	ID        string `json:"id"`
	ImageURL  string `json:"image_url"`
	IsPrimary bool   `json:"is_primary"`
}

type TentTypeRateResponse struct {
	ID            string  `json:"id"`
	StartDate     string  `json:"start_date"`
	EndDate       string  `json:"end_date"`
	PricePerNight float64 `json:"price_per_night"`
	Description   string  `json:"description,omitempty"`
	IsActive      bool    `json:"is_active"`
}

type TentTypeAmenityResponse struct {
	AmenityID string `json:"amenity_id"`
	Name      string `json:"name"`
}

type AddImageRequest struct {
	ImageURL  string `json:"image_url" validate:"required,max=255"`
	IsPrimary bool   `json:"is_primary"`
}
