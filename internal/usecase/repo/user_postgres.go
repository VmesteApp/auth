package repo

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strconv"

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
	return u.doSaveUser(ctx, email, passHash, entity.UserRole)
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

func (u *UserRepository) Admins(ctx context.Context) ([]entity.Admin, error) {
	sql := `SELECT id, email FROM users WHERE role = $1`

	rows, err := u.Pool.Query(ctx, sql, entity.AdminRole)
	if err != nil {
		return nil, fmt.Errorf("can't find admins: %w", err)
	}
	defer rows.Close()

	var admins = make([]entity.Admin, 0)

	for rows.Next() {
		var admin entity.Admin

		err = rows.Scan(&admin.UserID, &admin.Email)
		if err != nil {
			return nil, fmt.Errorf("can't scan admin: %w", err)
		}

		admins = append(admins, admin)
	}

	return admins, nil
}

func (u *UserRepository) DeleteAdmin(ctx context.Context, userID uint64) error {
	sql := `DELETE FROM users WHERE id = $1`

	_, err := u.Pool.Exec(ctx, sql, userID)
	if err != nil {
		return fmt.Errorf("can't delete user: %w", err)
	}

	return nil
}

func (u *UserRepository) SaveAdmin(ctx context.Context, email string, passHash []byte) error {
	return u.doSaveUser(ctx, email, passHash, entity.AdminRole)
}

func (u *UserRepository) doSaveUser(ctx context.Context, email string, passHash []byte, role entity.Role) error {
	sql := `INSERT INTO users (email, pass_hash, role) VALUES ($1, $2, $3)`

	_, err := u.Pool.Exec(ctx, sql, email, passHash, role)
	if err != nil {
		if code, _ := err.(*pgconn.PgError); code.Code == "23505" {
			return entity.ErrUserExists
		}

		return fmt.Errorf("can't to save user: %w", err)
	}

	return nil
}

func (u *UserRepository) VkProfile(ctx context.Context, userID uint64) (entity.VkProfile, error) {
	sql := `SELECT provider_id FROM social_logins WHERE provider = 'vk' AND user_id = $1`

	var parsedVkID string

	err := u.Pool.QueryRow(ctx, sql, userID).Scan(&parsedVkID)
	fmt.Println(reflect.TypeOf(err))

	if err != nil {
		if err.Error() == "no rows in result set" {
			return entity.VkProfile{}, entity.ErrUserNotFound
		}

		return entity.VkProfile{}, fmt.Errorf("can't to get social logins: %w", err)
	}

	vkID, err := strconv.Atoi(parsedVkID)
	if err != nil {
		return entity.VkProfile{}, fmt.Errorf("can't parse provider id: %w", err)
	}

	return entity.VkProfile{
		UserID: userID,
		VkID:   vkID,
	}, nil
}
