package entity

type VkTokenInfo struct {
	Date    int `json:"date"`
	Expire  int `json:"expire"`
	Success int `json:"success"`
	UserId  int `json:"user_id"`
}
