package usecase

import (
	"github.com/google/uuid"
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

func (u *usecase) ListAll() ([]dto.UserResponse, error) {
	users, err := u.repo.FindAll()
	if err != nil {
		return nil, appErr.Internal("failed to list users")
	}

	var result []dto.UserResponse
	for _, user := range users {
		result = append(result, dto.UserResponse{
			ID:          user.ID.String(),
			Name:        user.Name,
			Email:       user.Email,
			PhoneNumber: user.PhoneNumber,
			RoleID:      user.RoleID.String(),
		})
	}
	return result, nil
}

func (u *usecase) GetUserByID(id string) (*dto.UserResponse, error) {
	user, err := u.repo.FindByID(id)
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

func (u *usecase) UpdateUser(id string, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
	user, err := u.repo.FindByID(id)
	if err != nil {
		return nil, appErr.NotFound("user not found")
	}

	if req.Name != "" {
		user.Name = req.Name
	}
	if req.PhoneNumber != "" {
		user.PhoneNumber = req.PhoneNumber
	}
	if req.RoleID != "" {
		roleUUID, err := uuid.Parse(req.RoleID)
		if err != nil {
			return nil, appErr.BadRequest("invalid role id")
		}
		user.RoleID = roleUUID
	}

	if err := u.repo.Update(user); err != nil {
		return nil, appErr.Internal("failed to update user")
	}

	return &dto.UserResponse{
		ID:          user.ID.String(),
		Name:        user.Name,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		RoleID:      user.RoleID.String(),
	}, nil
}

func (u *usecase) DeleteUser(id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return appErr.BadRequest("invalid user id")
	}

	if _, err := u.repo.FindByID(id); err != nil {
		return appErr.NotFound("user not found")
	}

	return u.repo.Delete(uid)
}
