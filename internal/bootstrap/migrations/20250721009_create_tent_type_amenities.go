package migrations

import (
	"time"

	"github.com/google/uuid"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

type TentTypeAmenity struct {
	ID         uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TentTypeID uuid.UUID      `gorm:"type:uuid;not null"`
	AmenityID  uuid.UUID      `gorm:"type:uuid;not null"`
	CreatedAt  time.Time      `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt  time.Time      `gorm:"default:CURRENT_TIMESTAMP"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

var CreateTentTypeAmenities = &gormigrate.Migration{
	ID: "20250721009",
	Migrate: func(tx *gorm.DB) error {
		return tx.Migrator().CreateTable(&TentTypeAmenity{})
	},
	Rollback: func(tx *gorm.DB) error {
		return tx.Migrator().DropTable("tent_type_amenities")
	},
}
