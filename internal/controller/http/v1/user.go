package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/VmesteApp/auth-service/internal/usecase"
	"github.com/VmesteApp/auth-service/pkg/logger"
)

type VkAuthRequestBody struct {
	VkAccessToken string `json:"vkAccessToken"`
}

type userRoutes struct {
	u usecase.User
	l logger.Interface
}

func newUserRoutes(handler *gin.RouterGroup, u usecase.User, l logger.Interface) {
	r := &userRoutes{u, l}

	handler.POST("/login/vk", r.loginByVk)
}

type doLoginByVkRequest struct {
	VkAccessToken string `json:"vkAccessToken" binding:"required"`
}

func (r *userRoutes) loginByVk(ctx *gin.Context) {
	var request doLoginByVkRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		r.l.Error(err, "http - v1 - doTranslate")
		errorResponse(ctx, http.StatusBadRequest, "invalid request body")

		return
	}

	token, err := r.u.VkLogin(ctx.Request.Context(), request.VkAccessToken)
	// TODO: add handling of custom BadToken
	if err != nil {
		r.l.Error(err, "http - v1 - loginByVk")
		errorResponse(ctx, http.StatusInternalServerError, "auth service problems")

		return
	}

	ctx.JSON(http.StatusOK, map[string]string{
		"token": token,
	})
}
