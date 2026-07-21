package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/ramdhanrizkij/nest-glamping-api/internal/features/tents/domain"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) domain.Repository {
	return &repository{db: db}
}

func (r *repository) Create(tent *domain.Tent) error {
	return r.db.Create(tent).Error
}

func (r *repository) FindAll() ([]domain.Tent, error) {
	var tents []domain.Tent
	err := r.db.Find(&tents).Error
	return tents, err
}

func (r *repository) FindByID(id uuid.UUID) (*domain.Tent, error) {
	var tent domain.Tent
	err := r.db.Where("id = ?", id).First(&tent).Error
	if err != nil {
		return nil, err
	}
	return &tent, nil
}

func (r *repository) FindByTentTypeID(tentTypeID uuid.UUID) ([]domain.Tent, error) {
	var tents []domain.Tent
	err := r.db.Where("tent_type_id = ?", tentTypeID).Find(&tents).Error
	return tents, err
}

func (r *repository) Update(tent *domain.Tent) error {
	return r.db.Save(tent).Error
}

func (r *repository) Delete(id uuid.UUID) error {
	return r.db.Delete(&domain.Tent{}, id).Error
}

func (r *repository) FindAvailableTents(tentTypeID uuid.UUID, checkIn, checkOut time.Time) ([]domain.Tent, error) {
	var availableTents []domain.Tent

	err := r.db.Raw(`
		SELECT t.* FROM tents t
		WHERE t.tent_type_id = ?
		AND t.status = 'available'
		AND t.deleted_at IS NULL
		AND t.id NOT IN (
			SELECT bt.tent_id FROM booking_tents bt
			JOIN bookings b ON b.id = bt.booking_id
			WHERE bt.tent_id = t.id
			AND b.status NOT IN ('cancelled', 'completed')
			AND b.check_in_date < ?
			AND b.check_out_date > ?
		)
	`, tentTypeID, checkOut, checkIn).Scan(&availableTents).Error

	return availableTents, err
}
