package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TentType struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name        string         `gorm:"type:varchar(255);not null" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	Capacity    int            `gorm:"type:int;not null;default:2" json:"capacity"`
	BasePrice   float64        `gorm:"type:decimal(12,2);not null" json:"base_price"`
	CreatedAt   time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type TentTypeRate struct {
	ID            uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	TentTypeID    uuid.UUID      `gorm:"type:uuid;not null" json:"tent_type_id"`
	StartDate     time.Time      `gorm:"type:date;not null" json:"start_date"`
	EndDate       time.Time      `gorm:"type:date;not null" json:"end_date"`
	PricePerNight float64        `gorm:"type:decimal(12,2);not null" json:"price_per_night"`
	Description   string         `gorm:"type:varchar(255)" json:"description"`
	IsActive      bool           `gorm:"type:boolean;default:true" json:"is_active"`
	CreatedAt     time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

type TentTypeImage struct {
	ID         uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	TentTypeID uuid.UUID      `gorm:"type:uuid;not null" json:"tent_type_id"`
	ImageURL   string         `gorm:"type:varchar(255);not null" json:"image_url"`
	IsPrimary  bool           `gorm:"type:boolean;default:false" json:"is_primary"`
	CreatedAt  time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

type TentTypeAmenity struct {
	ID         uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	TentTypeID uuid.UUID      `gorm:"type:uuid;not null" json:"tent_type_id"`
	AmenityID  uuid.UUID      `gorm:"type:uuid;not null" json:"amenity_id"`
	CreatedAt  time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}
