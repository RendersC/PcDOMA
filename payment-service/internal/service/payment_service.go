package service

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/pcdoma/payment-service/internal/domain"
	"gorm.io/gorm"
)

var (
	ErrPaymentNotFound = errors.New("payment not found")
	ErrAlreadyProcessed = errors.New("payment already processed")
)

type PaymentRepository interface {
	Create(ctx context.Context, p *domain.Payment) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Payment, error)
	GetByBookingID(ctx context.Context, bookingID uuid.UUID) (*domain.Payment, error)
	ListByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*domain.Payment, int64, error)
	ListAll(ctx context.Context, limit, offset int) ([]*domain.Payment, int64, error)
	Update(ctx context.Context, p *domain.Payment) error
}

type PaymentService struct {
	repo PaymentRepository
}

func NewPaymentService(repo PaymentRepository) *PaymentService {
	return &PaymentService{repo: repo}
}

func (s *PaymentService) Create(ctx context.Context, req domain.InternalCreatePaymentRequest) (*domain.Payment, error) {
	currency := req.Currency
	if currency == "" {
		currency = "KZT"
	}
	p := &domain.Payment{
		ID:        uuid.New(),
		BookingID: req.BookingID,
		UserID:    req.UserID,
		Amount:    req.Amount,
		Currency:  currency,
		Method:    domain.MethodCard,
		Status:    domain.StatusPending,
	}
	if err := s.repo.Create(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *PaymentService) Process(ctx context.Context, id, userID uuid.UUID) (*domain.Payment, error) {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPaymentNotFound
		}
		return nil, err
	}
	if p == nil {
		return nil, ErrPaymentNotFound
	}
	if p.Status != domain.StatusPending {
		return nil, ErrAlreadyProcessed
	}

	p.Status = domain.StatusProcessing
	s.repo.Update(ctx, p)

	// Mock payment: 95% success
	time.Sleep(100 * time.Millisecond)
	now := time.Now()
	p.ProcessedAt = &now
	if rand.Float64() < 0.95 {
		p.Status = domain.StatusCompleted
	} else {
		p.Status = domain.StatusFailed
	}

	if err := s.repo.Update(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *PaymentService) Refund(ctx context.Context, id uuid.UUID) (*domain.Payment, error) {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil || p == nil {
		return nil, ErrPaymentNotFound
	}
	p.Status = domain.StatusRefunded
	now := time.Now()
	p.ProcessedAt = &now
	if err := s.repo.Update(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *PaymentService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Payment, error) {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil || p == nil {
		return nil, ErrPaymentNotFound
	}
	return p, nil
}

func (s *PaymentService) ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*domain.Payment, int64, error) {
	return s.repo.ListByUserID(ctx, userID, limit, offset)
}

func (s *PaymentService) ListAll(ctx context.Context, limit, offset int) ([]*domain.Payment, int64, error) {
	return s.repo.ListAll(ctx, limit, offset)
}
