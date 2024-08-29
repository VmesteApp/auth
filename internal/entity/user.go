package entity

import "errors"

type User struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	PassHash string
	Role     Role `json:"role"`
}
type Role string

const (
	UserRole Role = "user"
	AdminRole Role = "admin"
	SuperAdminRole Role = "superadmin"
)

var ErrUserExists = errors.New("user exists")
