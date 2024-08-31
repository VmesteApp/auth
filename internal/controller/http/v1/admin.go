package v1

import (
	"net/http"

	"github.com/VmesteApp/auth-service/pkg/logger"
	"github.com/gin-gonic/gin"
)

type adminRoutes struct {
	l logger.Interface
}

func newAdminRoutes(handler *gin.RouterGroup, l logger.Interface) {
	routes := &adminRoutes{
		l: l,
	}

	handler.GET("/", routes.doGetAllAdmins)
	handler.POST("/", routes.doCreateNewAdmin)
	handler.DELETE("/:id", routes.doDeleteAdmin)
}

func (a *adminRoutes) doGetAllAdmins(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, nil)
}

type doCreateNewAdminRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (a *adminRoutes) doCreateNewAdmin(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, nil)
}

func (a *adminRoutes) doDeleteAdmin(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, nil)
}
