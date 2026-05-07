package domain

import (
	"time"

	"github.com/google/uuid"
)

type BookingStatus string
type BookingType string
type RentalType string

const (
	StatusPendingPayment BookingStatus = "pending_payment"
	StatusPaid           BookingStatus = "paid"
	StatusPendingWorker  BookingStatus = "pending_worker"
	StatusAccepted       BookingStatus = "accepted"
	StatusRejected       BookingStatus = "rejected"
	StatusDelivering     BookingStatus = "delivering"
	StatusActive         BookingStatus = "active"
	StatusCompleted      BookingStatus = "completed"
	StatusCancelled      BookingStatus = "cancelled"

	TypeCustom BookingType = "custom"
	TypeSetup  BookingType = "setup"

	RentalHourly RentalType = "hourly"
	RentalDaily  RentalType = "daily"
)

type Booking struct {
	ID              uuid.UUID     `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID          uuid.UUID     `gorm:"type:uuid;not null;index" json:"user_id"`
	BookingType     BookingType   `gorm:"type:varchar(10);not null" json:"booking_type"`
	PCID            string        `gorm:"type:varchar(24)" json:"pc_id,omitempty"`
	SetupID         string        `gorm:"type:varchar(24)" json:"setup_id,omitempty"`
	LocationID      string        `gorm:"type:varchar(50)" json:"location_id"`
	WorkerID        *uuid.UUID    `gorm:"type:uuid" json:"worker_id,omitempty"`
	RentalType      RentalType    `gorm:"type:varchar(10)" json:"rental_type"`
	StartTime       time.Time     `gorm:"not null" json:"start_time"`
	EndTime         time.Time     `gorm:"not null" json:"end_time"`
	BasePrice       float64       `gorm:"not null" json:"base_price"`
	ExtrasPrice     float64       `json:"extras_price"`
	Discount        float64       `json:"discount"`
	TotalPrice      float64       `gorm:"not null" json:"total_price"`
	Currency        string        `gorm:"default:'KZT'" json:"currency"`
	Status          BookingStatus `gorm:"type:varchar(30);default:'pending_payment'" json:"status"`
	RejectionReason string        `json:"rejection_reason,omitempty"`
	PaymentID       *uuid.UUID    `gorm:"type:uuid" json:"payment_id,omitempty"`
	AcceptedAt      *time.Time    `json:"accepted_at,omitempty"`
	DeliveredAt     *time.Time    `json:"delivered_at,omitempty"`
	CompletedAt     *time.Time    `json:"completed_at,omitempty"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`

	Peripherals []BookingPeripheral `gorm:"foreignKey:BookingID" json:"peripherals,omitempty"`
}

type BookingPeripheral struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	BookingID    uuid.UUID `gorm:"type:uuid;not null;index" json:"booking_id"`
	PeripheralID string    `gorm:"type:varchar(24)" json:"peripheral_id"`
	PricePerUnit float64   `json:"price_per_unit"`
}

type CreateBookingRequest struct {
	BookingType   BookingType `json:"booking_type" binding:"required"`
	PCID          string      `json:"pc_id"`
	SetupID       string      `json:"setup_id"`
	PeripheralIDs []string    `json:"peripheral_ids"`
	StartTime     time.Time   `json:"start_time" binding:"required"`
	EndTime       time.Time   `json:"end_time" binding:"required"`
	RentalType    RentalType  `json:"rental_type" binding:"required"`
}

type BookingEvent struct {
	Event     string    `json:"event"`
	BookingID uuid.UUID `json:"booking_id"`
	UserID    uuid.UUID `json:"user_id"`
	LocationID string   `json:"location_id"`
	PCName    string    `json:"pc_name"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Items     []string  `json:"items"`
	Total     float64   `json:"total"`
}

type WorkerActionRequest struct {
	RejectionReason string `json:"rejection_reason"`
}
