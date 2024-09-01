package v1

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/VmesteApp/auth-service/internal/entity"
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
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (a *adminRoutes) doCreateNewAdmin(ctx *gin.Context) {
	var request doCreateNewAdminRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		errorResponse(ctx, http.StatusBadRequest, "invalid request")

		return
	}

	err := a.u.CreateAdmin(ctx.Request.Context(), request.Email, request.Password)
	if errors.Is(err, entity.ErrUserExists) {
		errorResponse(ctx, http.StatusConflict, "email already used")

		return
	}
	if err != nil {
		a.l.Error(err, "http - v1 - doCreateNewAdmin")
		errorResponse(ctx, http.StatusInternalServerError, "SSO service problems")

		return
	}

	ctx.JSON(http.StatusOK, nil)
}

func (a *adminRoutes) doDeleteAdmin(ctx *gin.Context) {
	id := ctx.Param("id")

	userID, err := strconv.Atoi(id)
	if err != nil {
		errorResponse(ctx, http.StatusBadRequest, "invalid ID")

		return
	}

	err = a.u.DeleteAdmin(ctx.Request.Context(), uint64(userID))
	if err != nil {
		a.l.Error(err, "http - v1 - doDeleteAdmin")
		errorResponse(ctx, http.StatusInternalServerError, "SSO service problems")

		return
	}

	ctx.JSON(http.StatusOK, nil)
}
