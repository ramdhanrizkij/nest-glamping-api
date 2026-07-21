package repository

import (
	"github.com/google/uuid"
	"github.com/ramdhanrizkij/nest-glamping-api/internal/features/amenities/domain"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) domain.Repository {
	return &repository{db: db}
}

func (r *repository) Create(amenity *domain.Amenity) error {
	return r.db.Create(amenity).Error
}

func (r *repository) FindAll() ([]domain.Amenity, error) {
	var amenities []domain.Amenity
	err := r.db.Find(&amenities).Error
	return amenities, err
}

func (r *repository) FindByID(id uuid.UUID) (*domain.Amenity, error) {
	var amenity domain.Amenity
	err := r.db.Where("id = ?", id).First(&amenity).Error
	if err != nil {
		return nil, err
	}
	return &amenity, nil
}

func (r *repository) Update(amenity *domain.Amenity) error {
	return r.db.Save(amenity).Error
}

func (r *repository) Delete(id uuid.UUID) error {
	return r.db.Delete(&domain.Amenity{}, id).Error
}
