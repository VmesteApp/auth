package usecase

import (
	"context"

	"github.com/VmesteApp/auth-service/internal/entity"
)

// User Routes
type (
	User interface {
		CreateAccount(ctx context.Context, email, password string) error
		Login(ctx context.Context, email, password string) (*entity.User, string, error)
		VkLogin(ctx context.Context, userAccessToken string) (*entity.User, string, error)
	}
	UserRepo interface {
		SaveUser(ctx context.Context, email string, hassPash []byte) error
		User(ctx context.Context, email string) (*entity.User, error)
		SaveSocialUser(ctx context.Context, provider, providerID string) (*entity.User, error)
		SocialUser(ctx context.Context, provider, providerID string) (*entity.User, error)
	}
	VkWebApi interface {
		ValidateUserAccessToken(userAccessToken string) (*entity.VkTokenInfo, error)
	}
)

// Admin Routes
type (
	Admin interface {
		Admins(ctx context.Context) ([]entity.Admin, error)
		CreateAdmin(ctx context.Context, email, password string) error
		DeleteAdmin(ctx context.Context, userID uint64) error
	}
	AdminRepo interface {
		Admins(ctx context.Context) ([]entity.Admin, error)
		SaveAdmin(ctx context.Context, email string, passHash []byte) error
		DeleteAdmin(ctx context.Context, userID uint64) error
	}
)
