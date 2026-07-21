package usecase

import (
	"github.com/google/uuid"
	"github.com/ramdhanrizkij/nest-glamping-api/internal/features/amenities/domain"
	"github.com/ramdhanrizkij/nest-glamping-api/internal/features/amenities/dto"
	appErr "github.com/ramdhanrizkij/nest-glamping-api/internal/shared/errors"
)

type usecase struct {
	repo domain.Repository
}

func NewUsecase(repo domain.Repository) domain.Service {
	return &usecase{repo: repo}
}

func (u *usecase) Create(req dto.CreateAmenityRequest) (*dto.AmenityResponse, error) {
	amenity := &domain.Amenity{
		ID:          uuid.New(),
		Name:        req.Name,
		IconURL:     req.IconURL,
		Description: req.Description,
	}

	if err := u.repo.Create(amenity); err != nil {
		return nil, appErr.Internal("failed to create amenity")
	}

	return &dto.AmenityResponse{
		ID:          amenity.ID.String(),
		Name:        amenity.Name,
		IconURL:     amenity.IconURL,
		Description: amenity.Description,
	}, nil
}

func (u *usecase) List() ([]dto.AmenityResponse, error) {
	amenities, err := u.repo.FindAll()
	if err != nil {
		return nil, appErr.Internal("failed to list amenities")
	}

	var result []dto.AmenityResponse
	for _, a := range amenities {
		result = append(result, dto.AmenityResponse{
			ID:          a.ID.String(),
			Name:        a.Name,
			IconURL:     a.IconURL,
			Description: a.Description,
		})
	}
	return result, nil
}

func (u *usecase) FindByID(id string) (*dto.AmenityResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, appErr.BadRequest("invalid amenity id")
	}

	amenity, err := u.repo.FindByID(uid)
	if err != nil {
		return nil, appErr.NotFound("amenity not found")
	}

	return &dto.AmenityResponse{
		ID:          amenity.ID.String(),
		Name:        amenity.Name,
		IconURL:     amenity.IconURL,
		Description: amenity.Description,
	}, nil
}

func (u *usecase) Update(id string, req dto.UpdateAmenityRequest) (*dto.AmenityResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, appErr.BadRequest("invalid amenity id")
	}

	amenity, err := u.repo.FindByID(uid)
	if err != nil {
		return nil, appErr.NotFound("amenity not found")
	}

	if req.Name != "" {
		amenity.Name = req.Name
	}
	if req.IconURL != "" {
		amenity.IconURL = req.IconURL
	}
	if req.Description != "" {
		amenity.Description = req.Description
	}

	if err := u.repo.Update(amenity); err != nil {
		return nil, appErr.Internal("failed to update amenity")
	}

	return &dto.AmenityResponse{
		ID:          amenity.ID.String(),
		Name:        amenity.Name,
		IconURL:     amenity.IconURL,
		Description: amenity.Description,
	}, nil
}

func (u *usecase) Delete(id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return appErr.BadRequest("invalid amenity id")
	}

	if _, err := u.repo.FindByID(uid); err != nil {
		return appErr.NotFound("amenity not found")
	}

	return u.repo.Delete(uid)
}
