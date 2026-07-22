package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/ramdhanrizkij/nest-glamping-api/internal/features/payments/domain"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) domain.Repository {
	return &repository{db: db}
}

func (r *repository) Create(payment *domain.Payment) error {
	return r.db.Create(payment).Error
}

func (r *repository) FindByID(id uuid.UUID) (*domain.Payment, error) {
	var payment domain.Payment
	err := r.db.Where("id = ?", id).First(&payment).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *repository) FindByBookingID(bookingID uuid.UUID) (*domain.Payment, error) {
	var payment domain.Payment
	err := r.db.Where("booking_id = ?", bookingID).First(&payment).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *repository) UpdateStatus(id uuid.UUID, status string, gatewayRef string) error {
	updates := map[string]interface{}{
		"payment_status": status,
	}
	if gatewayRef != "" {
		updates["gateway_ref"] = gatewayRef
	}
	if status == "paid" {
		now := time.Now()
		updates["payment_date"] = now
	}
	return r.db.Model(&domain.Payment{}).Where("id = ?", id).Updates(updates).Error
}
