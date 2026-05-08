package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/pcdoma/booking-service/internal/domain"
	"gorm.io/gorm"
)

type BookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) *BookingRepository {
	return &BookingRepository{db: db}
}

func (r *BookingRepository) Create(ctx context.Context, b *domain.Booking) error {
	return r.db.WithContext(ctx).Create(b).Error
}

func (r *BookingRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Booking, error) {
	var b domain.Booking
	err := r.db.WithContext(ctx).Preload("Peripherals").First(&b, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *BookingRepository) ListByUser(ctx context.Context, userID uuid.UUID, status string, limit, offset int) ([]*domain.Booking, int64, error) {
	var bookings []*domain.Booking
	var total int64
	query := r.db.WithContext(ctx).Where("user_id = ?", userID)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	query.Model(&domain.Booking{}).Count(&total)
	err := query.Preload("Peripherals").Order("created_at DESC").Limit(limit).Offset(offset).Find(&bookings).Error
	return bookings, total, err
}

func (r *BookingRepository) ListByLocation(ctx context.Context, locationID string, status string, limit, offset int) ([]*domain.Booking, int64, error) {
	var bookings []*domain.Booking
	var total int64
	query := r.db.WithContext(ctx)
	if locationID != "" {
		query = query.Where("location_id = ?", locationID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	query.Model(&domain.Booking{}).Count(&total)
	err := query.Preload("Peripherals").Order("created_at DESC").Limit(limit).Offset(offset).Find(&bookings).Error
	return bookings, total, err
}

func (r *BookingRepository) ListAll(ctx context.Context, limit, offset int) ([]*domain.Booking, int64, error) {
	var bookings []*domain.Booking
	var total int64
	r.db.WithContext(ctx).Model(&domain.Booking{}).Count(&total)
	err := r.db.WithContext(ctx).Preload("Peripherals").Order("created_at DESC").Limit(limit).Offset(offset).Find(&bookings).Error
	return bookings, total, err
}

func (r *BookingRepository) Update(ctx context.Context, b *domain.Booking) error {
	return r.db.WithContext(ctx).Save(b).Error
}
