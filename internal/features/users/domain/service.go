package domain

import "github.com/ramdhanrizkij/nest-glamping-api/internal/features/users/dto"

type Service interface {
	GetProfile(userID string) (*dto.UserResponse, error)
	UpdateProfile(userID string, req dto.UpdateProfileRequest) (*dto.UserResponse, error)
}
