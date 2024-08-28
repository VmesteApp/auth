package usecase

import (
	"context"
)

type UserUseCase struct {
	repo UserRepo
	api  VkWebApi
}

// New - make user usecase.
func New(repo UserRepo, webapi VkWebApi) *UserUseCase {
	return &UserUseCase{
		repo: repo,
		api:  webapi,
	}
}

// CreateAccount implements User.
func (u *UserUseCase) CreateAccount(ctx context.Context) {
	panic("unimplemented")
}

// Login implements User.
func (u *UserUseCase) Login(ctx context.Context) {
	panic("unimplemented")
}

func (u *UserUseCase) VkLogin(ctx context.Context, userAccessToken string) (string, error) {
	tokenInfo, err := u.api.ValidateUserAccessToken(userAccessToken)
	if err != nil {
		return "", err
	}

	_ = tokenInfo.UserId

	return "jwt token", nil
}
