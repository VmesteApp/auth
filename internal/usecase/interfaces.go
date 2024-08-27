package usecase

import (
	"context"

	"github.com/VmesteApp/auth-service/internal/entity"
)

type (
	User interface {
		CreateAccount(ctx context.Context)
		Login(ctx context.Context)
		VkLogin(ctx context.Context, userAccessToken string) (string, error)
	}
	UserRepo interface{}
	VkWebApi interface {
		ValidateUserAccessToken(userAccessToken string) (*entity.VkTokenInfo, error)
	}
)
