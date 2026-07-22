package domain

import "github.com/ramdhanrizkij/nest-glamping-api/internal/features/payments/dto"

type Service interface {
	InitiatePayment(userID, bookingID string, req dto.PayRequest) (*dto.PaymentResponse, error)
	GetPaymentByBookingID(userID, bookingID string) (*dto.PaymentResponse, error)
	HandleCallback(req dto.CallbackRequest) error
}
