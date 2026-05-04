package domain

import (
	"time"

	"github.com/google/uuid"
)

type PaymentStatus string
type PaymentMethod string

const (
	StatusPending    PaymentStatus = "pending"
	StatusProcessing PaymentStatus = "processing"
	StatusCompleted  PaymentStatus = "completed"
	StatusFailed     PaymentStatus = "failed"
	StatusRefunded   PaymentStatus = "refunded"

	MethodCard  PaymentMethod = "card"
	MethodKaspi PaymentMethod = "kaspi"
)

type Payment struct {
	ID          uuid.UUID     `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	BookingID   uuid.UUID     `gorm:"type:uuid;not null;index" json:"booking_id"`
	UserID      uuid.UUID     `gorm:"type:uuid;not null;index" json:"user_id"`
	Amount      float64       `gorm:"not null" json:"amount"`
	Currency    string        `gorm:"default:'KZT'" json:"currency"`
	Method      PaymentMethod `gorm:"type:varchar(20);not null" json:"method"`
	Status      PaymentStatus `gorm:"type:varchar(20);default:'pending'" json:"status"`
	ProcessedAt *time.Time    `json:"processed_at,omitempty"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

type CreatePaymentRequest struct {
	BookingID uuid.UUID     `json:"booking_id" binding:"required"`
	Amount    float64       `json:"amount" binding:"required,gt=0"`
	Currency  string        `json:"currency"`
	Method    PaymentMethod `json:"method" binding:"required"`
}

type ProcessPaymentRequest struct {
	Method PaymentMethod `json:"method" binding:"required"`
}

// Internal request from booking-service
type InternalCreatePaymentRequest struct {
	BookingID uuid.UUID `json:"booking_id" binding:"required"`
	UserID    uuid.UUID `json:"user_id" binding:"required"`
	Amount    float64   `json:"amount" binding:"required,gt=0"`
	Currency  string    `json:"currency"`
}
