package migrations

import (
	"time"

	"github.com/google/uuid"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

type BookingTent struct {
	ID            uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	BookingID     uuid.UUID      `gorm:"type:uuid;not null"`
	TentID        uuid.UUID      `gorm:"type:uuid;not null"`
	PricePerNight float64        `gorm:"type:decimal(12,2);not null"`
	CreatedAt     time.Time      `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt     time.Time      `gorm:"default:CURRENT_TIMESTAMP"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

var CreateBookingTents = &gormigrate.Migration{
	ID: "20250721012",
	Migrate: func(tx *gorm.DB) error {
		return tx.Migrator().CreateTable(&BookingTent{})
	},
	Rollback: func(tx *gorm.DB) error {
		return tx.Migrator().DropTable("booking_tents")
	},
}
