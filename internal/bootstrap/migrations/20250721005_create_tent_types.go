package migrations

import (
	"time"

	"github.com/google/uuid"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

type TentType struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name        string         `gorm:"type:varchar(255);not null"`
	Description string         `gorm:"type:text"`
	Capacity    int            `gorm:"type:int;not null;default:2"`
	BasePrice   float64        `gorm:"type:decimal(12,2);not null"`
	CreatedAt   time.Time      `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time      `gorm:"default:CURRENT_TIMESTAMP"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

var CreateTentTypes = &gormigrate.Migration{
	ID: "20250721005",
	Migrate: func(tx *gorm.DB) error {
		return tx.Migrator().CreateTable(&TentType{})
	},
	Rollback: func(tx *gorm.DB) error {
		return tx.Migrator().DropTable("tent_types")
	},
}
