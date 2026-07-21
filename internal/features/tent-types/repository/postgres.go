package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/ramdhanrizkij/nest-glamping-api/internal/features/tent-types/domain"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) domain.Repository {
	return &repository{db: db}
}

// --- TentType ---

func (r *repository) Create(tentType *domain.TentType) error {
	return r.db.Create(tentType).Error
}

func (r *repository) FindAll() ([]domain.TentType, error) {
	var tentTypes []domain.TentType
	err := r.db.Find(&tentTypes).Error
	return tentTypes, err
}

func (r *repository) FindByID(id uuid.UUID) (*domain.TentType, error) {
	var tentType domain.TentType
	err := r.db.Where("id = ?", id).First(&tentType).Error
	if err != nil {
		return nil, err
	}
	return &tentType, nil
}

func (r *repository) Update(tentType *domain.TentType) error {
	return r.db.Save(tentType).Error
}

func (r *repository) Delete(id uuid.UUID) error {
	return r.db.Delete(&domain.TentType{}, id).Error
}

// --- Rates ---

func (r *repository) CreateRate(rate *domain.TentTypeRate) error {
	return r.db.Create(rate).Error
}

func (r *repository) FindRatesByTentTypeID(tentTypeID uuid.UUID) ([]domain.TentTypeRate, error) {
	var rates []domain.TentTypeRate
	err := r.db.Where("tent_type_id = ?", tentTypeID).Order("start_date").Find(&rates).Error
	return rates, err
}

func (r *repository) FindRateByID(id uuid.UUID) (*domain.TentTypeRate, error) {
	var rate domain.TentTypeRate
	err := r.db.Where("id = ?", id).First(&rate).Error
	if err != nil {
		return nil, err
	}
	return &rate, nil
}

func (r *repository) UpdateRate(rate *domain.TentTypeRate) error {
	return r.db.Save(rate).Error
}

func (r *repository) DeleteRate(id uuid.UUID) error {
	return r.db.Delete(&domain.TentTypeRate{}, id).Error
}

func (r *repository) FindOverlappingRate(tentTypeID uuid.UUID, startDate, endDate time.Time, excludeRateID uuid.UUID) (*domain.TentTypeRate, error) {
	var rate domain.TentTypeRate
	query := r.db.Where("tent_type_id = ? AND is_active = true AND start_date < ? AND end_date > ?", tentTypeID, endDate, startDate)
	if excludeRateID != uuid.Nil {
		query = query.Where("id != ?", excludeRateID)
	}
	err := query.First(&rate).Error
	if err != nil {
		return nil, err
	}
	return &rate, nil
}

// --- Images ---

func (r *repository) CreateImage(image *domain.TentTypeImage) error {
	return r.db.Create(image).Error
}

func (r *repository) FindImagesByTentTypeID(tentTypeID uuid.UUID) ([]domain.TentTypeImage, error) {
	var images []domain.TentTypeImage
	err := r.db.Where("tent_type_id = ?", tentTypeID).Order("is_primary DESC, created_at ASC").Find(&images).Error
	return images, err
}

func (r *repository) FindImageByID(id uuid.UUID) (*domain.TentTypeImage, error) {
	var image domain.TentTypeImage
	err := r.db.Where("id = ?", id).First(&image).Error
	if err != nil {
		return nil, err
	}
	return &image, nil
}

func (r *repository) DeleteImage(id uuid.UUID) error {
	return r.db.Delete(&domain.TentTypeImage{}, id).Error
}

func (r *repository) SetPrimaryImage(tentTypeID uuid.UUID, imageID uuid.UUID) error {
	return r.db.Model(&domain.TentTypeImage{}).
		Where("tent_type_id = ? AND id = ?", tentTypeID, imageID).
		Update("is_primary", true).Error
}

func (r *repository) ClearPrimaryImage(tentTypeID uuid.UUID) error {
	return r.db.Model(&domain.TentTypeImage{}).
		Where("tent_type_id = ?", tentTypeID).
		Update("is_primary", false).Error
}

// --- Amenities ---

func (r *repository) FindAmenitiesByTentTypeID(tentTypeID uuid.UUID) ([]uuid.UUID, error) {
	var amenityIDs []uuid.UUID
	err := r.db.Model(&domain.TentTypeAmenity{}).
		Where("tent_type_id = ?", tentTypeID).
		Pluck("amenity_id", &amenityIDs).Error
	return amenityIDs, err
}

func (r *repository) SetAmenities(tentTypeID uuid.UUID, amenityIDs []uuid.UUID) error {
	if err := r.RemoveAllAmenities(tentTypeID); err != nil {
		return err
	}
	var tta []domain.TentTypeAmenity
	for _, aid := range amenityIDs {
		tta = append(tta, domain.TentTypeAmenity{
			ID:         uuid.New(),
			TentTypeID: tentTypeID,
			AmenityID:  aid,
		})
	}
	if len(tta) == 0 {
		return nil
	}
	return r.db.Create(&tta).Error
}

func (r *repository) RemoveAllAmenities(tentTypeID uuid.UUID) error {
	return r.db.Where("tent_type_id = ?", tentTypeID).Delete(&domain.TentTypeAmenity{}).Error
}
