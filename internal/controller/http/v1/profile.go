package v1

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/VmesteApp/auth-service/internal/entity"
	"github.com/VmesteApp/auth-service/internal/usecase"
	"github.com/VmesteApp/auth-service/pkg/logger"
)

type profileRoutes struct {
	u usecase.Profile
	l logger.Interface
}

func newProfileRoutes(handler *gin.RouterGroup, u usecase.Profile, l logger.Interface) {
	r := &profileRoutes{u, l}

	handler.GET("/:id/vk", r.doVkProfile)
}

// @Summary     Get VK profile
// @Description Get VK profile by user id
// @ID          vk-profile
// @Tags  	    profiles
// @Param       id   path      int  true  "User ID"
// @Accept      json
// @Success     200
// @Failure     400
// @Failure     401
// @Failure     409
// @Failure     500
// @Produce     json
// @Router      /profile/{id}/vk [get]
// @Security    BearerAuth
func (r *profileRoutes) doVkProfile(ctx *gin.Context) {
	id := ctx.Param("id")

	userID, err := strconv.Atoi(id)
	if err != nil {
		errorResponse(ctx, http.StatusBadRequest, "invalid ID")

		return
	}

	vkProfile, err := r.u.VkProfile(ctx.Request.Context(), uint64(userID))
	if errors.Is(err, entity.ErrUserNotFound) {
		errorResponse(ctx, http.StatusConflict, "user not found")

		return
	}
	if err != nil {
		r.l.Error(err, "http - v1 - doVkProfile")
		errorResponse(ctx, http.StatusInternalServerError, "auth service problems")

		return
	}

	ctx.JSON(http.StatusOK, vkProfile)
}
