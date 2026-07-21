package domain

import "github.com/ramdhanrizkij/nest-glamping-api/internal/features/amenities/dto"

type Service interface {
	Create(req dto.CreateAmenityRequest) (*dto.AmenityResponse, error)
	List() ([]dto.AmenityResponse, error)
	FindByID(id string) (*dto.AmenityResponse, error)
	Update(id string, req dto.UpdateAmenityRequest) (*dto.AmenityResponse, error)
	Delete(id string) error
}
