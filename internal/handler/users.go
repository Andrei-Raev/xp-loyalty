package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/Andrei-Raev/xp-loyalty/internal/model"
)

type UserService interface {
	Create(ctx context.Context, id, username, avatarURL, nickname string, role model.Role) error
	GetByUsername(ctx context.Context, username string) (model.User, error)
	Prizes(ctx context.Context, username string) ([]model.UserPrize, error)
}

type UserHandler struct {
	userService UserService
}

func NewUserHandler(userService UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

type getUserResponse struct {
	model.CredentialsSecure
	Nickname             string            `json:"nickname"`
	AvatarURL            string            `json:"avatar_url"`
	XPoints              int               `json:"XPoints"`
	RegistrationTime     time.Time         `json:"registration_time"`
	Prizes               []model.UserPrize `json:"prizes"`
	LastDailyCardsUpdate time.Time         `json:"-"`
}

// @Summary get user by username
// @Tags users
// @Param username path string true "username"
// @Router /api/users/{username} [get]
// @Security ApiKeyAuth
func (h UserHandler) Get(ctx *gin.Context) {
	username, err := ParsePath(ctx, "username")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, E(err))
	}

	user, err := h.userService.GetByUsername(ctx.Request.Context(), username)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, E(err))
		return
	}

	prizes, err := h.userService.Prizes(ctx.Request.Context(), username)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, E(err))
		return
	}

	resp := getUserResponse{
		CredentialsSecure:    user.CredentialsSecure,
		Nickname:             user.Nickname,
		AvatarURL:            user.AvatarURL,
		XPoints:              user.XPoints,
		RegistrationTime:     user.RegistrationTime,
		LastDailyCardsUpdate: user.LastDailyCardsUpdate,
		Prizes:               prizes,
	}

	ctx.JSON(http.StatusOK, resp)
}

// @Summary get user by token
// @Tags users
// @Router /api/users/profile [get]
// @Security ApiKeyAuth
func (h UserHandler) Profile(ctx *gin.Context) {
	c, ok := ctx.Get(model.CtxCredentialsKey)
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, M("user does not exist"))
		return
	}

	credentials, ok := c.(model.Credentials)
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, M("wrong token"))
		return
	}

	user, err := h.userService.GetByUsername(ctx.Request.Context(), credentials.Username)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, E(err))
		return
	}

	prizes, err := h.userService.Prizes(ctx.Request.Context(), credentials.Username)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, E(err))
		return
	}

	resp := getUserResponse{
		CredentialsSecure:    user.CredentialsSecure,
		Nickname:             user.Nickname,
		AvatarURL:            user.AvatarURL,
		XPoints:              user.XPoints,
		RegistrationTime:     user.RegistrationTime,
		LastDailyCardsUpdate: user.LastDailyCardsUpdate,
		Prizes:               prizes,
	}

	ctx.JSON(http.StatusOK, resp)
}
