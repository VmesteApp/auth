package entity

import (
	"errors"
)

type User struct {
	ID           uint64         `json:"id"`
	Email        string         `json:"email"`
	Role         Role           `json:"role"`
	SocialLogins []*SocialLogin `json:"socialLogins,omitempty"`

	PassHash []byte
}

type Role string

type SocialLogin struct {
	ID         uint64    `json:"id"`
	UserID     uint64    `json:"userId"`
	ProviderID string    `json:"providerId"`
	Provider   string    `json:"provider"`
}

const (
	UserRole       Role = "user"
	AdminRole      Role = "admin"
	SuperAdminRole Role = "superadmin"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserExists         = errors.New("user exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)
