package usecase

import (
	"github.com/google/uuid"
	bookingDomain "github.com/ramdhanrizkij/nest-glamping-api/internal/features/bookings/domain"
	"github.com/ramdhanrizkij/nest-glamping-api/internal/features/payments/domain"
	"github.com/ramdhanrizkij/nest-glamping-api/internal/features/payments/dto"
	pg "github.com/ramdhanrizkij/nest-glamping-api/pkg/payment_gateway"
	appErr "github.com/ramdhanrizkij/nest-glamping-api/internal/shared/errors"
)

type usecase struct {
	paymentRepo  domain.Repository
	bookingRepo  bookingDomain.Repository
	gateway      pg.PaymentGateway
}

func NewUsecase(paymentRepo domain.Repository, bookingRepo bookingDomain.Repository, gateway pg.PaymentGateway) domain.Service {
	return &usecase{
		paymentRepo: paymentRepo,
		bookingRepo: bookingRepo,
		gateway:     gateway,
	}
}

func (u *usecase) InitiatePayment(userID, bookingID string, req dto.PayRequest) (*dto.PaymentResponse, error) {
	bookingUUID, err := uuid.Parse(bookingID)
	if err != nil {
		return nil, appErr.BadRequest("invalid booking id")
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, appErr.BadRequest("invalid user id")
	}

	booking, _, err := u.bookingRepo.FindBookingByIDWithTents(bookingUUID)
	if err != nil {
		return nil, appErr.NotFound("booking not found")
	}

	if booking.UserID != userUUID {
		return nil, appErr.NotFound("booking not found")
	}

	if booking.Status != "pending" {
		return nil, appErr.BadRequest("can only pay for pending bookings")
	}

	payment, _ := u.paymentRepo.FindByBookingID(bookingUUID)
	if payment == nil {
		payment = &domain.Payment{
			ID:            uuid.New(),
			BookingID:     bookingUUID,
			Amount:        booking.TotalAmount,
			PaymentMethod: req.PaymentMethod,
			PaymentStatus: "unpaid",
		}
		if err := u.paymentRepo.Create(payment); err != nil {
			return nil, appErr.Internal("failed to create payment")
		}
	} else if payment.PaymentStatus == "paid" {
		return nil, appErr.BadRequest("booking already paid")
	} else {
		payment.PaymentMethod = req.PaymentMethod
	}

	result, err := u.gateway.CreatePayment(booking.BookingCode, booking.TotalAmount, req.PaymentMethod)
	if err != nil {
		return nil, appErr.Internal("failed to initiate payment")
	}

	return &dto.PaymentResponse{
		ID:            payment.ID.String(),
		BookingID:     booking.ID.String(),
		Amount:        payment.Amount,
		PaymentMethod: payment.PaymentMethod,
		PaymentStatus: payment.PaymentStatus,
		PaymentURL:    result.PaymentURL,
	}, nil
}

func (u *usecase) GetPaymentByBookingID(userID, bookingID string) (*dto.PaymentResponse, error) {
	bookingUUID, err := uuid.Parse(bookingID)
	if err != nil {
		return nil, appErr.BadRequest("invalid booking id")
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, appErr.BadRequest("invalid user id")
	}

	booking, _, err := u.bookingRepo.FindBookingByIDWithTents(bookingUUID)
	if err != nil {
		return nil, appErr.NotFound("booking not found")
	}

	if booking.UserID != userUUID {
		return nil, appErr.NotFound("booking not found")
	}

	payment, err := u.paymentRepo.FindByBookingID(bookingUUID)
	if err != nil {
		return nil, appErr.NotFound("payment not found")
	}

	resp := &dto.PaymentResponse{
		ID:            payment.ID.String(),
		BookingID:     payment.BookingID.String(),
		Amount:        payment.Amount,
		PaymentMethod: payment.PaymentMethod,
		PaymentStatus: payment.PaymentStatus,
		GatewayRef:    payment.GatewayRef,
	}

	if payment.PaymentDate != nil {
		resp.PaymentDate = payment.PaymentDate.Format("2006-01-02 15:04:05")
	}

	return resp, nil
}

func (u *usecase) HandleCallback(req dto.CallbackRequest) error {
	payment, err := u.paymentRepo.FindByID(uuid.MustParse(req.ExternalID))
	if err != nil {
		return appErr.NotFound("payment not found")
	}

	switch req.Status {
	case "paid":
		u.paymentRepo.UpdateStatus(payment.ID, "paid", req.TransactionID)
		u.bookingRepo.UpdateBookingStatus(payment.BookingID, "paid")
	case "failed":
		u.paymentRepo.UpdateStatus(payment.ID, "failed", req.TransactionID)
		u.bookingRepo.UpdateBookingStatus(payment.BookingID, "cancelled")
	case "expired":
		u.paymentRepo.UpdateStatus(payment.ID, "failed", "")
		u.bookingRepo.UpdateBookingStatus(payment.BookingID, "cancelled")
	}

	return nil
}
