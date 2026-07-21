package domain

import (
	"github.com/ramdhanrizkij/nest-glamping-api/internal/features/tent-types/dto"
)

type Service interface {
	Create(req dto.CreateTentTypeRequest) (*dto.TentTypeResponse, error)
	List() ([]dto.TentTypeResponse, error)
	FindByID(id string) (*dto.TentTypeDetailResponse, error)
	Update(id string, req dto.UpdateTentTypeRequest) (*dto.TentTypeResponse, error)
	Delete(id string) error

	// Images
	AddImage(tentTypeID string, req dto.AddImageRequest) (*dto.TentTypeImageResponse, error)
	DeleteImage(imageID string) error
	SetPrimaryImage(tentTypeID string, imageID string) error

	// Rates
	CreateRate(tentTypeID string, req dto.CreateRateRequest) (*dto.TentTypeRateResponse, error)
	ListRates(tentTypeID string) ([]dto.TentTypeRateResponse, error)
	UpdateRate(rateID string, req dto.UpdateRateRequest) (*dto.TentTypeRateResponse, error)
	DeleteRate(rateID string) error
}
