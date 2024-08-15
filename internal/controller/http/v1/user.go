package v1

import (
	"github.com/VmesteApp/auth-service/internal/usecase"
	"github.com/VmesteApp/auth-service/pkg/logger"
	"github.com/gin-gonic/gin"
)

func newUserRoutes(handler *gin.RouterGroup, t usecase.User, l logger.Interface) {}
