package domain

import "github.com/google/uuid"

type Repository interface {
	Create(payment *Payment) error
	FindByID(id uuid.UUID) (*Payment, error)
	FindByBookingID(bookingID uuid.UUID) (*Payment, error)
	UpdateStatus(id uuid.UUID, status string, gatewayRef string) error
}

type BookingRepository interface {
	UpdateBookingStatus(id uuid.UUID, status string) error
	FindBookingByID(id uuid.UUID) (BookingInterface, error)
}

type BookingInterface interface {
	GetID() uuid.UUID
	GetStatus() string
}
