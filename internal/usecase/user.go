package usecase

type UserUseCase struct {
	repo UserRepo
}

// New - make user usecase.
func New(repo UserRepo) *UserUseCase {
	return &UserUseCase{
		repo: repo,
	}
}
