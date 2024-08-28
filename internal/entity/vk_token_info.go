package entity

import "errors"

type VkTokenInfo struct {
	Date    int `json:"date"`
	Expire  int `json:"expire"`
	Success int `json:"success"`
	UserId  int `json:"user_id"`
}

var (
	ErrVkTokenExpired = errors.New("vk access_token expired")
	ErrBadVkToken     = errors.New("vk access_token is bad")
)
