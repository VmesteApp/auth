package entity

import "errors"

type User struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	PassHash string
}

var ErrUserExists = errors.New("user exists")
