package dto

type PayRequest struct {
	PaymentMethod string `json:"payment_method" validate:"required"`
}

type PaymentResponse struct {
	ID            string  `json:"id"`
	BookingID     string  `json:"booking_id"`
	Amount        float64 `json:"amount"`
	PaymentMethod string  `json:"payment_method"`
	PaymentStatus string  `json:"payment_status"`
	PaymentDate   string  `json:"payment_date,omitempty"`
	GatewayRef    string  `json:"gateway_ref,omitempty"`
	PaymentURL    string  `json:"payment_url,omitempty"`
}

type CallbackRequest struct {
	ExternalID    string `json:"external_id" validate:"required"`
	Status        string `json:"status" validate:"required"`
	TransactionID string `json:"transaction_id"`
}
