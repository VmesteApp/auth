package repo

import (
	"context"
	"database/sql"
	"fmt"

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

// TODO: are there need join of SocialLogin?
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

func (u *UserRepository) SaveSocialUser(ctx context.Context, provider, providerID string) (*entity.User, error) {
	tx, err := u.Postgres.Pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed start transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	var newUserId uint64
	err = tx.QueryRow(ctx, "INSERT INTO users (email, pass_hash) VALUES (NULL, NULL) RETURNING id").Scan(&newUserId)
	if err != nil {
		return nil, fmt.Errorf("failed to insert user: %w", err)
	}

	sql := `
		INSERT INTO social_logins 
			(user_id, provider, provider_id) 
			VALUES ($1, $2, $3)
	`

	fmt.Println(newUserId, provider, providerID)
	_, err = tx.Exec(ctx, sql, newUserId, provider, providerID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert social_logins: %w", err)
	}

	return &entity.User{
		ID:   newUserId,
		Role: entity.UserRole,
	}, nil
}

func (u *UserRepository) SocialUser(ctx context.Context, provider, providerId string) (*entity.User, error) {
	query := `
	SELECT 
		u.id, u.email, u.pass_hash, u.role 
		FROM users u 
		JOIN social_logins s 
			ON s.user_id = u.id 
		WHERE s.provider = $1 AND s.provider_id = $2;	
	`

	rows, err := u.Pool.Query(ctx, query, provider, providerId)
	if err != nil {
		return nil, fmt.Errorf("can't to find user by social login: %w", err)
	}
	defer rows.Close()

	var user entity.User

	if rows.Next() {
		var email sql.NullString
		var passHash sql.Null[[]byte]

		err := rows.Scan(&user.ID, &email, &passHash, &user.Role)
		if err != nil {
			return nil, fmt.Errorf("can't to scan user: %w", err)
		}

		if email.Valid {
			user.Email = email.String
		}
		if passHash.Valid {
			user.PassHash = passHash.V
		}

		return &user, nil
	}

	return nil, entity.ErrUserNotFound
}
