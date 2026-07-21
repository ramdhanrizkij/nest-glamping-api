package migrations

import (
	"time"

	"github.com/google/uuid"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

type TentTypeImage struct {
	ID         uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TentTypeID uuid.UUID      `gorm:"type:uuid;not null"`
	ImageURL   string         `gorm:"type:varchar(255);not null"`
	IsPrimary  bool           `gorm:"type:boolean;default:false"`
	CreatedAt  time.Time      `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt  time.Time      `gorm:"default:CURRENT_TIMESTAMP"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

var CreateTentTypeImages = &gormigrate.Migration{
	ID: "20250721007",
	Migrate: func(tx *gorm.DB) error {
		return tx.Migrator().CreateTable(&TentTypeImage{})
	},
	Rollback: func(tx *gorm.DB) error {
		return tx.Migrator().DropTable("tent_type_images")
	},
}
