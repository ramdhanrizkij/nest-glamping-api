package usecase

import (
	"github.com/ramdhanrizkij/nest-glamping-api/internal/features/users/domain"
	"github.com/ramdhanrizkij/nest-glamping-api/internal/features/users/dto"
	appErr "github.com/ramdhanrizkij/nest-glamping-api/internal/shared/errors"
)

type usecase struct {
	repo domain.Repository
}

func NewUsecase(repo domain.Repository) domain.Service {
	return &usecase{repo: repo}
}

func (u *usecase) GetProfile(userID string) (*dto.UserResponse, error) {
	user, err := u.repo.FindByID(userID)
	if err != nil {
		return nil, appErr.NotFound("user not found")
	}

	return &dto.UserResponse{
		ID:          user.ID.String(),
		Name:        user.Name,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		RoleID:      user.RoleID.String(),
	}, nil
}

func (u *usecase) UpdateProfile(userID string, req dto.UpdateProfileRequest) (*dto.UserResponse, error) {
	user, err := u.repo.FindByID(userID)
	if err != nil {
		return nil, appErr.NotFound("user not found")
	}

	if req.Name != "" {
		user.Name = req.Name
	}
	if req.PhoneNumber != "" {
		user.PhoneNumber = req.PhoneNumber
	}

	if err := u.repo.Update(user); err != nil {
		return nil, appErr.Internal("failed to update profile")
	}

	return &dto.UserResponse{
		ID:          user.ID.String(),
		Name:        user.Name,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		RoleID:      user.RoleID.String(),
	}, nil
}
