package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Booking struct {
	ID              uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID          uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"`
	BookingCode     string         `gorm:"type:varchar(20);uniqueIndex;not null" json:"booking_code"`
	CheckInDate     time.Time      `gorm:"type:date;not null" json:"check_in_date"`
	CheckInTime     time.Time      `gorm:"type:time" json:"check_in_time"`
	CheckOutDate    time.Time      `gorm:"type:date;not null" json:"check_out_date"`
	CheckOutTime    time.Time      `gorm:"type:time" json:"check_out_time"`
	TotalAmount     float64        `gorm:"type:decimal(12,2);not null" json:"total_amount"`
	Status          string         `gorm:"type:varchar(50);default:'pending'" json:"status"`
	GuestCount      int            `gorm:"type:int;not null;default:1" json:"guest_count"`
	SpecialRequests string         `gorm:"type:text" json:"special_requests"`
	CreatedAt       time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

type BookingTent struct {
	ID            uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	BookingID     uuid.UUID      `gorm:"type:uuid;not null" json:"booking_id"`
	TentID        uuid.UUID      `gorm:"type:uuid;not null" json:"tent_id"`
	PricePerNight float64        `gorm:"type:decimal(12,2);not null" json:"price_per_night"`
	CreatedAt     time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}
