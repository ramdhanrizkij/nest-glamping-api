package migrations

import (
	"time"

	"github.com/google/uuid"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

type RefreshToken struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index"`
	Token     string    `gorm:"type:varchar(500);uniqueIndex;not null"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

var CreateRefreshTokens = &gormigrate.Migration{
	ID: "20250721004",
	Migrate: func(tx *gorm.DB) error {
		return tx.Migrator().CreateTable(&RefreshToken{})
	},
	Rollback: func(tx *gorm.DB) error {
		return tx.Migrator().DropTable("refresh_tokens")
	},
}
