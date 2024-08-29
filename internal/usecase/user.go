package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/VmesteApp/auth-service/internal/entity"
	"github.com/VmesteApp/auth-service/pkg/jwt"
)

type UserUseCase struct {
	repo        UserRepo
	api         VkWebApi
	tokenSecret string
	tokenTTL    time.Duration
}

// New - make user usecase.
func New(repo UserRepo, webapi VkWebApi, tokenSecret string, tokenTTL time.Duration) *UserUseCase {
	return &UserUseCase{
		repo:        repo,
		api:         webapi,
		tokenSecret: tokenSecret,
		tokenTTL:    tokenTTL,
	}
}

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

func (u *UserUseCase) Login(ctx context.Context, email, password string) (string, error) {
	user, err := u.repo.User(ctx, email)
	if err != nil {
		if errors.Is(err, entity.ErrUserNotFound) {
			return "", entity.ErrUserNotFound
		}

		return "", fmt.Errorf("can't get user by email: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		return "", entity.ErrInvalidCredentials
	}

	payload := map[string]any{
		"uid":  user.ID,
		"role": user.Role,
	}

	token, err := jwt.NewToken(payload, u.tokenSecret, u.tokenTTL)
	if err != nil {
		return "", fmt.Errorf("can't generate token: %w", err)
	}

	return token, nil
}

func (u *UserUseCase) VkLogin(ctx context.Context, userAccessToken string) (string, error) {
	tokenInfo, err := u.api.ValidateUserAccessToken(userAccessToken)
	if err != nil {
		return "", err
	}

	_ = tokenInfo.UserId

	return "jwt token", nil
}
