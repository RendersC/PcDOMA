package events

import (
	"context"
	"encoding/json"
	"log"

	"github.com/pcdoma/booking-service/internal/domain"
	"github.com/redis/go-redis/v9"
)

const (
	ChannelBooking = "booking.events"
	ChannelWorker  = "worker.events"
)

type Publisher struct {
	redis *redis.Client
}

func NewPublisher(rdb *redis.Client) *Publisher {
	return &Publisher{redis: rdb}
}

func (p *Publisher) PublishBookingConfirmed(ctx context.Context, event domain.BookingEvent) {
	event.Event = "booking.confirmed"
	p.publish(ctx, ChannelBooking, event)
	p.publish(ctx, ChannelWorker, event)
}

func (p *Publisher) PublishBookingCancelled(ctx context.Context, event domain.BookingEvent) {
	event.Event = "booking.cancelled"
	p.publish(ctx, ChannelBooking, event)
	p.publish(ctx, ChannelWorker, event)
}

func (p *Publisher) publish(ctx context.Context, channel string, payload interface{}) {
	if p.redis == nil {
		return
	}
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("failed to marshal event: %v", err)
		return
	}
	if err := p.redis.Publish(ctx, channel, data).Err(); err != nil {
		log.Printf("failed to publish to %s: %v", channel, err)
	}
}
