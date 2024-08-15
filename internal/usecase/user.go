package usecase

import "context"

type UserUseCase struct {
	repo UserRepo
}

// New - make user usecase.
func New(repo UserRepo) *UserUseCase {
	return &UserUseCase{
		repo: repo,
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

