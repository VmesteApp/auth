package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	"github.com/VmesteApp/auth-service/config"
	v1 "github.com/VmesteApp/auth-service/internal/controller/http/v1"
	"github.com/VmesteApp/auth-service/internal/usecase"
	"github.com/VmesteApp/auth-service/internal/usecase/repo"
	"github.com/VmesteApp/auth-service/internal/usecase/webapi"
	"github.com/VmesteApp/auth-service/pkg/httpserver"
	"github.com/VmesteApp/auth-service/pkg/logger"
	"github.com/VmesteApp/auth-service/pkg/postgres"
)

func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)
	l.Info("logger init")

	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))

	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()
	l.Info("connected to database")

	err = InitSuperAdmin(pg, cfg.SuperAdminConfig)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - app.InitSuperAdmin: %w", err))
	}
	l.Info("%s is current superadmin!", cfg.SuperAdminConfig.Email)

	userUseCase := usecase.New(repo.NewUserRepository(pg), webapi.New(cfg.AppId, cfg.ServiceKey), cfg.JwtConfig.Secret, cfg.JwtConfig.TTL)

	handler := gin.New()
	v1.NewRouter(handler, l, userUseCase, cfg)

	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
