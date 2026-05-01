package domain

import (
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	RoleUser   Role = "user"
	RoleWorker Role = "worker"
	RoleAdmin  Role = "admin"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name         string    `gorm:"not null" json:"name"`
	Email        string    `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string    `gorm:"not null" json:"-"`
	Role         Role      `gorm:"type:varchar(10);default:'user'" json:"role"`
	Phone        string    `json:"phone,omitempty"`
	LocationID   string    `json:"location_id,omitempty"` // for workers: which location they serve
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type RefreshToken struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index"`
	Token     string    `gorm:"uniqueIndex;not null"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time
}

type RegisterRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=100"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Phone    string `json:"phone"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

type UpdateUserRequest struct {
	Name       string `json:"name" binding:"omitempty,min=2,max=100"`
	Phone      string `json:"phone"`
	LocationID string `json:"location_id"`
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

type ValidateTokenResponse struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	Name   string `json:"name"`
	Valid  bool   `json:"valid"`
}
