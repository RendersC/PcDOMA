package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/pcdoma/auth-service/internal/domain"
	"github.com/pcdoma/auth-service/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	List(ctx context.Context, limit, offset int) ([]*domain.User, int64, error)
	Update(ctx context.Context, id uuid.UUID, req domain.UpdateUserRequest, callerID uuid.UUID, callerRole string) (*domain.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
	ChangePassword(ctx context.Context, id uuid.UUID, req domain.ChangePasswordRequest) error
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (s *userService) List(ctx context.Context, limit, offset int) ([]*domain.User, int64, error) {
	return s.userRepo.List(ctx, limit, offset)
}

func (s *userService) Update(ctx context.Context, id uuid.UUID, req domain.UpdateUserRequest, callerID uuid.UUID, callerRole string) (*domain.User, error) {
	if callerRole != string(domain.RoleAdmin) && callerID != id {
		return nil, ErrForbidden
	}

	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.LocationID != "" && callerRole == string(domain.RoleAdmin) {
		user.LocationID = req.LocationID
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) Delete(ctx context.Context, id uuid.UUID) error {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}
	return s.userRepo.Delete(ctx, id)
}

func (s *userService) ChangePassword(ctx context.Context, id uuid.UUID, req domain.ChangePasswordRequest) error {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.OldPassword)); err != nil {
		return ErrInvalidCredentials
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.PasswordHash = string(hash)
	return s.userRepo.Update(ctx, user)
}
