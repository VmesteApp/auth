// Package v1 implements routing paths. Each services in own file.
package v1

import (
	"net/http"

	"github.com/VmesteApp/auth-service/internal/usecase"
	"github.com/VmesteApp/auth-service/pkg/logger"
	"github.com/gin-gonic/gin"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// NewRouter -.
// Swagger spec:
// @title       vmesteapp.ru
// @description auth service
// @version     1.0
// @host        localhost:8080
// @BasePath    /v1
func NewRouter(handler *gin.Engine, l logger.Interface, t usecase.User) {
	// Options
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	// K8s probe
	handler.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	// Prometheus metrics
	handler.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Routers
	h := handler.Group("/v1")

	{
		newUserRoutes(h, t, l)
	}
}
