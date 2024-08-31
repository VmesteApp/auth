package app

import (
	"context"

	"github.com/VmesteApp/auth-service/config"
	"github.com/VmesteApp/auth-service/internal/entity"
	"github.com/VmesteApp/auth-service/pkg/postgres"
	"golang.org/x/crypto/bcrypt"
)

func InitSuperAdmin(pg *postgres.Postgres, cfg config.SuperAdminConfig) error {
	_, err := pg.Pool.Exec(context.Background(), "DELETE FROM users WHERE role = $1", entity.SuperAdminRole)
	if err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(cfg.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = pg.Pool.Exec(context.Background(), "INSERT INTO users (email, pass_hash, role) VALUES ($1, $2, $3)", cfg.Email, hashedPassword, entity.SuperAdminRole)
	if err != nil {
		return err
	}

	return nil
}
