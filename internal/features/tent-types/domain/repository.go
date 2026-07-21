package domain

import (
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	Create(tentType *TentType) error
	FindAll() ([]TentType, error)
	FindByID(id uuid.UUID) (*TentType, error)
	Update(tentType *TentType) error
	Delete(id uuid.UUID) error

	// Rates
	CreateRate(rate *TentTypeRate) error
	FindRatesByTentTypeID(tentTypeID uuid.UUID) ([]TentTypeRate, error)
	FindRateByID(id uuid.UUID) (*TentTypeRate, error)
	UpdateRate(rate *TentTypeRate) error
	DeleteRate(id uuid.UUID) error
	FindOverlappingRate(tentTypeID uuid.UUID, startDate, endDate time.Time, excludeRateID uuid.UUID) (*TentTypeRate, error)

	// Images
	CreateImage(image *TentTypeImage) error
	FindImagesByTentTypeID(tentTypeID uuid.UUID) ([]TentTypeImage, error)
	FindImageByID(id uuid.UUID) (*TentTypeImage, error)
	DeleteImage(id uuid.UUID) error
	SetPrimaryImage(tentTypeID uuid.UUID, imageID uuid.UUID) error
	ClearPrimaryImage(tentTypeID uuid.UUID) error

	// Amenities
	FindAmenitiesByTentTypeID(tentTypeID uuid.UUID) ([]uuid.UUID, error)
	SetAmenities(tentTypeID uuid.UUID, amenityIDs []uuid.UUID) error
	RemoveAllAmenities(tentTypeID uuid.UUID) error
}
