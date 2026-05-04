package service

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/pcdoma/auth-service/internal/domain"
	"github.com/pcdoma/auth-service/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrEmailTaken         = errors.New("email already in use")
	ErrInvalidToken       = errors.New("invalid or expired token")
	ErrUserNotFound       = errors.New("user not found")
	ErrForbidden          = errors.New("forbidden")
)

type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	Name   string `json:"name"`
	jwt.RegisteredClaims
}

type AuthService interface {
	Register(ctx context.Context, req domain.RegisterRequest) (*domain.User, *domain.TokenPair, error)
	Login(ctx context.Context, req domain.LoginRequest) (*domain.User, *domain.TokenPair, error)
	Refresh(ctx context.Context, refreshToken string) (*domain.TokenPair, error)
	Logout(ctx context.Context, refreshToken string) error
	ValidateToken(ctx context.Context, tokenStr string) (*domain.ValidateTokenResponse, error)
}

type authService struct {
	userRepo    repository.UserRepository
	tokenRepo   repository.RefreshTokenRepository
	jwtSecret   string
	accessTTL   time.Duration
	refreshTTL  time.Duration
}

func NewAuthService(
	userRepo repository.UserRepository,
	tokenRepo repository.RefreshTokenRepository,
	jwtSecret string,
	accessTTL, refreshTTL time.Duration,
) AuthService {
	return &authService{
		userRepo:   userRepo,
		tokenRepo:  tokenRepo,
		jwtSecret:  jwtSecret,
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}
}

func (s *authService) Register(ctx context.Context, req domain.RegisterRequest) (*domain.User, *domain.TokenPair, error) {
	existing, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, nil, err
	}
	if existing != nil {
		return nil, nil, ErrEmailTaken
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, nil, err
	}

	user := &domain.User{
		ID:           uuid.New(),
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hash),
		Phone:        req.Phone,
		Role:         domain.RoleUser,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, nil, err
	}

	tokens, err := s.generateTokenPair(ctx, user)
	if err != nil {
		return nil, nil, err
	}

	return user, tokens, nil
}

func (s *authService) Login(ctx context.Context, req domain.LoginRequest) (*domain.User, *domain.TokenPair, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, nil, err
	}
	if user == nil {
		return nil, nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, nil, ErrInvalidCredentials
	}

	tokens, err := s.generateTokenPair(ctx, user)
	if err != nil {
		return nil, nil, err
	}

	return user, tokens, nil
}

func (s *authService) Refresh(ctx context.Context, refreshToken string) (*domain.TokenPair, error) {
	rt, err := s.tokenRepo.GetByToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}
	if rt == nil || time.Now().After(rt.ExpiresAt) {
		return nil, ErrInvalidToken
	}

	user, err := s.userRepo.GetByID(ctx, rt.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	if err := s.tokenRepo.DeleteByToken(ctx, refreshToken); err != nil {
		return nil, err
	}

	return s.generateTokenPair(ctx, user)
}

func (s *authService) Logout(ctx context.Context, refreshToken string) error {
	return s.tokenRepo.DeleteByToken(ctx, refreshToken)
}

func (s *authService) ValidateToken(_ context.Context, tokenStr string) (*domain.ValidateTokenResponse, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})
	if err != nil || !token.Valid {
		return &domain.ValidateTokenResponse{Valid: false}, nil
	}

	return &domain.ValidateTokenResponse{
		UserID: claims.UserID,
		Role:   claims.Role,
		Name:   claims.Name,
		Valid:  true,
	}, nil
}

func (s *authService) generateTokenPair(ctx context.Context, user *domain.User) (*domain.TokenPair, error) {
	now := time.Now()

	accessClaims := &Claims{
		UserID: user.ID.String(),
		Role:   string(user.Role),
		Name:   user.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.accessTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
			Subject:   user.ID.String(),
		},
	}

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, err
	}

	rawRefresh := uuid.New().String()
	rt := &domain.RefreshToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     rawRefresh,
		ExpiresAt: now.Add(s.refreshTTL),
	}

	if err := s.tokenRepo.Create(ctx, rt); err != nil {
		return nil, err
	}

	return &domain.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: rawRefresh,
		ExpiresIn:    int64(s.accessTTL.Seconds()),
	}, nil
}
