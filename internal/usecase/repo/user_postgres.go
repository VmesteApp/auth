package repo

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgconn"

	"github.com/VmesteApp/auth-service/internal/entity"
	"github.com/VmesteApp/auth-service/pkg/postgres"
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
	sql, args, err := u.Builder.
		Select("id", "email", "pass_hash", "role").
		From("users").
		Where(squirrel.Eq{"email": email}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("can't to find user by email: %w", err)
	}

	rows, err := u.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("can't to find user by email: %w", err)
	}
	defer rows.Close()

	var user entity.User
	if rows.Next() {
		err := rows.Scan(&user.ID, &user.Email, &user.PassHash, &user.Role)
		if err != nil {
			return nil, fmt.Errorf("can't to scan user: %w", err)
		}
		return &user, nil
	}

	return nil, entity.ErrUserNotFound
}
