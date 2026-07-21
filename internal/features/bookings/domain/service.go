package domain

import "github.com/ramdhanrizkij/nest-glamping-api/internal/features/bookings/dto"

type Service interface {
	CreateBooking(userID string, req dto.CreateBookingRequest) (*dto.BookingResponse, error)
	ListMyBookings(userID string) ([]dto.BookingResponse, error)
	GetBookingDetail(bookingID, userID string, isAdmin bool) (*dto.BookingDetailResponse, error)
	CancelBooking(bookingID, userID string) error
	ListAllBookings() ([]dto.BookingResponse, error)
	ConfirmBooking(bookingID string) error
}
