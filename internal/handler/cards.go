package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/wintermonth2298/xp-loyalty/internal/model"
)

type CardsService interface {
	GetStatic(ctx context.Context, ids []string) (model.CardsStatic, error)
	CreateStatic(ctx context.Context, card model.CardStatic) error
	DeleteStatic(ctx context.Context, ids []string) (err error)
	Update(ctx context.Context, id string, progress int, doneOption float32) (string, int, error)
	GetFormattedCards(ctx context.Context, ownerUsername string) (model.Cards, model.Cards, error)
	ViewCard(ctx context.Context, id string) error
}

type CardsUserService interface {
	GetByUsername(ctx context.Context, username string) (model.User, error)
	Update(ctx context.Context, user model.User) error
}

type CardsHandler struct {
	cardsService CardsService
	userService  CardsUserService
}

func NewCardsStaticHandler(cardsService CardsService, userService CardsUserService) *CardsHandler {
	return &CardsHandler{cardsService: cardsService, userService: userService}
}

type createStaticCardInput struct {
	Title            string             `json:"title"`
	ShortDescription string             `json:"short_description"`
	LongDescription  string             `json:"long_description"`
	Goal             string             `json:"goal"`
	Type             string             `json:"type"`
	Pool             string             `json:"pool"`
	ChainName        string             `json:"chain_name"`
	BackgroundURL    string             `json:"background_url"`
	ChainOrder       int                `json:"chain_order"`
	OrdSettings      *model.OrdSettings `json:"ordinary_settings"`
	PrgSettings      *model.PrgSettings `json:"progress_settings"`
	OptSettings      *model.OptSettings `json:"options_settings"`
}

// @Summary create static card
// @Tags cards
// @Param input body createStaticCardInput false "create static card input"
// @Router /api/cards [post]
// @Security ApiKeyAuth
func (h CardsHandler) CreateStatic(ctx *gin.Context) {
	inp := new(createStaticCardInput)
	if err := ctx.BindJSON(inp); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, E(err))
		return
	}

	card := model.CardStatic{
		Title:            inp.Title,
		ShortDescription: inp.ShortDescription,
		LongDescription:  inp.LongDescription,
		Goal:             inp.Goal,
		Type:             inp.Type,
		Pool:             inp.Pool,
		ChainName:        inp.ChainName,
		ChainOrder:       inp.ChainOrder,
		BackgroundURL:    inp.BackgroundURL,
		OrdSettings:      inp.OrdSettings,
		PrgSettings:      inp.PrgSettings,
		OptSettings:      inp.OptSettings,
	}

	if err := h.cardsService.CreateStatic(ctx.Request.Context(), card); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, E(err))
		return
	}

	ctx.JSON(http.StatusOK, M("ok"))
}

type deleteCardStaticsInput struct {
	IDs []string `json:"ids"`
}

// @Summary delete cards
// @Tags cards
// @Param input body deleteCardStaticsInput true "delete static cards input"
// @Router /api/cards [delete]
// @Security ApiKeyAuth
func (h CardsHandler) DeleteStatic(ctx *gin.Context) {
	inp := new(deleteCardStaticsInput)
	if err := ctx.BindJSON(inp); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, E(err))
		return
	}

	if err := h.cardsService.DeleteStatic(ctx.Request.Context(), inp.IDs); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, E(err))
		return
	}

	ctx.JSON(http.StatusOK, M("ok"))
}

type updateCardInput struct {
	ID         string  `json:"card_id"`
	Progress   int     `json:"progress"`
	DoneOption float32 `json:"done_option"`
}

// @Summary update card
// @Tags cards
// @Param input body updateCardInput true "update card input"
// @Router /api/cards/done [post]
// @Security ApiKeyAuth
func (h CardsHandler) UpdateCard(ctx *gin.Context) {
	inp := new(updateCardInput)
	if err := ctx.BindJSON(inp); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, E(err))
		return
	}

	username, xpoints, err := h.cardsService.Update(ctx.Request.Context(), inp.ID, inp.Progress, inp.DoneOption)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, E(err))
		return
	}

	user, err := h.userService.GetByUsername(ctx.Request.Context(), username)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, E(err))
		return
	}

	user.XPoints += xpoints
	if err := h.userService.Update(ctx.Request.Context(), user); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, E(err))
		return
	}

	ctx.JSON(http.StatusOK, M("ok"))
}

type getUserCardsResponse struct {
	PendingCards []model.Card `json:"pending_cards"`
	DoneCards    []model.Card `json:"done_cards"`
}

// @Summary get all user cards by username
// @Tags cards
// @Param username path string true "username"
// @Router /api/cards/{username} [get]
// @Security ApiKeyAuth
func (h CardsHandler) GetUserCards(ctx *gin.Context) {
	username, err := ParsePath(ctx, "username")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, E(err))
	}
	pending, done, err := h.cardsService.GetFormattedCards(ctx.Request.Context(), username)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, E(err))
		return
	}
	ctx.JSON(http.StatusOK, getUserCardsResponse{PendingCards: pending, DoneCards: done})
}

// @Summary get all user cards by token
// @Tags cards
// @Router /api/cards/profile [get]
// @Security ApiKeyAuth
func (h CardsHandler) GetProfileCards(ctx *gin.Context) {
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

	pending, done, err := h.cardsService.GetFormattedCards(ctx.Request.Context(), credentials.Username)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, E(err))
		return
	}
	ctx.JSON(http.StatusOK, getUserCardsResponse{PendingCards: pending, DoneCards: done})
}

type getStaticCardsResponse struct {
	Cards model.CardsStatic `json:"cards"`
}

// @Summary get all static cards
// @Tags cards
// @Router /api/cards [get]
// @Security ApiKeyAuth
func (h CardsHandler) GetAllStatic(ctx *gin.Context) {
	cards, err := h.cardsService.GetStatic(ctx.Request.Context(), []string{"*"})
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, E(err))
		return
	}
	ctx.JSON(http.StatusOK, getStaticCardsResponse{Cards: cards})
}

type viewCardInput struct {
	CardID string `json:"card_id"`
}

// @Summary view card
// @Tags cards
// @Param input body viewCardInput true "view card"
// @Router /api/cards/view [post]
// @Security ApiKeyAuth
func (h CardsHandler) ViewCard(ctx *gin.Context) {
	inp := new(viewCardInput)
	if err := ctx.BindJSON(inp); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, E(err))
		return
	}

	err := h.cardsService.ViewCard(ctx.Request.Context(), inp.CardID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, E(err))
		return
	}

	ctx.JSON(http.StatusOK, M("ok"))
}
