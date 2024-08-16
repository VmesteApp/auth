package v1

import (
	"github.com/gin-gonic/gin"

	"github.com/VmesteApp/auth-service/internal/usecase"
	"github.com/VmesteApp/auth-service/pkg/logger"
)

func newUserRoutes(handler *gin.RouterGroup, t usecase.User, l logger.Interface) {}
