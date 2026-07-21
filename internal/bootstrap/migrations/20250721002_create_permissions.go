package migrations

import (
	"time"

	"github.com/google/uuid"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

type Permission struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name        string    `gorm:"type:varchar(100);uniqueIndex;not null"`
	Description string    `gorm:"type:text"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type RolePermission struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	RoleID       uuid.UUID `gorm:"type:uuid;not null"`
	PermissionID uuid.UUID `gorm:"type:uuid;not null"`
	CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

var CreatePermissions = &gormigrate.Migration{
	ID: "20250721002",
	Migrate: func(tx *gorm.DB) error {
		if err := tx.Migrator().CreateTable(&Permission{}); err != nil {
			return err
		}
		return tx.Migrator().CreateTable(&RolePermission{})
	},
	Rollback: func(tx *gorm.DB) error {
		if err := tx.Migrator().DropTable("role_permissions"); err != nil {
			return err
		}
		return tx.Migrator().DropTable("permissions")
	},
}
