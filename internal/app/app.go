package app

import (
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"

	"github.com/VmesteApp/auth-service/config"
	profileGRPC "github.com/VmesteApp/auth-service/internal/controller/grpc/profile"
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

	// DB
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))

	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()
	l.Info("connected to database")

	// Permissions
	err = InitSuperAdmin(pg, cfg.SuperAdminConfig)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - app.InitSuperAdmin: %w", err))
	}
	l.Info("%s is current superadmin!", cfg.SuperAdminConfig.Email)

	// Usecases
	userRepository := repo.NewUserRepository(pg)

	userUseCase := usecase.New(userRepository, webapi.New(cfg.AppId, cfg.ServiceKey), cfg.JwtConfig.Secret, cfg.JwtConfig.TTL)
	adminUseCase := usecase.NewAdminUseCase(userRepository)
	profileUseCase := usecase.NewProfileUseCase(userRepository)

	// HTTP
	handler := gin.New()
	v1.NewRouter(handler, l, userUseCase, adminUseCase, profileUseCase, cfg)

	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// gRPC
	gRPCServer := grpc.NewServer()
	profileGRPC.Register(gRPCServer, userRepository)

	gRPClistener, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPC.Port))
	if err != nil {
		l.Fatal("failed list gRPC: %s", err)
	}

	l.Info("gRPC is runnig: %s", slog.String("addr", gRPClistener.Addr().String()))

	if err := gRPCServer.Serve(gRPClistener); err != nil {
		l.Fatal("failed serve gRPC listener: %s", err)
	}

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
	gRPCServer.GracefulStop()
}
