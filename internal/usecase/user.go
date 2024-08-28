package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/VmesteApp/auth-service/internal/entity"
	"golang.org/x/crypto/bcrypt"
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
func (u *UserUseCase) CreateAccount(ctx context.Context, email, password string) error {
	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("can't generate password hash: %w", err)
	}

	err = u.repo.SaveUser(ctx, email, passHash)
	if err != nil {
		if errors.Is(err, entity.ErrUserExists) {
			return err
		}

		return fmt.Errorf("can't save user: %w", err)
	}

	return nil
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
