package domain

import "github.com/google/uuid"

type Repository interface {
	FindRoleByName(name string) (*Role, error)
	FindRoleByID(id uuid.UUID) (*Role, error)
	SaveRefreshToken(token *RefreshToken) error
	FindRefreshToken(token string) (*RefreshToken, error)
	DeleteRefreshToken(token string) error
	DeleteUserRefreshTokens(userID uuid.UUID) error
}
