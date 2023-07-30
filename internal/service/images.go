package service

import (
	"context"

	"github.com/wintermonth2298/xp-loyalty/internal/model"
)

type ImageRepository interface {
	Create(ctx context.Context, img model.Image) error
	GetAvatars(ctx context.Context) ([]model.Image, error)
	GetPrizes(ctx context.Context) ([]model.Image, error)
	GetCardsBackgrounds(ctx context.Context) ([]model.Image, error)
}

type ImageService struct {
	imageRepo ImageRepository
}

func NewImagesService(imgRepo ImageRepository) *ImageService {
	return &ImageService{imageRepo: imgRepo}
}

func (s *ImageService) Create(ctx context.Context, url string, imgtype string) error {
	return s.imageRepo.Create(ctx, model.Image{URL: url, Type: imgtype})
}

func (s *ImageService) GetAvatars(ctx context.Context) ([]model.Image, error) {
	return s.imageRepo.GetAvatars(ctx)
}

func (s *ImageService) GetPrizes(ctx context.Context) ([]model.Image, error) {
	return s.imageRepo.GetPrizes(ctx)
}

func (s *ImageService) GetCardsBackgrounds(ctx context.Context) ([]model.Image, error) {
	return s.imageRepo.GetCardsBackgrounds(ctx)
}
