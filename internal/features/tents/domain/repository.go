package domain

import (
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	Create(tent *Tent) error
	FindAll() ([]Tent, error)
	FindByID(id uuid.UUID) (*Tent, error)
	FindByTentTypeID(tentTypeID uuid.UUID) ([]Tent, error)
	Update(tent *Tent) error
	Delete(id uuid.UUID) error
	FindAvailableTents(tentTypeID uuid.UUID, checkIn, checkOut time.Time) ([]Tent, error)
}
