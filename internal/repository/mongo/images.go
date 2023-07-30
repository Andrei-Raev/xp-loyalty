package mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/Andrei-Raev/xp-loyalty/internal/model"
)

type ImagesRepository struct {
	db *mongo.Collection
}

func NewImagesRepository(db *mongo.Database) *ImagesRepository {
	return &ImagesRepository{db: db.Collection("images")}
}

func (r *ImagesRepository) Create(ctx context.Context, img model.Image) error {
	_, err := r.db.InsertOne(ctx, toMongoImage(img))
	if err != nil {
		return fmt.Errorf("error images Create(): %w", err)
	}
	return nil
}

func (r *ImagesRepository) GetAvatars(ctx context.Context) ([]model.Image, error) {
	var images []mongoImage

	cursor, err := r.db.Find(ctx, bson.M{"type": model.TypeImageAvatar})
	if err != nil {
		return []model.Image{}, err
	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, &images)
	if err != nil {
		return nil, fmt.Errorf("error images GetAvatars(): %w", err)
	}
	return toModelImages(images), nil
}

func (r *ImagesRepository) GetPrizes(ctx context.Context) ([]model.Image, error) {
	var images []mongoImage

	cursor, err := r.db.Find(ctx, bson.M{"type": model.TypeImagePrize})
	if err != nil {
		return []model.Image{}, fmt.Errorf("error images GetPrizes(): %w", err)
	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, &images)
	if err != nil {
		return nil, fmt.Errorf("error images GetPrizes(): %w", err)
	}
	return toModelImages(images), nil
}

func (r *ImagesRepository) GetCardsBackgrounds(ctx context.Context) ([]model.Image, error) {
	var images []mongoImage

	cursor, err := r.db.Find(ctx, bson.M{"type": model.TypeCardBackground})
	if err != nil {
		return []model.Image{}, fmt.Errorf("error images GetCardsBackgrounds(): %w", err)
	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, &images)
	if err != nil {
		return nil, fmt.Errorf("error images GetCardsBackgrounds(): %w", err)
	}
	return toModelImages(images), nil
}

type mongoImage struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	URL  string             `bson:"url"`
	Type string             `bson:"type"`
}

func toMongoImage(i model.Image) mongoImage {
	id, _ := primitive.ObjectIDFromHex(i.ID)
	return mongoImage{
		ID:   id,
		URL:  i.URL,
		Type: i.Type,
	}
}

func toModelImage(i mongoImage) model.Image {
	return model.Image{
		ID:   i.ID.Hex(),
		URL:  i.URL,
		Type: i.Type,
	}
}

func toModelImages(imgs []mongoImage) []model.Image {
	modelImgs := make([]model.Image, len(imgs))
	for i := range imgs {
		modelImgs[i] = toModelImage(imgs[i])
	}
	return modelImgs
}
