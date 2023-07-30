package handler

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/wintermonth2298/xp-loyalty/internal/model"
)

type ImageService interface {
	Create(ctx context.Context, url string, imgtype string) error
	GetAvatars(ctx context.Context) ([]model.Image, error)
	GetPrizes(ctx context.Context) ([]model.Image, error)
	GetCardsBackgrounds(ctx context.Context) ([]model.Image, error)
}

type ImageHandler struct {
	imageService ImageService
}

func NewImagesHandler(imgService ImageService) *ImageHandler {
	return &ImageHandler{imageService: imgService}
}

// @Summary upload avatar image
// @Tags images
// @Accept multipart/form-data
// @Param image formData file true "avatar image"
// @Router /api/images/upload/avatar [post]
// @Security ApiKeyAuth
func (h ImageHandler) UploadAvatarImage(ctx *gin.Context) {
	file, err := ctx.FormFile("image")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, E(err))
		return
	}

	name := uniqueImagName(file.Filename, model.TypeImageAvatar)
	if err := ctx.SaveUploadedFile(file, name); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, E(err))
		return
	}

	imageURL := fmt.Sprintf("http://localhost:8000/%s", name)

	data := map[string]interface{}{
		"image_url": imageURL,
		"header":    file.Header,
		"size":      file.Size,
	}

	if err := h.imageService.Create(ctx.Request.Context(), imageURL, model.TypeImageAvatar); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, E(err))
	}
	ctx.JSON(http.StatusOK, data)
}

// @Summary upload prize image
// @Tags images
// @Accept multipart/form-data
// @Param image formData file true "prize image"
// @Router /api/images/upload/prize [post]
// @Security ApiKeyAuth
func (h ImageHandler) UploadPrizeImage(ctx *gin.Context) {
	file, err := ctx.FormFile("image")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, E(err))
		return
	}

	name := uniqueImagName(file.Filename, model.TypeImagePrize)
	if err := ctx.SaveUploadedFile(file, name); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, E(err))
		return
	}

	imageURL := fmt.Sprintf("http://localhost:8000/%s", name)

	data := map[string]interface{}{
		"image_url": imageURL,
		"header":    file.Header,
		"size":      file.Size,
	}

	if err := h.imageService.Create(ctx.Request.Context(), imageURL, model.TypeImagePrize); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, E(err))
		return
	}

	ctx.JSON(http.StatusOK, data)
}

// @Summary upload card background image
// @Tags images
// @Accept multipart/form-data
// @Param image formData file true "prize image"
// @Router /api/images/upload/card-background [post]
// @Security ApiKeyAuth
func (h ImageHandler) UploadCardsBackground(ctx *gin.Context) {
	file, err := ctx.FormFile("image")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, E(err))
		return
	}

	name := uniqueImagName(file.Filename, model.TypeCardBackground)
	if err := ctx.SaveUploadedFile(file, name); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, E(err))
		return
	}

	imageURL := fmt.Sprintf("http://localhost:8000/%s", name)

	data := map[string]interface{}{
		"image_url": imageURL,
		"header":    file.Header,
		"size":      file.Size,
	}

	if err := h.imageService.Create(ctx.Request.Context(), imageURL, model.TypeCardBackground); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, E(err))
		return
	}

	ctx.JSON(http.StatusOK, data)
}

type imagesResponse struct {
	Images []model.Image `json:"images"`
}

// @Summary get avatar images
// @Tags images
// @Router /api/images/avatar [get]
// @Security ApiKeyAuth
func (h ImageHandler) GetAvatarImages(ctx *gin.Context) {
	images, err := h.imageService.GetAvatars(ctx.Request.Context())
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, E(err))
		return
	}

	ctx.JSON(http.StatusOK, imagesResponse{Images: images})
}

// @Summary get prize images
// @Tags images
// @Router /api/images/prize [get]
// @Security ApiKeyAuth
func (h ImageHandler) GetPrizeImages(ctx *gin.Context) {
	images, err := h.imageService.GetPrizes(ctx.Request.Context())
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, E(err))
		return
	}

	ctx.JSON(http.StatusOK, imagesResponse{Images: images})
}

// @Summary get background images
// @Tags images
// @Router /api/images/card-background [get]
// @Security ApiKeyAuth
func (h ImageHandler) GetCardsBackgrounds(ctx *gin.Context) {
	images, err := h.imageService.GetCardsBackgrounds(ctx.Request.Context())
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, E(err))
		return
	}

	ctx.JSON(http.StatusOK, imagesResponse{Images: images})
}

func uniqueImagName(imageName string, imgtype string) string {
	uniqueID := uuid.New()
	unique := strings.Replace(uniqueID.String(), "-", "", -1)

	ext := strings.Split(imageName, ".")[1]
	name := strings.Split(imageName, ".")[0]

	result := fmt.Sprintf("%s-%s.%s", name, unique, ext)
	result = fmt.Sprintf("static/images/%s/%s", imgtype, result)

	return result
}
