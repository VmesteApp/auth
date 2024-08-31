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
	sql := `
	SELECT 
    u.id, u.email, u.pass_hash, u.role,
    ARRAY(SELECT id FROM social_logins WHERE user_id = u.id) AS social_login_ids,
    ARRAY(SELECT provider FROM social_logins WHERE user_id = u.id) AS providers,
    ARRAY(SELECT provider_id FROM social_logins WHERE user_id = u.id) AS provider_ids
		FROM users u
		WHERE u.email = $1
		GROUP BY u.id, u.email, u.pass_hash, u.role;
	`

	rows, err := u.Pool.Query(ctx, sql, email)
	if err != nil {
		return nil, fmt.Errorf("can't to find user by email: %w", err)
	}
	defer rows.Close()

	var user entity.User
	if rows.Next() {
		var ids []*uint64
		var providers, providerIds []*string

		err := rows.Scan(&user.ID, &user.Email, &user.PassHash, &user.Role, &ids, &providers, &providerIds)
		if err != nil {
			return nil, fmt.Errorf("can't to scan user: %w", err)
		}

		for i := 0; i < len(ids); i++ {
			el := entity.SocialLogin{ID: *ids[0], UserID: user.ID, ProviderID: *providerIds[i], Provider: *providers[i]}
			user.SocialLogins = append(user.SocialLogins, &el)
		}

		return &user, nil
	}

	return nil, entity.ErrUserNotFound
}
