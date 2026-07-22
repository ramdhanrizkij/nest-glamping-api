package migrations

import (
	"time"

	"github.com/google/uuid"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

type Payment struct {
	ID            uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	BookingID     uuid.UUID      `gorm:"type:uuid;not null"`
	Amount        float64        `gorm:"type:decimal(12,2);not null"`
	PaymentMethod string         `gorm:"type:varchar(100)"`
	PaymentStatus string         `gorm:"type:varchar(50);default:'unpaid'"`
	PaymentDate   *time.Time     `gorm:"type:timestamp"`
	GatewayRef    string         `gorm:"type:varchar(255)"`
	CreatedAt     time.Time      `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt     time.Time      `gorm:"default:CURRENT_TIMESTAMP"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

var CreatePayments = &gormigrate.Migration{
	ID: "20250721013",
	Migrate: func(tx *gorm.DB) error {
		return tx.Migrator().CreateTable(&Payment{})
	},
	Rollback: func(tx *gorm.DB) error {
		return tx.Migrator().DropTable("payments")
	},
}
