package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Tent struct {
	ID         uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	TentTypeID uuid.UUID      `gorm:"type:uuid;not null" json:"tent_type_id"`
	Code       string         `gorm:"type:varchar(100);not null" json:"name_or_number"`
	Status     string         `gorm:"type:varchar(50);default:'available'" json:"status"`
	CreatedAt  time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}
