package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/pcdoma/catalog-service/internal/domain"
	"github.com/pcdoma/catalog-service/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// ── Mock repositories ───────────────────────────────────────────

type mockPCRepo struct{ mock.Mock }

func (m *mockPCRepo) Create(ctx context.Context, pc *domain.PC) error {
	return m.Called(ctx, pc).Error(0)
}
func (m *mockPCRepo) GetByID(ctx context.Context, id string) (*domain.PC, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.PC), args.Error(1)
}
func (m *mockPCRepo) List(ctx context.Context, filter domain.PCFilter) ([]*domain.PC, int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]*domain.PC), args.Get(1).(int64), args.Error(2)
}
func (m *mockPCRepo) Update(ctx context.Context, id string, update *domain.UpdatePCRequest) (*domain.PC, error) {
	args := m.Called(ctx, id, update)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.PC), args.Error(1)
}
func (m *mockPCRepo) UpdateStatus(ctx context.Context, id string, status domain.PCStatus) error {
	return m.Called(ctx, id, status).Error(0)
}
func (m *mockPCRepo) Delete(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}
func (m *mockPCRepo) UpdateRating(ctx context.Context, id string, avg float64, count int) error {
	return m.Called(ctx, id, avg, count).Error(0)
}
func (m *mockPCRepo) GetLocations(ctx context.Context) ([]string, error) {
	args := m.Called(ctx)
	return args.Get(0).([]string), args.Error(1)
}

type mockReviewRepo struct{ mock.Mock }

func (m *mockReviewRepo) Create(ctx context.Context, review *domain.Review) error {
	return m.Called(ctx, review).Error(0)
}
func (m *mockReviewRepo) ListByPCID(ctx context.Context, pcID string, limit, offset int) ([]*domain.Review, int64, error) {
	args := m.Called(ctx, pcID, limit, offset)
	return args.Get(0).([]*domain.Review), args.Get(1).(int64), args.Error(2)
}
func (m *mockReviewRepo) GetRatingStats(ctx context.Context, pcID string) (float64, int, error) {
	args := m.Called(ctx, pcID)
	return args.Get(0).(float64), args.Get(1).(int), args.Error(2)
}
func (m *mockReviewRepo) Delete(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

// newService wires the service with nil redis (caching paths are nil-guarded).
func newService(pc *mockPCRepo, rev *mockReviewRepo) *service.PCService {
	return service.NewPCService(pc, rev, nil, 0)
}

// ── Tests ───────────────────────────────────────────────────────

func TestCreate_PersistsPC(t *testing.T) {
	pcRepo := &mockPCRepo{}
	svc := newService(pcRepo, &mockReviewRepo{})
	pcRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.PC")).Return(nil)

	pc, err := svc.Create(context.Background(), domain.CreatePCRequest{
		Name:         "Titan RGB",
		Specs:        domain.PCSpecs{CPU: "i9", GPU: "RTX 5080"},
		PricePerHour: 1500,
		PricePerDay:  9000,
	})

	require.NoError(t, err)
	assert.Equal(t, "Titan RGB", pc.Name)
	assert.Equal(t, "RTX 5080", pc.Specs.GPU)
	pcRepo.AssertExpectations(t)
}

func TestGetByID_NotFound(t *testing.T) {
	pcRepo := &mockPCRepo{}
	svc := newService(pcRepo, &mockReviewRepo{})
	pcRepo.On("GetByID", mock.Anything, "missing").Return(nil, nil)

	_, err := svc.GetByID(context.Background(), "missing")

	assert.Error(t, err)
}

func TestList_AppliesDefaultLimit(t *testing.T) {
	pcRepo := &mockPCRepo{}
	svc := newService(pcRepo, &mockReviewRepo{})
	// Expect the service to fill Limit=20 before hitting the repo.
	pcRepo.On("List", mock.Anything, mock.MatchedBy(func(f domain.PCFilter) bool {
		return f.Limit == 20
	})).Return([]*domain.PC{{Name: "PC1"}}, int64(1), nil)

	pcs, total, err := svc.List(context.Background(), domain.PCFilter{})

	require.NoError(t, err)
	assert.Len(t, pcs, 1)
	assert.Equal(t, int64(1), total)
	pcRepo.AssertExpectations(t)
}

func TestGetAvailability_AvailablePC(t *testing.T) {
	pcRepo := &mockPCRepo{}
	svc := newService(pcRepo, &mockReviewRepo{})
	pcRepo.On("GetByID", mock.Anything, "pc-1").Return(&domain.PC{
		Status:       domain.StatusAvailable,
		PricePerHour: 1500,
		PricePerDay:  9000,
	}, nil)

	avail, err := svc.GetAvailability(context.Background(), "pc-1")

	require.NoError(t, err)
	assert.True(t, avail.Available)
	assert.Equal(t, 1500.0, avail.PricePerHour)
}

func TestGetAvailability_BookedPC(t *testing.T) {
	pcRepo := &mockPCRepo{}
	svc := newService(pcRepo, &mockReviewRepo{})
	pcRepo.On("GetByID", mock.Anything, "pc-2").Return(&domain.PC{
		Status: domain.StatusBooked,
	}, nil)

	avail, err := svc.GetAvailability(context.Background(), "pc-2")

	require.NoError(t, err)
	assert.False(t, avail.Available)
}

func TestCreateReview_RecalculatesRating(t *testing.T) {
	pcRepo := &mockPCRepo{}
	reviewRepo := &mockReviewRepo{}
	svc := newService(pcRepo, reviewRepo)

	reviewRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Review")).Return(nil)
	reviewRepo.On("GetRatingStats", mock.Anything, "pc-1").Return(4.5, 2, nil)
	pcRepo.On("UpdateRating", mock.Anything, "pc-1", 4.5, 2).Return(nil)

	review, err := svc.CreateReview(context.Background(), "pc-1", "u-1", "Alice", domain.CreateReviewRequest{
		Rating:  5,
		Comment: "Great rig",
	})

	require.NoError(t, err)
	assert.Equal(t, 5, review.Rating)
	pcRepo.AssertExpectations(t)
	reviewRepo.AssertExpectations(t)
}

func TestDelete_PropagatesRepoError(t *testing.T) {
	pcRepo := &mockPCRepo{}
	svc := newService(pcRepo, &mockReviewRepo{})
	pcRepo.On("Delete", mock.Anything, "pc-9").Return(errors.New("delete failed"))

	err := svc.Delete(context.Background(), "pc-9")

	assert.Error(t, err)
}
