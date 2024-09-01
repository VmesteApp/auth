package v1

import (
	"net/http"

	"github.com/VmesteApp/auth-service/internal/usecase"
	"github.com/VmesteApp/auth-service/pkg/logger"
	"github.com/gin-gonic/gin"
)

type adminRoutes struct {
	u usecase.Admin
	l logger.Interface
}

func newAdminRoutes(handler *gin.RouterGroup, u usecase.Admin, l logger.Interface) {
	routes := &adminRoutes{
		l: l,
		u: u,
	}

	handler.GET("/", routes.doGetAllAdmins)
	handler.POST("/", routes.doCreateNewAdmin)
	handler.DELETE("/:id", routes.doDeleteAdmin)
}

func (a *adminRoutes) doGetAllAdmins(ctx *gin.Context) {
	admins, err := a.u.Admins(ctx.Request.Context())
	if err != nil {
		a.l.Error(err, "http - v1 - doGetAllAdmins")
		errorResponse(ctx, http.StatusInternalServerError, "SSO service problems")

		return
	}

	ctx.JSON(http.StatusOK, admins)
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
