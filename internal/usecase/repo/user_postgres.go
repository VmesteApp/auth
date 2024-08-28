package repo

import (
	"context"
	"fmt"

	"github.com/VmesteApp/auth-service/internal/entity"
	"github.com/VmesteApp/auth-service/pkg/postgres"
	"github.com/jackc/pgconn"
)

type UserRepository struct {
	*postgres.Postgres
}

func NewUserRepository(pg *postgres.Postgres) *UserRepository {
	return &UserRepository{pg}
}

func (u *UserRepository) SaveUser(ctx context.Context, email string, passHash []byte) error {
	sql, args, err := u.Builder.
		Insert("users").
		Columns("email", "pass_hash").
		Values(email, passHash).
		ToSql()
	if err != nil {
		return fmt.Errorf("can't to save user: %w", err)
	}

	_, err = u.Pool.Exec(ctx, sql, args...)
	if err != nil {
		if code, _ := err.(*pgconn.PgError); code.Code == "23505" {
			return entity.ErrUserExists
		}

		return fmt.Errorf("can't to save user: %w", err)
	}

	return nil
}

func (u *UserRepository) User(ctx context.Context, email string) (*entity.User, error) {
	return nil, nil
}
