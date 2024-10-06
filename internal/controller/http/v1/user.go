package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/VmesteApp/auth-service/internal/entity"
	"github.com/VmesteApp/auth-service/internal/usecase"
	"github.com/VmesteApp/auth-service/pkg/logger"
)

type userRoutes struct {
	u usecase.User
	l logger.Interface
}

func newUserRoutes(handler *gin.RouterGroup, u usecase.User, l logger.Interface) {
	r := &userRoutes{u, l}

	handler.POST("/register", r.doRegisterNewUser)
	handler.POST("/login", r.doLoginByEmail)
	handler.POST("/login/vk", r.doLoginByVk)
}

type doRegisterNewUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (r *userRoutes) doRegisterNewUser(ctx *gin.Context) {
	var request doRegisterNewUserRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		r.l.Error(err, "http - v1 - doRegisterNewUser")
		errorResponse(ctx, http.StatusBadRequest, "invalid request body")

		return
	}

	err := r.u.CreateAccount(ctx.Request.Context(), request.Email, request.Password)
	if errors.Is(err, entity.ErrUserExists) {
		errorResponse(ctx, http.StatusConflict, "user already exists")

		return
	}
	if err != nil {
		r.l.Error(err, "http - v1 - doRegisterNewUser")
		errorResponse(ctx, http.StatusInternalServerError, "auth service problems")

		return
	}

	ctx.JSON(http.StatusOK, nil)
}

type doLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type doLoginResponse struct {
	Token  string      `json:"token"`
	UserID uint64      `json:"userId"`
	Role   entity.Role `json:"role"`
}

// @Summary     Login By Email
// @Description Login By Email (for admin and superadmin users)
// @ID          loginByEmail
// @Tags  	    login
// @Accept      json
// @Produce     json
// @Router      /login [post]
func (r *userRoutes) doLoginByEmail(ctx *gin.Context) {
	var request doLoginRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		r.l.Error(err, "http - v1 - doLoginByEmail")
		errorResponse(ctx, http.StatusBadRequest, "invalid request body")

		return
	}

	user, token, err := r.u.Login(ctx.Request.Context(), request.Email, request.Password)
	if errors.Is(err, entity.ErrUserNotFound) {
		errorResponse(ctx, http.StatusConflict, "user not found")

		return
	}
	if errors.Is(err, entity.ErrInvalidCredentials) {
		errorResponse(ctx, http.StatusUnauthorized, "wrong credentials")

		return
	}
	if err != nil {
		r.l.Error(err, "http - v1 - doLoginByEmail")
		errorResponse(ctx, http.StatusInternalServerError, "auth service problems")

		return
	}

	res := doLoginResponse{
		Token:  token,
		UserID: user.ID,
		Role:   user.Role,
	}

	ctx.JSON(http.StatusOK, res)
}

type doLoginByVkRequest struct {
	VkAccessToken string `json:"vkAccessToken" binding:"required"`
}

func (r *userRoutes) doLoginByVk(ctx *gin.Context) {
	var request doLoginByVkRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		r.l.Error(err, "http - v1 - doLoginByVk")
		errorResponse(ctx, http.StatusBadRequest, "invalid request body")

		return
	}

	user, token, err := r.u.VkLogin(ctx.Request.Context(), request.VkAccessToken)
	if errors.Is(err, entity.ErrBadVkToken) {
		errorResponse(ctx, http.StatusBadRequest, "wrong access_token")
		return
	}
	if errors.Is(err, entity.ErrVkTokenExpired) {
		errorResponse(ctx, http.StatusUnauthorized, "access_token is expired")
		return
	}
	if err != nil {
		r.l.Error(err, "http - v1 - loginByVk")
		errorResponse(ctx, http.StatusInternalServerError, "auth service problems")

		return
	}

	res := doLoginResponse{
		Token:  token,
		UserID: user.ID,
		Role:   user.Role,
	}

	ctx.JSON(http.StatusOK, res)
}
