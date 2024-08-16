package usecase

import "context"

type (
	User interface {
		CreateAccount(ctx context.Context)
		Login(ctx context.Context)
	}
	UserRepo interface{}
)
