package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Andrei-Raev/xp-loyalty/internal/model"
)

type AuthService interface {
	SignUp(ctx context.Context, credentials model.Credentials) (string, error)
	SignIn(ctx context.Context, username, password string) (string, error)
	CheckAccess(ctx context.Context, token string, role model.Role) (model.Credentials, error)
}

type AuthUserService interface {
	Create(ctx context.Context, id, username, avatarURL, nickname string, role model.Role) error
}

type AuthAdminService interface {
	Create(ctx context.Context, admin model.Admin) error
}

type AuthHandler struct {
	authService  AuthService
	userService  AuthUserService
	adminService AuthAdminService
}

func NewAuthHandler(authService AuthService, userService AuthUserService, adminService AuthAdminService) *AuthHandler {
	return &AuthHandler{
		authService:  authService,
		userService:  userService,
		adminService: adminService,
	}
}

type signUpAdminInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// @Summary sign up admin
// @Tags auth
// @Param input body signUpAdminInput true "sign up info"
// @Router /api/auth/sign-up-admin [post]
// @Security ApiKeyAuth
func (h AuthHandler) SignUpAdmin(ctx *gin.Context) {
	inp := new(signUpAdminInput)
	if err := ctx.BindJSON(inp); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, E(err))
		return
	}

	credentials := model.Credentials{
		CredentialsSecure: model.CredentialsSecure{
			Role:     model.RoleAdmin,
			Username: inp.Username,
		},
		Password: inp.Password,
	}
	id, err := h.authService.SignUp(ctx.Request.Context(), credentials)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, E(err))
		return
	}

	admin := model.Admin{
		CredentialsSecure: model.CredentialsSecure{
			ID:       id,
			Username: inp.Username,
			Role:     model.RoleAdmin,
		},
	}
	if err := h.adminService.Create(ctx.Request.Context(), admin); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, E(err))
	}

	ctx.JSON(http.StatusOK, M("ok"))
}

type signUpUserInput struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Nickname  string `json:"nickname"`
	AvatarURL string `json:"avatar_url"`
}

// @Summary sign up user
// @Tags auth
// @Param input body signUpUserInput true "sign up info"
// @Router /api/auth/sign-up-user [post]
// @Security ApiKeyAuth
func (h AuthHandler) SignUpUser(ctx *gin.Context) {
	inp := new(signUpUserInput)
	if err := ctx.BindJSON(inp); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, E(err))
		return
	}

	credentials := model.Credentials{
		CredentialsSecure: model.CredentialsSecure{
			Role:     model.RoleUser,
			Username: inp.Username,
		},
		Password: inp.Password,
	}
	id, err := h.authService.SignUp(ctx.Request.Context(), credentials)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, E(err))
		return
	}

	if err := h.userService.Create(ctx.Request.Context(), id, inp.Username, inp.AvatarURL, inp.Nickname, model.RoleUser); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, E(err))
	}

	ctx.JSON(http.StatusOK, M("ok"))
}

type signInInput struct {
	Username string `json:"not username"`
	Password string `json:"password"`
}

type signInResponse struct {
	Token string `json:"token"`
}

// @Summary sign in
// @Tags auth
// @Param input body signInInput true "login credentials"
// @Router /api/auth/sign-in [post]
func (h AuthHandler) SignIn(ctx *gin.Context) {
	inp := new(signInInput)
	if err := ctx.BindJSON(inp); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, E(err))
		return
	}

	token, err := h.authService.SignIn(ctx.Request.Context(), inp.Username, inp.Password)
	if err != nil {
		if errors.Is(err, model.ErrWrongPassword) {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, E(err))
			return
		}

		if errors.Is(err, model.ErrUserNotFound) {
			ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, E(err))
			return
		}

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, E(err))
		return
	}

	ctx.JSON(http.StatusOK, signInResponse{Token: token})
}

// auth Middleware
func (m AuthHandler) WithAuth(role model.Role) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		credentials, err := m.authService.CheckAccess(ctx.Request.Context(), token, role)
		if err != nil {
			if errors.Is(err, model.ErrInvalidAccessToken) {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, E(err))
				return
			}
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, E(err))
			return
		}

		ctx.Set(model.CtxCredentialsKey, credentials)
	}
}
