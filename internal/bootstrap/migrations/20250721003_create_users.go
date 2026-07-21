package migrations

import (
	"time"

	"github.com/google/uuid"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	RoleID       uuid.UUID `gorm:"type:uuid;not null"`
	Name         string    `gorm:"type:varchar(255);not null"`
	Email        string    `gorm:"type:varchar(255);uniqueIndex;not null"`
	PasswordHash string    `gorm:"type:varchar(255);not null"`
	PhoneNumber  string    `gorm:"type:varchar(20)"`
	CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

var CreateUsers = &gormigrate.Migration{
	ID: "20250721003",
	Migrate: func(tx *gorm.DB) error {
		return tx.Migrator().CreateTable(&User{})
	},
	Rollback: func(tx *gorm.DB) error {
		return tx.Migrator().DropTable("users")
	},
}
