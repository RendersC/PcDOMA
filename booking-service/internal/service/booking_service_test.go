package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/pcdoma/booking-service/internal/clients"
	"github.com/pcdoma/booking-service/internal/domain"
	"github.com/pcdoma/booking-service/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// ── Mocks ───────────────────────────────────────────────────────

type mockRepo struct{ mock.Mock }

func (m *mockRepo) Create(ctx context.Context, b *domain.Booking) error {
	return m.Called(ctx, b).Error(0)
}
func (m *mockRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Booking, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Booking), args.Error(1)
}
func (m *mockRepo) ListByUser(ctx context.Context, userID uuid.UUID, status string, limit, offset int) ([]*domain.Booking, int64, error) {
	args := m.Called(ctx, userID, status, limit, offset)
	return args.Get(0).([]*domain.Booking), args.Get(1).(int64), args.Error(2)
}
func (m *mockRepo) ListByLocation(ctx context.Context, locationID, status string, limit, offset int) ([]*domain.Booking, int64, error) {
	args := m.Called(ctx, locationID, status, limit, offset)
	return args.Get(0).([]*domain.Booking), args.Get(1).(int64), args.Error(2)
}
func (m *mockRepo) ListAll(ctx context.Context, limit, offset int) ([]*domain.Booking, int64, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]*domain.Booking), args.Get(1).(int64), args.Error(2)
}
func (m *mockRepo) Update(ctx context.Context, b *domain.Booking) error {
	return m.Called(ctx, b).Error(0)
}

type mockCatalog struct{ mock.Mock }

func (m *mockCatalog) GetPCAvailability(ctx context.Context, pcID string) (*clients.PCAvailability, error) {
	args := m.Called(ctx, pcID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*clients.PCAvailability), args.Error(1)
}
func (m *mockCatalog) GetSetupDetails(ctx context.Context, setupID string) (*clients.SetupDetails, error) {
	args := m.Called(ctx, setupID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*clients.SetupDetails), args.Error(1)
}
func (m *mockCatalog) UpdatePCStatus(ctx context.Context, pcID, status string) error {
	return m.Called(ctx, pcID, status).Error(0)
}

type mockPayment struct{ mock.Mock }

func (m *mockPayment) CreatePayment(ctx context.Context, bookingID, userID uuid.UUID, amount float64, currency string) (*clients.CreatePaymentResponse, error) {
	args := m.Called(ctx, bookingID, userID, amount, currency)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*clients.CreatePaymentResponse), args.Error(1)
}

type mockPublisher struct{ mock.Mock }

func (m *mockPublisher) PublishBookingConfirmed(ctx context.Context, event domain.BookingEvent) {
	m.Called(ctx, event)
}
func (m *mockPublisher) PublishBookingCancelled(ctx context.Context, event domain.BookingEvent) {
	m.Called(ctx, event)
}

// ── Tests ───────────────────────────────────────────────────────

func TestCreate_CustomBooking_Success(t *testing.T) {
	repo := &mockRepo{}
	catalog := &mockCatalog{}
	payment := &mockPayment{}
	pub := &mockPublisher{}
	svc := service.NewBookingService(repo, catalog, payment, pub)

	start := time.Now()
	end := start.Add(2 * time.Hour)

	catalog.On("GetPCAvailability", mock.Anything, "pc-1").
		Return(&clients.PCAvailability{Available: true, PricePerHour: 1000, PricePerDay: 8000}, nil)
	repo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Booking")).Return(nil)
	payment.On("CreatePayment", mock.Anything, mock.Anything, mock.Anything, 2000.0, "KZT").
		Return(&clients.CreatePaymentResponse{ID: uuid.New(), Status: "pending"}, nil)
	repo.On("Update", mock.Anything, mock.AnythingOfType("*domain.Booking")).Return(nil)
	catalog.On("UpdatePCStatus", mock.Anything, "pc-1", "booked").Return(nil)
	pub.On("PublishBookingConfirmed", mock.Anything, mock.Anything).Return()

	b, err := svc.Create(context.Background(), uuid.New(), domain.CreateBookingRequest{
		BookingType: domain.TypeCustom,
		PCID:        "pc-1",
		StartTime:   start,
		EndTime:     end,
		RentalType:  domain.RentalHourly,
	})

	require.NoError(t, err)
	assert.Equal(t, 2000.0, b.BasePrice)
	assert.Equal(t, 2000.0, b.TotalPrice)
	assert.Equal(t, domain.StatusPendingWorker, b.Status)
	assert.NotNil(t, b.PaymentID)
	catalog.AssertExpectations(t)
	payment.AssertExpectations(t)
	pub.AssertExpectations(t)
}

func TestCreate_CustomBooking_PCNotAvailable(t *testing.T) {
	repo := &mockRepo{}
	catalog := &mockCatalog{}
	svc := service.NewBookingService(repo, catalog, &mockPayment{}, &mockPublisher{})

	catalog.On("GetPCAvailability", mock.Anything, "pc-x").
		Return(&clients.PCAvailability{Available: false}, nil)

	_, err := svc.Create(context.Background(), uuid.New(), domain.CreateBookingRequest{
		BookingType: domain.TypeCustom,
		PCID:        "pc-x",
		StartTime:   time.Now(),
		EndTime:     time.Now().Add(time.Hour),
		RentalType:  domain.RentalHourly,
	})

	assert.ErrorIs(t, err, service.ErrPCNotAvailable)
}

func TestCreate_CustomBooking_MissingPCID(t *testing.T) {
	svc := service.NewBookingService(&mockRepo{}, &mockCatalog{}, &mockPayment{}, &mockPublisher{})

	_, err := svc.Create(context.Background(), uuid.New(), domain.CreateBookingRequest{
		BookingType: domain.TypeCustom,
		StartTime:   time.Now(),
		EndTime:     time.Now().Add(time.Hour),
		RentalType:  domain.RentalHourly,
	})

	assert.Error(t, err)
}

func TestWorkerAccept_Success(t *testing.T) {
	repo := &mockRepo{}
	svc := service.NewBookingService(repo, &mockCatalog{}, &mockPayment{}, &mockPublisher{})
	id := uuid.New()
	workerID := uuid.New()

	repo.On("GetByID", mock.Anything, id).Return(&domain.Booking{ID: id, Status: domain.StatusPendingWorker}, nil)
	repo.On("Update", mock.Anything, mock.AnythingOfType("*domain.Booking")).Return(nil)

	b, err := svc.WorkerAccept(context.Background(), id, workerID)

	require.NoError(t, err)
	assert.Equal(t, domain.StatusAccepted, b.Status)
	require.NotNil(t, b.WorkerID)
	assert.Equal(t, workerID, *b.WorkerID)
	assert.NotNil(t, b.AcceptedAt)
}

func TestWorkerAccept_InvalidTransition(t *testing.T) {
	repo := &mockRepo{}
	svc := service.NewBookingService(repo, &mockCatalog{}, &mockPayment{}, &mockPublisher{})
	id := uuid.New()
	repo.On("GetByID", mock.Anything, id).Return(&domain.Booking{ID: id, Status: domain.StatusCompleted}, nil)

	_, err := svc.WorkerAccept(context.Background(), id, uuid.New())

	assert.ErrorIs(t, err, service.ErrInvalidTransition)
}

func TestWorkerComplete_FromDelivering(t *testing.T) {
	repo := &mockRepo{}
	catalog := &mockCatalog{}
	svc := service.NewBookingService(repo, catalog, &mockPayment{}, &mockPublisher{})
	id := uuid.New()

	repo.On("GetByID", mock.Anything, id).Return(&domain.Booking{ID: id, Status: domain.StatusDelivering, PCID: "pc-1"}, nil)
	repo.On("Update", mock.Anything, mock.AnythingOfType("*domain.Booking")).Return(nil)
	catalog.On("UpdatePCStatus", mock.Anything, "pc-1", "available").Return(nil)

	b, err := svc.WorkerComplete(context.Background(), id)

	require.NoError(t, err)
	assert.Equal(t, domain.StatusCompleted, b.Status)
	assert.NotNil(t, b.CompletedAt)
	catalog.AssertExpectations(t)
}

func TestCancel_Success(t *testing.T) {
	repo := &mockRepo{}
	catalog := &mockCatalog{}
	pub := &mockPublisher{}
	svc := service.NewBookingService(repo, catalog, &mockPayment{}, pub)
	id := uuid.New()
	userID := uuid.New()

	repo.On("GetByID", mock.Anything, id).Return(&domain.Booking{ID: id, UserID: userID, Status: domain.StatusPendingWorker, PCID: "pc-1"}, nil)
	repo.On("Update", mock.Anything, mock.AnythingOfType("*domain.Booking")).Return(nil)
	catalog.On("UpdatePCStatus", mock.Anything, "pc-1", "available").Return(nil)
	pub.On("PublishBookingCancelled", mock.Anything, mock.Anything).Return()

	err := svc.Cancel(context.Background(), id, userID, "user")

	require.NoError(t, err)
	pub.AssertExpectations(t)
}

func TestCancel_ForbiddenForOtherUser(t *testing.T) {
	repo := &mockRepo{}
	svc := service.NewBookingService(repo, &mockCatalog{}, &mockPayment{}, &mockPublisher{})
	id := uuid.New()
	repo.On("GetByID", mock.Anything, id).Return(&domain.Booking{ID: id, UserID: uuid.New(), Status: domain.StatusPendingWorker}, nil)

	err := svc.Cancel(context.Background(), id, uuid.New(), "user")

	assert.ErrorIs(t, err, service.ErrForbidden)
}

func TestGetByID_ForbiddenForOtherUser(t *testing.T) {
	repo := &mockRepo{}
	svc := service.NewBookingService(repo, &mockCatalog{}, &mockPayment{}, &mockPublisher{})
	id := uuid.New()
	repo.On("GetByID", mock.Anything, id).Return(&domain.Booking{ID: id, UserID: uuid.New()}, nil)

	_, err := svc.GetByID(context.Background(), id, uuid.New(), "user")

	assert.ErrorIs(t, err, service.ErrForbidden)
}
