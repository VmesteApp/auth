package entity

import "errors"

type User struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	PassHash []byte
	Role     Role `json:"role"`
}
type Role string

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
