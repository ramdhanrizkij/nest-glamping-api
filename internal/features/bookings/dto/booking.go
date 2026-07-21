package dto

type CreateBookingRequest struct {
	TentTypeID     string   `json:"tent_type_id" validate:"required"`
	CheckInDate    string   `json:"check_in_date" validate:"required"`
	CheckOutDate   string   `json:"check_out_date" validate:"required"`
	GuestCount     int      `json:"guest_count" validate:"required,min=1"`
	SpecialRequests string  `json:"special_requests" validate:"omitempty"`
	TentIDs        []string `json:"tent_ids" validate:"required,min=1"`
}

type BookingResponse struct {
	ID              string              `json:"id"`
	BookingCode     string              `json:"booking_code"`
	CheckInDate     string              `json:"check_in_date"`
	CheckOutDate    string              `json:"check_out_date"`
	TotalAmount     float64             `json:"total_amount"`
	Status          string              `json:"status"`
	GuestCount      int                 `json:"guest_count"`
	SpecialRequests string              `json:"special_requests,omitempty"`
	Tents           []BookingTentItem   `json:"tents,omitempty"`
}

type BookingTentItem struct {
	ID            string  `json:"id"`
	TentID        string  `json:"tent_id"`
	TentName      string  `json:"tent_name,omitempty"`
	PricePerNight float64 `json:"price_per_night"`
}

type BookingDetailResponse struct {
	BookingResponse
	CreatedAt string `json:"created_at"`
}
