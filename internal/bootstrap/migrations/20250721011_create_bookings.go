package migrations

import (
	"time"

	"github.com/google/uuid"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

type Booking struct {
	ID             uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID         uuid.UUID      `gorm:"type:uuid;not null"`
	BookingCode    string         `gorm:"type:varchar(20);uniqueIndex;not null"`
	CheckInDate    time.Time      `gorm:"type:date;not null"`
	CheckInTime    time.Time      `gorm:"type:time"`
	CheckOutDate   time.Time      `gorm:"type:date;not null"`
	CheckOutTime   time.Time      `gorm:"type:time"`
	TotalAmount    float64        `gorm:"type:decimal(12,2);not null"`
	Status         string         `gorm:"type:varchar(50);default:'pending'"`
	GuestCount     int            `gorm:"type:int;not null;default:1"`
	SpecialRequests string        `gorm:"type:text"`
	CreatedAt      time.Time      `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt      time.Time      `gorm:"default:CURRENT_TIMESTAMP"`
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

var CreateBookings = &gormigrate.Migration{
	ID: "20250721011",
	Migrate: func(tx *gorm.DB) error {
		return tx.Migrator().CreateTable(&Booking{})
	},
	Rollback: func(tx *gorm.DB) error {
		return tx.Migrator().DropTable("bookings")
	},
}
