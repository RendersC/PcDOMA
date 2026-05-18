package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/pcdoma/payment-service/internal/domain"
	"github.com/pcdoma/payment-service/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// ── Mock repository ─────────────────────────────────────────────

type mockPaymentRepo struct{ mock.Mock }

func (m *mockPaymentRepo) Create(ctx context.Context, p *domain.Payment) error {
	return m.Called(ctx, p).Error(0)
}
func (m *mockPaymentRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Payment, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Payment), args.Error(1)
}
func (m *mockPaymentRepo) GetByBookingID(ctx context.Context, bookingID uuid.UUID) (*domain.Payment, error) {
	args := m.Called(ctx, bookingID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Payment), args.Error(1)
}
func (m *mockPaymentRepo) ListByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*domain.Payment, int64, error) {
	args := m.Called(ctx, userID, limit, offset)
	return args.Get(0).([]*domain.Payment), args.Get(1).(int64), args.Error(2)
}
func (m *mockPaymentRepo) ListAll(ctx context.Context, limit, offset int) ([]*domain.Payment, int64, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]*domain.Payment), args.Get(1).(int64), args.Error(2)
}
func (m *mockPaymentRepo) Update(ctx context.Context, p *domain.Payment) error {
	return m.Called(ctx, p).Error(0)
}

// ── Tests ───────────────────────────────────────────────────────

func TestCreate_DefaultsCurrencyToKZT(t *testing.T) {
	repo := &mockPaymentRepo{}
	svc := service.NewPaymentService(repo)
	repo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Payment")).Return(nil)

	p, err := svc.Create(context.Background(), domain.InternalCreatePaymentRequest{
		BookingID: uuid.New(),
		UserID:    uuid.New(),
		Amount:    5000,
	})

	require.NoError(t, err)
	assert.Equal(t, "KZT", p.Currency)
	assert.Equal(t, domain.StatusPending, p.Status)
	assert.Equal(t, domain.MethodCard, p.Method)
	repo.AssertExpectations(t)
}

func TestCreate_RepoError(t *testing.T) {
	repo := &mockPaymentRepo{}
	svc := service.NewPaymentService(repo)
	repo.On("Create", mock.Anything, mock.Anything).Return(errors.New("db down"))

	_, err := svc.Create(context.Background(), domain.InternalCreatePaymentRequest{
		BookingID: uuid.New(), UserID: uuid.New(), Amount: 100,
	})

	assert.Error(t, err)
}

func TestProcess_ReachesTerminalState(t *testing.T) {
	repo := &mockPaymentRepo{}
	svc := service.NewPaymentService(repo)
	id := uuid.New()
	pending := &domain.Payment{ID: id, Status: domain.StatusPending}

	repo.On("GetByID", mock.Anything, id).Return(pending, nil)
	repo.On("Update", mock.Anything, mock.AnythingOfType("*domain.Payment")).Return(nil)

	p, err := svc.Process(context.Background(), id, uuid.New())

	require.NoError(t, err)
	// Mock gateway resolves to either completed or failed, never left pending/processing.
	assert.Contains(t, []domain.PaymentStatus{domain.StatusCompleted, domain.StatusFailed}, p.Status)
	assert.NotNil(t, p.ProcessedAt)
}

func TestProcess_AlreadyProcessed(t *testing.T) {
	repo := &mockPaymentRepo{}
	svc := service.NewPaymentService(repo)
	id := uuid.New()
	done := &domain.Payment{ID: id, Status: domain.StatusCompleted}
	repo.On("GetByID", mock.Anything, id).Return(done, nil)

	_, err := svc.Process(context.Background(), id, uuid.New())

	assert.ErrorIs(t, err, service.ErrAlreadyProcessed)
}

func TestProcess_NotFound(t *testing.T) {
	repo := &mockPaymentRepo{}
	svc := service.NewPaymentService(repo)
	id := uuid.New()
	repo.On("GetByID", mock.Anything, id).Return(nil, gorm.ErrRecordNotFound)

	_, err := svc.Process(context.Background(), id, uuid.New())

	assert.ErrorIs(t, err, service.ErrPaymentNotFound)
}

func TestRefund_SetsRefundedStatus(t *testing.T) {
	repo := &mockPaymentRepo{}
	svc := service.NewPaymentService(repo)
	id := uuid.New()
	completed := &domain.Payment{ID: id, Status: domain.StatusCompleted}
	repo.On("GetByID", mock.Anything, id).Return(completed, nil)
	repo.On("Update", mock.Anything, mock.AnythingOfType("*domain.Payment")).Return(nil)

	p, err := svc.Refund(context.Background(), id)

	require.NoError(t, err)
	assert.Equal(t, domain.StatusRefunded, p.Status)
	assert.NotNil(t, p.ProcessedAt)
}

func TestGetByID_NotFound(t *testing.T) {
	repo := &mockPaymentRepo{}
	svc := service.NewPaymentService(repo)
	id := uuid.New()
	repo.On("GetByID", mock.Anything, id).Return(nil, errors.New("missing"))

	_, err := svc.GetByID(context.Background(), id)

	assert.ErrorIs(t, err, service.ErrPaymentNotFound)
}
