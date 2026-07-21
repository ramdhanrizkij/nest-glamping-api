package domain

import "github.com/google/uuid"

type Repository interface {
	Create(amenity *Amenity) error
	FindAll() ([]Amenity, error)
	FindByID(id uuid.UUID) (*Amenity, error)
	Update(amenity *Amenity) error
	Delete(id uuid.UUID) error
}
