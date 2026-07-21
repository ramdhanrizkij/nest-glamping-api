package migrations

import (
	"time"

	"github.com/google/uuid"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

type Tent struct {
	ID         uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TentTypeID uuid.UUID      `gorm:"type:uuid;not null"`
	NameOrNum  string         `gorm:"type:varchar(100);not null"`
	Status     string         `gorm:"type:varchar(50);default:'available'"`
	CreatedAt  time.Time      `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt  time.Time      `gorm:"default:CURRENT_TIMESTAMP"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

var CreateTents = &gormigrate.Migration{
	ID: "20250721010",
	Migrate: func(tx *gorm.DB) error {
		return tx.Migrator().CreateTable(&Tent{})
	},
	Rollback: func(tx *gorm.DB) error {
		return tx.Migrator().DropTable("tents")
	},
}
