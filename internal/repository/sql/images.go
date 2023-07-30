package sql

import (
	"context"
	"fmt"
	"strconv"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/wintermonth2298/xp-loyalty/internal/model"
)

type ImagesRepository struct {
	db *sqlx.DB
}

func NewImagesRepository(db *sqlx.DB) *ImagesRepository {
	return &ImagesRepository{db: db}
}

func (r *ImagesRepository) Create(ctx context.Context, img model.Image) error {
	i := toSQLImage(img)

	query, args, err := psql.Insert("img").Columns("url", "type").Values(i.URL, i.Type).ToSql()
	if err != nil {
		return fmt.Errorf("imagesRepo - Create() - sq: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("imagesRepo - Create() - ExecContext(): %w", err)
	}

	return nil
}

func (r *ImagesRepository) GetAvatars(ctx context.Context) ([]model.Image, error) {
	var images []Image
	result, err := getAll(ctx, r.db, "img", sq.Eq{"img.type": model.TypeImageAvatar}, images)
	if err != nil {
		return nil, fmt.Errorf("imgsRepo - GetAvatars(): %w", err)
	}
	return toModelImages(result), nil
}

func (r *ImagesRepository) GetPrizes(ctx context.Context) ([]model.Image, error) {
	var images []Image
	result, err := getAll(ctx, r.db, "img", sq.Eq{"img.type": model.TypeImagePrize}, images)
	if err != nil {
		return nil, fmt.Errorf("imgsRepo - GetPrizes(): %w", err)
	}
	return toModelImages(result), nil
}

type Image struct {
	ID   int    `bd:"id"`
	URL  string `bd:"url"`
	Type string `bd:"type"`
}

func toSQLImage(u model.Image) Image {
	id, _ := strconv.Atoi(u.ID)
	return Image{
		ID:   id,
		URL:  u.URL,
		Type: u.Type,
	}
}

func toModelImages(u []Image) []model.Image {
	images := make([]model.Image, len(u))
	for i, img := range u {
		images[i] = model.Image{
			ID:   strconv.Itoa(img.ID),
			URL:  img.URL,
			Type: img.Type,
		}
	}
	return images
}
