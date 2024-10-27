// Package v1 implements routing paths. Each services in own file.
package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/VmesteApp/auth-service/docs"

	"github.com/VmesteApp/auth-service/config"
	"github.com/VmesteApp/auth-service/internal/entity"
	"github.com/VmesteApp/auth-service/internal/usecase"
	"github.com/VmesteApp/auth-service/pkg/logger"
	"github.com/VmesteApp/auth-service/pkg/middlewares"
)

func NewRouter(handler *gin.Engine, l logger.Interface, t usecase.User, a usecase.Admin, p usecase.Profile, cfg *config.Config) {
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	// API docs
	handler.GET("auth/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// K8s probe
	handler.GET("/auth/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	// Prometheus metrics
	handler.GET("/auth/metrics", gin.WrapH(promhttp.Handler()))

	// Routers
	{
		h := handler.Group("/auth")

		newUserRoutes(h, t, l)
	}

	{
		h := handler.Group(
			"/auth/admin",
			middlewares.AuthMiddleware(cfg.JwtConfig.Secret),
			middlewares.RoleMiddleware(string(entity.SuperAdminRole)),
		)

		newAdminRoutes(h, a, l)
	}

	{
		h := handler.Group("/auth/profile", middlewares.AuthMiddleware(cfg.JwtConfig.Secret))

		newProfileRoutes(h, p, l)
	}
}
