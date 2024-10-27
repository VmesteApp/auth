package usecase

import (
	"context"

	"github.com/VmesteApp/auth-service/internal/entity"
)

type ProfileUseCase struct {
	repo ProfileRepo
}

func NewProfileUseCase(repo Profile) *ProfileUseCase {
	return &ProfileUseCase{repo: repo}
}

func (u *ProfileUseCase) VkProfile(ctx context.Context, userID uint64) (entity.VkProfile, error) {
	profile, err := u.repo.VkProfile(ctx, userID)

	return profile, err
}
