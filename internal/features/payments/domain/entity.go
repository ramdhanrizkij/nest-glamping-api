package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Payment struct {
	ID            uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	BookingID     uuid.UUID      `gorm:"type:uuid;not null" json:"booking_id"`
	Amount        float64        `gorm:"type:decimal(12,2);not null" json:"amount"`
	PaymentMethod string         `gorm:"type:varchar(100)" json:"payment_method"`
	PaymentStatus string         `gorm:"type:varchar(50);default:'unpaid'" json:"payment_status"`
	PaymentDate   *time.Time     `gorm:"type:timestamp" json:"payment_date"`
	GatewayRef    string         `gorm:"type:varchar(255)" json:"gateway_ref"`
	CreatedAt     time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}
