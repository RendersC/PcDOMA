package events_test

import (
	"testing"

	"github.com/pcdoma/notification-service/internal/domain"
	"github.com/pcdoma/notification-service/internal/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildBookingNotification_Confirmed(t *testing.T) {
	notif, ok := events.BuildBookingNotification(domain.BookingEvent{
		Event:     "booking.confirmed",
		BookingID: "b-1",
		UserID:    "u-1",
		Total:     12000,
	})

	require.True(t, ok)
	assert.Equal(t, "u-1", notif.UserID)
	assert.Equal(t, domain.TypeBookingConfirmed, notif.Type)
	assert.Contains(t, notif.Body, "12000")
	assert.Equal(t, "b-1", notif.Data["booking_id"])
}

func TestBuildBookingNotification_Cancelled(t *testing.T) {
	notif, ok := events.BuildBookingNotification(domain.BookingEvent{
		Event:     "booking.cancelled",
		BookingID: "b-2",
		UserID:    "u-2",
	})

	require.True(t, ok)
	assert.Equal(t, domain.TypeBookingCancelled, notif.Type)
	assert.Equal(t, "b-2", notif.Data["booking_id"])
}

func TestBuildBookingNotification_UnknownEventIgnored(t *testing.T) {
	_, ok := events.BuildBookingNotification(domain.BookingEvent{Event: "booking.something"})
	assert.False(t, ok)
}

func TestBuildWorkerNotification_RoutesToLocation(t *testing.T) {
	notif, ok := events.BuildWorkerNotification(domain.BookingEvent{
		Event:      "booking.confirmed",
		BookingID:  "b-3",
		LocationID: "esil",
		Total:      5000,
	})

	require.True(t, ok)
	assert.Equal(t, "workers:esil", notif.UserID)
	assert.Equal(t, domain.TypeNewBooking, notif.Type)
	assert.Equal(t, "esil", notif.Data["location_id"])
}

func TestBuildWorkerNotification_NoLocationIgnored(t *testing.T) {
	_, ok := events.BuildWorkerNotification(domain.BookingEvent{
		Event: "booking.confirmed",
	})
	assert.False(t, ok)
}

func TestBuildWorkerNotification_NonConfirmedIgnored(t *testing.T) {
	_, ok := events.BuildWorkerNotification(domain.BookingEvent{
		Event:      "booking.cancelled",
		LocationID: "esil",
	})
	assert.False(t, ok)
}

func TestBuildPaymentNotification_Completed(t *testing.T) {
	notif, ok := events.BuildPaymentNotification(domain.PaymentEvent{
		Event:     "payment.completed",
		PaymentID: "p-1",
		UserID:    "u-1",
		Amount:    9000,
	})

	require.True(t, ok)
	assert.Equal(t, domain.TypePaymentCompleted, notif.Type)
	assert.Contains(t, notif.Body, "9000")
	assert.Equal(t, "p-1", notif.Data["payment_id"])
}

func TestBuildPaymentNotification_Failed(t *testing.T) {
	notif, ok := events.BuildPaymentNotification(domain.PaymentEvent{
		Event:     "payment.failed",
		PaymentID: "p-2",
		UserID:    "u-2",
	})

	require.True(t, ok)
	assert.Equal(t, domain.TypePaymentFailed, notif.Type)
}

func TestBuildPaymentNotification_UnknownIgnored(t *testing.T) {
	_, ok := events.BuildPaymentNotification(domain.PaymentEvent{Event: "payment.pending"})
	assert.False(t, ok)
}
