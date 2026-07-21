package repository

import (
	"github.com/google/uuid"
	"github.com/ramdhanrizkij/nest-glamping-api/internal/features/auth/domain"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) domain.Repository {
	return &repository{db: db}
}

func (r *repository) FindRoleByName(name string) (*domain.Role, error) {
	var role domain.Role
	err := r.db.Where("name = ?", name).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *repository) FindRoleByID(id uuid.UUID) (*domain.Role, error) {
	var role domain.Role
	err := r.db.Where("id = ?", id).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *repository) SaveRefreshToken(token *domain.RefreshToken) error {
	return r.db.Create(token).Error
}

func (r *repository) FindRefreshToken(token string) (*domain.RefreshToken, error) {
	var rt domain.RefreshToken
	err := r.db.Where("token = ? AND expires_at > NOW()", token).First(&rt).Error
	if err != nil {
		return nil, err
	}
	return &rt, nil
}

func (r *repository) DeleteRefreshToken(token string) error {
	return r.db.Where("token = ?", token).Delete(&domain.RefreshToken{}).Error
}

func (r *repository) DeleteUserRefreshTokens(userID uuid.UUID) error {
	return r.db.Where("user_id = ?", userID).Delete(&domain.RefreshToken{}).Error
}
