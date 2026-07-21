package migrations

import (
	"time"

	"github.com/google/uuid"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

type TentTypeRate struct {
	ID            uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TentTypeID    uuid.UUID      `gorm:"type:uuid;not null"`
	StartDate     time.Time      `gorm:"type:date;not null"`
	EndDate       time.Time      `gorm:"type:date;not null"`
	PricePerNight float64        `gorm:"type:decimal(12,2);not null"`
	Description   string         `gorm:"type:varchar(255)"`
	IsActive      bool           `gorm:"type:boolean;default:true"`
	CreatedAt     time.Time      `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt     time.Time      `gorm:"default:CURRENT_TIMESTAMP"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

var CreateTentTypeRates = &gormigrate.Migration{
	ID: "20250721006",
	Migrate: func(tx *gorm.DB) error {
		return tx.Migrator().CreateTable(&TentTypeRate{})
	},
	Rollback: func(tx *gorm.DB) error {
		return tx.Migrator().DropTable("tent_type_rates")
	},
}
