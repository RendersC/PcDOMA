package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/pcdoma/auth-service/internal/domain"
	"github.com/pcdoma/auth-service/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// ── Mock repositories ──────────────────────────────────────────

type mockUserRepo struct{ mock.Mock }

func (m *mockUserRepo) Create(ctx context.Context, user *domain.User) error {
	return m.Called(ctx, user).Error(0)
}
func (m *mockUserRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}
func (m *mockUserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}
func (m *mockUserRepo) Update(ctx context.Context, user *domain.User) error {
	return m.Called(ctx, user).Error(0)
}
func (m *mockUserRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}
func (m *mockUserRepo) List(ctx context.Context, limit, offset int) ([]*domain.User, int64, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]*domain.User), args.Get(1).(int64), args.Error(2)
}

type mockTokenRepo struct{ mock.Mock }

func (m *mockTokenRepo) Create(ctx context.Context, token *domain.RefreshToken) error {
	return m.Called(ctx, token).Error(0)
}
func (m *mockTokenRepo) GetByToken(ctx context.Context, token string) (*domain.RefreshToken, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.RefreshToken), args.Error(1)
}
func (m *mockTokenRepo) DeleteByToken(ctx context.Context, token string) error {
	return m.Called(ctx, token).Error(0)
}
func (m *mockTokenRepo) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	return m.Called(ctx, userID).Error(0)
}

// ── Helpers ─────────────────────────────────────────────────────

func newAuthService(userRepo *mockUserRepo, tokenRepo *mockTokenRepo) service.AuthService {
	return service.NewAuthService(userRepo, tokenRepo, "test-secret", 15*time.Minute, 168*time.Hour)
}

func hashPassword(t *testing.T, password string) string {
	t.Helper()
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatal(err)
	}
	return string(hash)
}

// ── Tests ───────────────────────────────────────────────────────

func TestRegister_Success(t *testing.T) {
	userRepo := &mockUserRepo{}
	tokenRepo := &mockTokenRepo{}
	svc := newAuthService(userRepo, tokenRepo)

	userRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(nil, nil)
	userRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil)
	tokenRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.RefreshToken")).Return(nil)

	user, tokens, err := svc.Register(context.Background(), domain.RegisterRequest{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	})

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.NotNil(t, tokens)
	assert.NotEmpty(t, tokens.AccessToken)
	assert.NotEmpty(t, tokens.RefreshToken)
	assert.Equal(t, domain.RoleUser, user.Role)
	userRepo.AssertExpectations(t)
}

func TestRegister_EmailTaken(t *testing.T) {
	userRepo := &mockUserRepo{}
	tokenRepo := &mockTokenRepo{}
	svc := newAuthService(userRepo, tokenRepo)

	existing := &domain.User{ID: uuid.New(), Email: "test@example.com"}
	userRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(existing, nil)

	_, _, err := svc.Register(context.Background(), domain.RegisterRequest{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	})

	assert.ErrorIs(t, err, service.ErrEmailTaken)
}

func TestLogin_Success(t *testing.T) {
	userRepo := &mockUserRepo{}
	tokenRepo := &mockTokenRepo{}
	svc := newAuthService(userRepo, tokenRepo)

	user := &domain.User{
		ID:           uuid.New(),
		Email:        "test@example.com",
		PasswordHash: hashPassword(t, "password123"),
		Role:         domain.RoleUser,
	}
	userRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(user, nil)
	tokenRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.RefreshToken")).Return(nil)

	gotUser, tokens, err := svc.Login(context.Background(), domain.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	})

	assert.NoError(t, err)
	assert.Equal(t, user.ID, gotUser.ID)
	assert.NotEmpty(t, tokens.AccessToken)
}

func TestLogin_WrongPassword(t *testing.T) {
	userRepo := &mockUserRepo{}
	tokenRepo := &mockTokenRepo{}
	svc := newAuthService(userRepo, tokenRepo)

	user := &domain.User{
		ID:           uuid.New(),
		Email:        "test@example.com",
		PasswordHash: hashPassword(t, "correct"),
	}
	userRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(user, nil)

	_, _, err := svc.Login(context.Background(), domain.LoginRequest{
		Email:    "test@example.com",
		Password: "wrong",
	})

	assert.ErrorIs(t, err, service.ErrInvalidCredentials)
}

func TestLogin_UserNotFound(t *testing.T) {
	userRepo := &mockUserRepo{}
	tokenRepo := &mockTokenRepo{}
	svc := newAuthService(userRepo, tokenRepo)

	userRepo.On("GetByEmail", mock.Anything, "unknown@example.com").Return(nil, nil)

	_, _, err := svc.Login(context.Background(), domain.LoginRequest{
		Email:    "unknown@example.com",
		Password: "password",
	})

	assert.ErrorIs(t, err, service.ErrInvalidCredentials)
}

func TestValidateToken_Valid(t *testing.T) {
	userRepo := &mockUserRepo{}
	tokenRepo := &mockTokenRepo{}
	svc := newAuthService(userRepo, tokenRepo)

	userRepo.On("GetByEmail", mock.Anything, "v@example.com").Return(nil, nil)
	userRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil)
	tokenRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.RefreshToken")).Return(nil)

	_, tokens, _ := svc.Register(context.Background(), domain.RegisterRequest{
		Name:     "V User",
		Email:    "v@example.com",
		Password: "password",
	})

	resp, err := svc.ValidateToken(context.Background(), tokens.AccessToken)
	assert.NoError(t, err)
	assert.True(t, resp.Valid)
	assert.Equal(t, string(domain.RoleUser), resp.Role)
}

func TestValidateToken_Invalid(t *testing.T) {
	userRepo := &mockUserRepo{}
	tokenRepo := &mockTokenRepo{}
	svc := newAuthService(userRepo, tokenRepo)

	resp, err := svc.ValidateToken(context.Background(), "invalid.token.here")
	assert.NoError(t, err)
	assert.False(t, resp.Valid)
}
