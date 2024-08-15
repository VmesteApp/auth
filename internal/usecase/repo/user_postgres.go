package repo

import "github.com/VmesteApp/auth-service/pkg/postgres"

type UserRepository struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) *UserRepository {
	return &UserRepository{pg}
}
