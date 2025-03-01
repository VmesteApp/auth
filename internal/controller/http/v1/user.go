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
	handler.POST("/login/vk", r.doVkLoginByLaunchParams)
	handler.POST("/login/vk/access-token", r.doVkLoginByAccessToken)
}

type doRegisterNewUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// @Summary     Create account
// @Description Create account by email and password
// @ID          register
// @Tags  	    login
// @Param 			request body doRegisterNewUserRequest true "query params"
// @Accept      json
// @Success     200
// @Failure     400
// @Failure     401
// @Failure     409
// @Failure     500
// @Produce     json
// @Router      /register [post]
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

// @Summary     Login by email
// @Description Login by email for admin and superadmin users
// @ID          login
// @Tags  	    login
// @Param 			request body doLoginRequest true "query params"
// @Accept      json
// @Success     200  {object}   doLoginResponse
// @Failure     400
// @Failure     401
// @Failure     409
// @Failure     500
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

type doLoginByVkAccessTokenRequest struct {
	VkAccessToken string `json:"vkAccessToken" binding:"required"`
}

// @Summary     Login by VK
// @Description Login by VK for users
// @ID          login-vk-access-token
// @Tags  	    login
// @Param 			request body doLoginByVkAccessTokenRequest true "query params"
// @Accept      json
// @Success     200  {object}  doLoginResponse
// @Failure     400
// @Failure     401
// @Failure     500
// @Produce     json
// @Router      /login/vk/access-token [post]
func (r *userRoutes) doVkLoginByAccessToken(ctx *gin.Context) {
	var request doLoginByVkAccessTokenRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		r.l.Error(err, "http - v1 - doVkLoginByAccessToken")
		errorResponse(ctx, http.StatusBadRequest, "invalid request body")

		return
	}

	user, token, err := r.u.VkLoginByAccessToken(ctx.Request.Context(), request.VkAccessToken)
	if errors.Is(err, entity.ErrBadVkToken) {
		errorResponse(ctx, http.StatusBadRequest, "wrong access_token")
		return
	}
	if errors.Is(err, entity.ErrVkTokenExpired) {
		errorResponse(ctx, http.StatusUnauthorized, "access_token is expired")
		return
	}
	if err != nil {
		r.l.Error(err, "http - v1 - doVkLoginByAccessToken")
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

type doVkLoginByLaunchParamsRequest struct {
	VkLaunchParams string `json:"vkLaunchParams" binding:"required"`
}

// @Summary     Login by VK
// @Description Login by VK for users
// @ID          login-vk
// @Tags  	    login
// @Param 			request body doVkLoginByLaunchParamsRequest true "query params"
// @Accept      json
// @Success     200  {object}  doLoginResponse
// @Failure     400
// @Failure     401
// @Failure     500
// @Produce     json
// @Router      /login/vk [post]
func (r *userRoutes) doVkLoginByLaunchParams(ctx *gin.Context) {
	var request doVkLoginByLaunchParamsRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		r.l.Error(err, "http - v1 - doVkLoginByLaunchParams")
		errorResponse(ctx, http.StatusBadRequest, "invalid request body")

		return
	}

	user, token, err := r.u.VkLogin(ctx.Request.Context(), request.VkLaunchParams)
	if errors.Is(err, entity.ErrBadVkLaunchParams) {
		errorResponse(ctx, http.StatusBadRequest, "wrong launch params")
		return
	}
	if err != nil {
		r.l.Error(err, "http - v1 - doVkLoginByAccessToken")
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
