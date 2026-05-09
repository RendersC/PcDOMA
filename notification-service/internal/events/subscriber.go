package events

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/pcdoma/notification-service/internal/domain"
	mongorepo "github.com/pcdoma/notification-service/internal/repository/mongo"
	"github.com/redis/go-redis/v9"
)

type Subscriber struct {
	redis  *redis.Client
	repo   *mongorepo.NotificationRepository
}

func NewSubscriber(rdb *redis.Client, repo *mongorepo.NotificationRepository) *Subscriber {
	return &Subscriber{redis: rdb, repo: repo}
}

func (s *Subscriber) Start(ctx context.Context) {
	pubsub := s.redis.Subscribe(ctx, "booking.events", "worker.events", "payment.events")
	ch := pubsub.Channel()

	go func() {
		defer pubsub.Close()
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-ch:
				if !ok {
					return
				}
				s.handle(ctx, msg.Channel, msg.Payload)
			}
		}
	}()

	log.Println("notification subscriber started, listening on booking.events, worker.events, payment.events")
}

func (s *Subscriber) handle(ctx context.Context, channel, payload string) {
	switch channel {
	case "booking.events":
		s.handleBookingEvent(ctx, payload)
	case "worker.events":
		s.handleWorkerEvent(ctx, payload)
	case "payment.events":
		s.handlePaymentEvent(ctx, payload)
	}
}

func (s *Subscriber) handleBookingEvent(ctx context.Context, payload string) {
	var event domain.BookingEvent
	if err := json.Unmarshal([]byte(payload), &event); err != nil {
		log.Printf("failed to unmarshal booking event: %v", err)
		return
	}

	notif, ok := BuildBookingNotification(event)
	if !ok {
		return
	}

	if err := s.repo.Create(ctx, notif); err != nil {
		log.Printf("failed to save notification: %v", err)
	} else {
		log.Printf("[EMAIL MOCK] To user %s: %s — %s", notif.UserID, notif.Title, notif.Body)
	}
}

func (s *Subscriber) handleWorkerEvent(ctx context.Context, payload string) {
	var event domain.BookingEvent
	if err := json.Unmarshal([]byte(payload), &event); err != nil {
		return
	}

	notif, ok := BuildWorkerNotification(event)
	if !ok {
		return
	}

	if err := s.repo.Create(ctx, notif); err != nil {
		log.Printf("failed to save worker notification: %v", err)
	} else {
		log.Printf("[WORKER NOTIFY] Location %s: new booking %s", event.LocationID, event.BookingID)
	}
}

func (s *Subscriber) handlePaymentEvent(ctx context.Context, payload string) {
	var event domain.PaymentEvent
	if err := json.Unmarshal([]byte(payload), &event); err != nil {
		return
	}

	notif, ok := BuildPaymentNotification(event)
	if !ok {
		return
	}
	s.repo.Create(ctx, notif)
}

// BuildBookingNotification maps a booking event to a user-facing notification.
// Returns ok=false for events that should not produce a notification.
func BuildBookingNotification(event domain.BookingEvent) (*domain.Notification, bool) {
	switch event.Event {
	case "booking.confirmed":
		return &domain.Notification{
			UserID: event.UserID,
			Type:   domain.TypeBookingConfirmed,
			Title:  "Бронирование подтверждено!",
			Body:   fmt.Sprintf("Ваш заказ принят. Итого: %.0f KZT", event.Total),
			Data:   map[string]string{"booking_id": event.BookingID},
		}, true
	case "booking.cancelled":
		return &domain.Notification{
			UserID: event.UserID,
			Type:   domain.TypeBookingCancelled,
			Title:  "Бронирование отменено",
			Body:   "Ваше бронирование было отменено",
			Data:   map[string]string{"booking_id": event.BookingID},
		}, true
	default:
		return nil, false
	}
}

// BuildWorkerNotification maps a booking event to a notification routed to the
// workers of a location. Only confirmed bookings with a location produce one.
func BuildWorkerNotification(event domain.BookingEvent) (*domain.Notification, bool) {
	if event.LocationID == "" || event.Event != "booking.confirmed" {
		return nil, false
	}
	return &domain.Notification{
		UserID: "workers:" + event.LocationID,
		Type:   domain.TypeNewBooking,
		Title:  "Новый заказ!",
		Body:   fmt.Sprintf("Заказ на %.0f KZT ожидает обработки", event.Total),
		Data:   map[string]string{"booking_id": event.BookingID, "location_id": event.LocationID},
	}, true
}

// BuildPaymentNotification maps a payment event to a user-facing notification.
func BuildPaymentNotification(event domain.PaymentEvent) (*domain.Notification, bool) {
	switch event.Event {
	case "payment.completed":
		return &domain.Notification{
			UserID: event.UserID,
			Type:   domain.TypePaymentCompleted,
			Title:  "Оплата прошла успешно",
			Body:   fmt.Sprintf("Оплачено %.0f KZT", event.Amount),
			Data:   map[string]string{"payment_id": event.PaymentID},
		}, true
	case "payment.failed":
		return &domain.Notification{
			UserID: event.UserID,
			Type:   domain.TypePaymentFailed,
			Title:  "Ошибка оплаты",
			Body:   "Не удалось обработать платёж. Попробуйте снова.",
			Data:   map[string]string{"payment_id": event.PaymentID},
		}, true
	default:
		return nil, false
	}
}
