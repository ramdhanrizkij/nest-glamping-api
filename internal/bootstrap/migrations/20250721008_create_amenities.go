package migrations

import (
	"time"

	"github.com/google/uuid"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

type Amenity struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name        string         `gorm:"type:varchar(100);not null"`
	IconURL     string         `gorm:"type:varchar(255)"`
	Description string         `gorm:"type:text"`
	CreatedAt   time.Time      `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time      `gorm:"default:CURRENT_TIMESTAMP"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

var CreateAmenities = &gormigrate.Migration{
	ID: "20250721008",
	Migrate: func(tx *gorm.DB) error {
		return tx.Migrator().CreateTable(&Amenity{})
	},
	Rollback: func(tx *gorm.DB) error {
		return tx.Migrator().DropTable("amenities")
	},
}
