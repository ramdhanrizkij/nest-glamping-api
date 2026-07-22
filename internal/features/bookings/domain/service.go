package domain

import "github.com/ramdhanrizkij/nest-glamping-api/internal/features/bookings/dto"

type Service interface {
	CreateBooking(userID string, req dto.CreateBookingRequest) (*dto.BookingResponse, error)
	ListMyBookings(userID string, page, perPage int) (*dto.BookingListResponse, error)
	GetBookingDetail(bookingID, userID string, isAdmin bool) (*dto.BookingDetailResponse, error)
	CancelBooking(bookingID, userID string) error
	ListAllBookings(page, perPage int) (*dto.BookingListResponse, error)
	ConfirmBooking(bookingID string) error
}
