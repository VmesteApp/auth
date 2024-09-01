package usecase

import (
	"context"

	"github.com/VmesteApp/auth-service/internal/entity"
)

type AdminUseCase struct {
	repo AdminRepo
}

func NewAdminUseCase(repo AdminRepo) *AdminUseCase {
	return &AdminUseCase{
		repo: repo,
	}
}

func (u *AdminUseCase) Admins(ctx context.Context) ([]entity.User, error) {
	return []entity.User{}, nil
}

func (u *AdminUseCase) CreateAdmin(ctx context.Context, email, password string) error {
	return nil
}

func (u *AdminUseCase) DeleteAdmin(ctx context.Context, userID uint64) error {
	return nil
}
