package domain

import "github.com/ramdhanrizkij/nest-glamping-api/internal/features/auth/dto"

type Service interface {
	Register(req dto.RegisterRequest) (*dto.AuthResponse, error)
	Login(req dto.LoginRequest) (*dto.AuthResponse, error)
	RefreshToken(refreshToken string) (*dto.AuthResponse, error)
	Logout(refreshToken string) error
}
