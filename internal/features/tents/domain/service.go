package domain

import (
	"time"

	"github.com/ramdhanrizkij/nest-glamping-api/internal/features/tents/dto"
)

type Service interface {
	Create(req dto.CreateTentRequest) (*dto.TentResponse, error)
	List() ([]dto.TentResponse, error)
	ListByTentTypeID(tentTypeID string) ([]dto.TentResponse, error)
	FindByID(id string) (*dto.TentResponse, error)
	Update(id string, req dto.UpdateTentRequest) (*dto.TentResponse, error)
	Delete(id string) error
	CheckAvailability(tentTypeID string, checkIn, checkOut time.Time) ([]dto.AvailableTentResponse, error)
}
