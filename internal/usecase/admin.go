package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/VmesteApp/auth-service/internal/entity"
	"golang.org/x/crypto/bcrypt"
)

type AdminUseCase struct {
	repo AdminRepo
}

func NewAdminUseCase(repo AdminRepo) *AdminUseCase {
	return &AdminUseCase{
		repo: repo,
	}
}

func (u *AdminUseCase) Admins(ctx context.Context) ([]entity.Admin, error) {
	admins, err := u.repo.Admins(ctx)
	if err != nil {
		return []entity.Admin{}, fmt.Errorf("can't get all admins: %w", err)
	}

	return admins, nil
}

func (u *AdminUseCase) CreateAdmin(ctx context.Context, email, password string) error {
	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("can't generate password hash: %w", err)
	}

	err = u.repo.SaveAdmin(ctx, email, passHash)
	if errors.Is(err, entity.ErrUserExists) {
		return err
	}
	if err != nil {
		return fmt.Errorf("can't save admin: %w", err)
	}

	return nil
}

func (u *AdminUseCase) DeleteAdmin(ctx context.Context, userID uint64) error {
	err := u.repo.DeleteAdmin(ctx, userID)
	if err != nil {
		return fmt.Errorf("can't delete admin: %w", err)
	}

	return nil
}
