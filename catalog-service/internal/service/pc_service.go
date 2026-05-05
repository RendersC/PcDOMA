package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pcdoma/catalog-service/internal/domain"
	"github.com/redis/go-redis/v9"
)

// PCRepo is the persistence contract the service needs for PCs.
// The concrete *mongo.PCRepository satisfies it; tests provide a mock.
type PCRepo interface {
	Create(ctx context.Context, pc *domain.PC) error
	GetByID(ctx context.Context, id string) (*domain.PC, error)
	List(ctx context.Context, filter domain.PCFilter) ([]*domain.PC, int64, error)
	Update(ctx context.Context, id string, update *domain.UpdatePCRequest) (*domain.PC, error)
	UpdateStatus(ctx context.Context, id string, status domain.PCStatus) error
	Delete(ctx context.Context, id string) error
	UpdateRating(ctx context.Context, id string, avg float64, count int) error
	GetLocations(ctx context.Context) ([]string, error)
}

// ReviewRepo is the persistence contract the service needs for reviews.
type ReviewRepo interface {
	Create(ctx context.Context, review *domain.Review) error
	ListByPCID(ctx context.Context, pcID string, limit, offset int) ([]*domain.Review, int64, error)
	GetRatingStats(ctx context.Context, pcID string) (avg float64, count int, err error)
	Delete(ctx context.Context, id string) error
}

type PCService struct {
	pcRepo     PCRepo
	reviewRepo ReviewRepo
	redis      *redis.Client
	cacheTTL   time.Duration
}

func NewPCService(pcRepo PCRepo, reviewRepo ReviewRepo, rdb *redis.Client, cacheTTL time.Duration) *PCService {
	return &PCService{
		pcRepo:     pcRepo,
		reviewRepo: reviewRepo,
		redis:      rdb,
		cacheTTL:   cacheTTL,
	}
}

func (s *PCService) Create(ctx context.Context, req domain.CreatePCRequest) (*domain.PC, error) {
	pc := &domain.PC{
		Name:         req.Name,
		Specs:        req.Specs,
		PricePerHour: req.PricePerHour,
		PricePerDay:  req.PricePerDay,
		Location:     req.Location,
		Images:       req.Images,
	}
	if err := s.pcRepo.Create(ctx, pc); err != nil {
		return nil, err
	}
	s.invalidateCache(ctx)
	return pc, nil
}

func (s *PCService) GetByID(ctx context.Context, id string) (*domain.PC, error) {
	pc, err := s.pcRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if pc == nil {
		return nil, fmt.Errorf("pc not found")
	}
	return pc, nil
}

func (s *PCService) List(ctx context.Context, filter domain.PCFilter) ([]*domain.PC, int64, error) {
	if filter.Limit == 0 {
		filter.Limit = 20
	}

	cacheKey := fmt.Sprintf("pcs:list:%v", filter)
	if s.redis != nil {
		cached, err := s.redis.Get(ctx, cacheKey).Result()
		if err == nil {
			var result struct {
				PCs   []*domain.PC `json:"pcs"`
				Total int64        `json:"total"`
			}
			if json.Unmarshal([]byte(cached), &result) == nil {
				return result.PCs, result.Total, nil
			}
		}
	}

	pcs, total, err := s.pcRepo.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	if s.redis != nil {
		data, _ := json.Marshal(map[string]interface{}{"pcs": pcs, "total": total})
		s.redis.Set(ctx, cacheKey, data, s.cacheTTL)
	}

	return pcs, total, nil
}

func (s *PCService) Update(ctx context.Context, id string, req domain.UpdatePCRequest) (*domain.PC, error) {
	pc, err := s.pcRepo.Update(ctx, id, &req)
	if err != nil {
		return nil, err
	}
	s.invalidateCache(ctx)
	return pc, nil
}

func (s *PCService) UpdateStatus(ctx context.Context, id string, status domain.PCStatus) error {
	err := s.pcRepo.UpdateStatus(ctx, id, status)
	if err == nil {
		s.invalidateCache(ctx)
	}
	return err
}

func (s *PCService) Delete(ctx context.Context, id string) error {
	err := s.pcRepo.Delete(ctx, id)
	if err == nil {
		s.invalidateCache(ctx)
	}
	return err
}

func (s *PCService) GetAvailability(ctx context.Context, id string) (*domain.AvailabilityResponse, error) {
	pc, err := s.pcRepo.GetByID(ctx, id)
	if err != nil || pc == nil {
		return nil, fmt.Errorf("pc not found")
	}
	return &domain.AvailabilityResponse{
		Available:    pc.Status == domain.StatusAvailable,
		PricePerHour: pc.PricePerHour,
		PricePerDay:  pc.PricePerDay,
	}, nil
}

func (s *PCService) GetLocations(ctx context.Context) ([]string, error) {
	return s.pcRepo.GetLocations(ctx)
}

func (s *PCService) invalidateCache(ctx context.Context) {
	if s.redis == nil {
		return
	}
	iter := s.redis.Scan(ctx, 0, "pcs:list:*", 0).Iterator()
	for iter.Next(ctx) {
		s.redis.Del(ctx, iter.Val())
	}
}

// ReviewService embedded

func (s *PCService) CreateReview(ctx context.Context, pcID, userID, userName string, req domain.CreateReviewRequest) (*domain.Review, error) {
	review := &domain.Review{
		PCID:     pcID,
		UserID:   userID,
		UserName: userName,
		Rating:   req.Rating,
		Comment:  req.Comment,
	}
	if err := s.reviewRepo.Create(ctx, review); err != nil {
		return nil, err
	}

	avg, count, err := s.reviewRepo.GetRatingStats(ctx, pcID)
	if err == nil {
		s.pcRepo.UpdateRating(ctx, pcID, avg, count)
	}

	return review, nil
}

func (s *PCService) ListReviews(ctx context.Context, pcID string, limit, offset int) ([]*domain.Review, int64, error) {
	if limit == 0 {
		limit = 20
	}
	return s.reviewRepo.ListByPCID(ctx, pcID, limit, offset)
}

func (s *PCService) DeleteReview(ctx context.Context, reviewID string) error {
	return s.reviewRepo.Delete(ctx, reviewID)
}
