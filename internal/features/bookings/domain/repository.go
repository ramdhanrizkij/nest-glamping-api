package domain

import "github.com/google/uuid"

type Repository interface {
	CreateBooking(booking *Booking, bookingTents []BookingTent) error
	FindBookingByID(id uuid.UUID) (*Booking, error)
	FindBookingByCode(code string) (*Booking, error)
	FindBookingTentsByBookingID(bookingID uuid.UUID) ([]BookingTent, error)
	FindBookingsByUserID(userID uuid.UUID) ([]Booking, error)
	FindAllBookings() ([]Booking, error)
	UpdateBookingStatus(id uuid.UUID, status string) error
	FindBookingByIDWithTents(id uuid.UUID) (*Booking, []BookingTent, error)
	FindUserBookingsWithTents(userID uuid.UUID) ([]Booking, error)
	FindAllBookingsWithTents() ([]Booking, error)
}
