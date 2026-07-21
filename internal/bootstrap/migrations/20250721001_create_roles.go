package migrations

import (
	"time"

	"github.com/google/uuid"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

type Role struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name        string    `gorm:"type:varchar(50);uniqueIndex;not null"`
	Description string    `gorm:"type:text"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

var CreateRoles = &gormigrate.Migration{
	ID: "20250721001",
	Migrate: func(tx *gorm.DB) error {
		if err := tx.Migrator().CreateTable(&Role{}); err != nil {
			return err
		}
		roles := []Role{
			{ID: uuid.New(), Name: "admin", Description: "System administrator"},
			{ID: uuid.New(), Name: "customer", Description: "Regular customer"},
		}
		return tx.Create(&roles).Error
	},
	Rollback: func(tx *gorm.DB) error {
		return tx.Migrator().DropTable("roles")
	},
}
