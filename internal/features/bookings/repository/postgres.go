package repository

import (
	"github.com/google/uuid"
	"github.com/ramdhanrizkij/nest-glamping-api/internal/features/bookings/domain"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) domain.Repository {
	return &repository{db: db}
}

func (r *repository) CreateBooking(booking *domain.Booking, bookingTents []domain.BookingTent) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(booking).Error; err != nil {
			return err
		}
		if len(bookingTents) > 0 {
			if err := tx.Create(&bookingTents).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *repository) FindBookingByID(id uuid.UUID) (*domain.Booking, error) {
	var booking domain.Booking
	err := r.db.Where("id = ?", id).First(&booking).Error
	if err != nil {
		return nil, err
	}
	return &booking, nil
}

func (r *repository) FindBookingByCode(code string) (*domain.Booking, error) {
	var booking domain.Booking
	err := r.db.Where("booking_code = ?", code).First(&booking).Error
	if err != nil {
		return nil, err
	}
	return &booking, nil
}

func (r *repository) FindBookingTentsByBookingID(bookingID uuid.UUID) ([]domain.BookingTent, error) {
	var bookingTents []domain.BookingTent
	err := r.db.Where("booking_id = ?", bookingID).Find(&bookingTents).Error
	return bookingTents, err
}

func (r *repository) FindBookingsByUserID(userID uuid.UUID) ([]domain.Booking, error) {
	var bookings []domain.Booking
	err := r.db.Where("user_id = ? AND deleted_at IS NULL", userID).Order("created_at DESC").Find(&bookings).Error
	return bookings, err
}

func (r *repository) FindAllBookings() ([]domain.Booking, error) {
	var bookings []domain.Booking
	err := r.db.Where("deleted_at IS NULL").Order("created_at DESC").Find(&bookings).Error
	return bookings, err
}

func (r *repository) UpdateBookingStatus(id uuid.UUID, status string) error {
	return r.db.Model(&domain.Booking{}).Where("id = ?", id).Update("status", status).Error
}

func (r *repository) FindBookingByIDWithTents(id uuid.UUID) (*domain.Booking, []domain.BookingTent, error) {
	var booking domain.Booking
	err := r.db.Where("id = ?", id).First(&booking).Error
	if err != nil {
		return nil, nil, err
	}

	var bookingTents []domain.BookingTent
	err = r.db.Where("booking_id = ?", id).Find(&bookingTents).Error
	if err != nil {
		return nil, nil, err
	}

	return &booking, bookingTents, nil
}

func (r *repository) FindUserBookingsWithTents(userID uuid.UUID) ([]domain.Booking, error) {
	var bookings []domain.Booking
	err := r.db.Where("user_id = ? AND deleted_at IS NULL", userID).Order("created_at DESC").Find(&bookings).Error
	if err != nil {
		return nil, err
	}
	return bookings, nil
}

func (r *repository) FindAllBookingsWithTents() ([]domain.Booking, error) {
	var bookings []domain.Booking
	err := r.db.Where("deleted_at IS NULL").Order("created_at DESC").Find(&bookings).Error
	if err != nil {
		return nil, err
	}
	return bookings, nil
}
