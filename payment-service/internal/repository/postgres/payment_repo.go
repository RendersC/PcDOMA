package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/pcdoma/payment-service/internal/domain"
	"gorm.io/gorm"
)

type paymentRepo struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) *paymentRepo {
	return &paymentRepo{db: db}
}

func (r *paymentRepo) Create(ctx context.Context, p *domain.Payment) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *paymentRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Payment, error) {
	var p domain.Payment
	err := r.db.WithContext(ctx).First(&p, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *paymentRepo) GetByBookingID(ctx context.Context, bookingID uuid.UUID) (*domain.Payment, error) {
	var p domain.Payment
	err := r.db.WithContext(ctx).First(&p, "booking_id = ?", bookingID).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *paymentRepo) ListByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*domain.Payment, int64, error) {
	var payments []*domain.Payment
	var total int64
	r.db.WithContext(ctx).Model(&domain.Payment{}).Where("user_id = ?", userID).Count(&total)
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).
		Order("created_at DESC").Limit(limit).Offset(offset).Find(&payments).Error
	return payments, total, err
}

func (r *paymentRepo) ListAll(ctx context.Context, limit, offset int) ([]*domain.Payment, int64, error) {
	var payments []*domain.Payment
	var total int64
	r.db.WithContext(ctx).Model(&domain.Payment{}).Count(&total)
	err := r.db.WithContext(ctx).Order("created_at DESC").Limit(limit).Offset(offset).Find(&payments).Error
	return payments, total, err
}

func (r *paymentRepo) Update(ctx context.Context, p *domain.Payment) error {
	return r.db.WithContext(ctx).Save(p).Error
}
