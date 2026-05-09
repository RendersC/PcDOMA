package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NotificationType string

const (
	TypeBookingConfirmed NotificationType = "booking_confirmed"
	TypeBookingCancelled NotificationType = "booking_cancelled"
	TypePaymentCompleted NotificationType = "payment_completed"
	TypePaymentFailed    NotificationType = "payment_failed"
	TypeNewBooking       NotificationType = "new_booking"       // for workers
	TypeOrderAccepted    NotificationType = "order_accepted"
	TypeOrderRejected    NotificationType = "order_rejected"
	TypeOrderDelivering  NotificationType = "order_delivering"
	TypeOrderCompleted   NotificationType = "order_completed"
)

type Notification struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    string             `bson:"user_id" json:"user_id"`
	Type      NotificationType   `bson:"type" json:"type"`
	Title     string             `bson:"title" json:"title"`
	Body      string             `bson:"body" json:"body"`
	Data      map[string]string  `bson:"data" json:"data"`
	IsRead    bool               `bson:"is_read" json:"is_read"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}

type BookingEvent struct {
	Event      string    `json:"event"`
	BookingID  string    `json:"booking_id"`
	UserID     string    `json:"user_id"`
	LocationID string    `json:"location_id"`
	PCName     string    `json:"pc_name"`
	StartTime  time.Time `json:"start_time"`
	Total      float64   `json:"total"`
}

type PaymentEvent struct {
	Event     string  `json:"event"`
	PaymentID string  `json:"payment_id"`
	BookingID string  `json:"booking_id"`
	UserID    string  `json:"user_id"`
	Amount    float64 `json:"amount"`
}
