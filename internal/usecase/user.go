package usecase

import (
	"context"
	"fmt"
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
		return "", fmt.Errorf("can't validate user access_token: %w", err)
	}

	fmt.Println(tokenInfo)
	if tokenInfo.Success == 1 {
		// TODO: add jwt token generation and user registration
		return "jwt token", nil
	}

	// TODO: add custom error
	return "", fmt.Errorf("wrong token")
}
