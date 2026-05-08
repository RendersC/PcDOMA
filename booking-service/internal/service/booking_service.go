package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/pcdoma/booking-service/internal/clients"
	"github.com/pcdoma/booking-service/internal/domain"
)

var (
	ErrPCNotAvailable    = errors.New("pc is not available for the selected time")
	ErrBookingNotFound   = errors.New("booking not found")
	ErrForbidden         = errors.New("forbidden")
	ErrInvalidTransition = errors.New("invalid status transition")
)

// Repository is the persistence contract for bookings.
type Repository interface {
	Create(ctx context.Context, b *domain.Booking) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Booking, error)
	ListByUser(ctx context.Context, userID uuid.UUID, status string, limit, offset int) ([]*domain.Booking, int64, error)
	ListByLocation(ctx context.Context, locationID, status string, limit, offset int) ([]*domain.Booking, int64, error)
	ListAll(ctx context.Context, limit, offset int) ([]*domain.Booking, int64, error)
	Update(ctx context.Context, b *domain.Booking) error
}

// CatalogClient is the synchronous dependency on catalog-service.
type CatalogClient interface {
	GetPCAvailability(ctx context.Context, pcID string) (*clients.PCAvailability, error)
	GetSetupDetails(ctx context.Context, setupID string) (*clients.SetupDetails, error)
	UpdatePCStatus(ctx context.Context, pcID, status string) error
}

// PaymentClient is the synchronous dependency on payment-service.
type PaymentClient interface {
	CreatePayment(ctx context.Context, bookingID, userID uuid.UUID, amount float64, currency string) (*clients.CreatePaymentResponse, error)
}

// EventPublisher publishes async booking events over Redis pub/sub.
type EventPublisher interface {
	PublishBookingConfirmed(ctx context.Context, event domain.BookingEvent)
	PublishBookingCancelled(ctx context.Context, event domain.BookingEvent)
}

type BookingService struct {
	repo          Repository
	catalogClient CatalogClient
	paymentClient PaymentClient
	publisher     EventPublisher
}

func NewBookingService(
	repo Repository,
	catalogClient CatalogClient,
	paymentClient PaymentClient,
	publisher EventPublisher,
) *BookingService {
	return &BookingService{
		repo:          repo,
		catalogClient: catalogClient,
		paymentClient: paymentClient,
		publisher:     publisher,
	}
}

func (s *BookingService) Create(ctx context.Context, userID uuid.UUID, req domain.CreateBookingRequest) (*domain.Booking, error) {
	var basePrice, extrasPrice, discount float64
	var locationID string

	if req.BookingType == domain.TypeCustom {
		if req.PCID == "" {
			return nil, errors.New("pc_id required for custom booking")
		}
		avail, err := s.catalogClient.GetPCAvailability(ctx, req.PCID)
		if err != nil || !avail.Available {
			return nil, ErrPCNotAvailable
		}
		hours := req.EndTime.Sub(req.StartTime).Hours()
		if req.RentalType == domain.RentalDaily {
			days := hours / 24
			basePrice = avail.PricePerDay * days
		} else {
			basePrice = avail.PricePerHour * hours
		}
		locationID = "default"
	} else if req.BookingType == domain.TypeSetup {
		if req.SetupID == "" {
			return nil, errors.New("setup_id required for setup booking")
		}
		setup, err := s.catalogClient.GetSetupDetails(ctx, req.SetupID)
		if err != nil {
			return nil, ErrPCNotAvailable
		}
		if setup.PCID != "" {
			avail, err := s.catalogClient.GetPCAvailability(ctx, setup.PCID)
			if err != nil || !avail.Available {
				return nil, ErrPCNotAvailable
			}
			req.PCID = setup.PCID
		}
		hours := req.EndTime.Sub(req.StartTime).Hours()
		if req.RentalType == domain.RentalDaily {
			days := hours / 24
			basePrice = setup.TotalPricePerDay * days
		} else {
			basePrice = setup.TotalPricePerHour * hours
		}
		req.PeripheralIDs = setup.PeripheralIDs
		locationID = "default"
	}

	totalPrice := basePrice + extrasPrice - discount
	if totalPrice < 0 {
		totalPrice = 0
	}

	booking := &domain.Booking{
		ID:          uuid.New(),
		UserID:      userID,
		BookingType: req.BookingType,
		PCID:        req.PCID,
		SetupID:     req.SetupID,
		LocationID:  locationID,
		RentalType:  req.RentalType,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		BasePrice:   basePrice,
		ExtrasPrice: extrasPrice,
		Discount:    discount,
		TotalPrice:  totalPrice,
		Currency:    "KZT",
		Status:      domain.StatusPendingPayment,
	}

	for _, pID := range req.PeripheralIDs {
		booking.Peripherals = append(booking.Peripherals, domain.BookingPeripheral{
			ID:           uuid.New(),
			BookingID:    booking.ID,
			PeripheralID: pID,
		})
	}

	if err := s.repo.Create(ctx, booking); err != nil {
		return nil, err
	}

	// Create payment record
	payment, err := s.paymentClient.CreatePayment(ctx, booking.ID, userID, totalPrice, "KZT")
	if err == nil {
		booking.PaymentID = &payment.ID
		booking.Status = domain.StatusPaid
		s.repo.Update(ctx, booking)

		// Move to worker queue
		booking.Status = domain.StatusPendingWorker
		s.repo.Update(ctx, booking)

		// Mark PC as booked in catalog
		if req.PCID != "" {
			s.catalogClient.UpdatePCStatus(ctx, req.PCID, "booked")
		}

		// Publish events
		s.publisher.PublishBookingConfirmed(ctx, domain.BookingEvent{
			BookingID:  booking.ID,
			UserID:     userID,
			LocationID: locationID,
			StartTime:  req.StartTime,
			EndTime:    req.EndTime,
			Total:      totalPrice,
		})
	}

	return booking, nil
}

func (s *BookingService) GetByID(ctx context.Context, id, userID uuid.UUID, role string) (*domain.Booking, error) {
	b, err := s.repo.GetByID(ctx, id)
	if err != nil || b == nil {
		return nil, ErrBookingNotFound
	}
	if role != "admin" && role != "worker" && b.UserID != userID {
		return nil, ErrForbidden
	}
	return b, nil
}

func (s *BookingService) ListByUser(ctx context.Context, userID uuid.UUID, status string, limit, offset int) ([]*domain.Booking, int64, error) {
	if limit == 0 {
		limit = 20
	}
	return s.repo.ListByUser(ctx, userID, status, limit, offset)
}

func (s *BookingService) Cancel(ctx context.Context, id, userID uuid.UUID, role string) error {
	b, err := s.repo.GetByID(ctx, id)
	if err != nil || b == nil {
		return ErrBookingNotFound
	}
	if role != "admin" && b.UserID != userID {
		return ErrForbidden
	}
	if b.Status == domain.StatusCompleted || b.Status == domain.StatusActive {
		return ErrInvalidTransition
	}

	b.Status = domain.StatusCancelled
	if err := s.repo.Update(ctx, b); err != nil {
		return err
	}

	if b.PCID != "" {
		s.catalogClient.UpdatePCStatus(ctx, b.PCID, "available")
	}

	s.publisher.PublishBookingCancelled(ctx, domain.BookingEvent{
		BookingID:  b.ID,
		UserID:     userID,
		LocationID: b.LocationID,
	})

	return nil
}

// Worker actions

func (s *BookingService) WorkerAccept(ctx context.Context, id, workerID uuid.UUID) (*domain.Booking, error) {
	b, err := s.repo.GetByID(ctx, id)
	if err != nil || b == nil {
		return nil, ErrBookingNotFound
	}
	if b.Status != domain.StatusPendingWorker {
		return nil, ErrInvalidTransition
	}
	now := time.Now()
	b.Status = domain.StatusAccepted
	b.WorkerID = &workerID
	b.AcceptedAt = &now
	s.repo.Update(ctx, b)
	return b, nil
}

func (s *BookingService) WorkerReject(ctx context.Context, id, workerID uuid.UUID, reason string) (*domain.Booking, error) {
	b, err := s.repo.GetByID(ctx, id)
	if err != nil || b == nil {
		return nil, ErrBookingNotFound
	}
	if b.Status != domain.StatusPendingWorker {
		return nil, ErrInvalidTransition
	}
	b.Status = domain.StatusRejected
	b.WorkerID = &workerID
	b.RejectionReason = reason
	s.repo.Update(ctx, b)

	if b.PCID != "" {
		s.catalogClient.UpdatePCStatus(ctx, b.PCID, "available")
	}
	return b, nil
}

func (s *BookingService) WorkerDeliver(ctx context.Context, id uuid.UUID) (*domain.Booking, error) {
	b, err := s.repo.GetByID(ctx, id)
	if err != nil || b == nil {
		return nil, ErrBookingNotFound
	}
	if b.Status != domain.StatusAccepted {
		return nil, ErrInvalidTransition
	}
	now := time.Now()
	b.Status = domain.StatusDelivering
	b.DeliveredAt = &now
	s.repo.Update(ctx, b)
	return b, nil
}

func (s *BookingService) WorkerComplete(ctx context.Context, id uuid.UUID) (*domain.Booking, error) {
	b, err := s.repo.GetByID(ctx, id)
	if err != nil || b == nil {
		return nil, ErrBookingNotFound
	}
	if b.Status != domain.StatusDelivering && b.Status != domain.StatusActive {
		return nil, ErrInvalidTransition
	}
	now := time.Now()
	b.Status = domain.StatusCompleted
	b.CompletedAt = &now
	s.repo.Update(ctx, b)

	if b.PCID != "" {
		s.catalogClient.UpdatePCStatus(ctx, b.PCID, "available")
	}
	return b, nil
}

func (s *BookingService) ListByLocation(ctx context.Context, locationID, status string, limit, offset int) ([]*domain.Booking, int64, error) {
	if limit == 0 {
		limit = 50
	}
	return s.repo.ListByLocation(ctx, locationID, status, limit, offset)
}

func (s *BookingService) ListAll(ctx context.Context, limit, offset int) ([]*domain.Booking, int64, error) {
	if limit == 0 {
		limit = 50
	}
	return s.repo.ListAll(ctx, limit, offset)
}
