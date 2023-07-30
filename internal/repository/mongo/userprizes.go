package mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/wintermonth2298/xp-loyalty/internal/model"
)

type AwardsRepository struct {
	db *mongo.Collection
}

func NewAwardsRepository(db *mongo.Database) *AwardsRepository {
	return &AwardsRepository{db: db.Collection("awards")}
}

func (r *AwardsRepository) Add(ctx context.Context, award model.UserPrize) error {
	_, err := r.db.InsertOne(ctx, mongoPrize(award))
	if err != nil {
		return fmt.Errorf("error prizes Create(): %w", err)
	}
	return nil
}

func (r *AwardsRepository) GetByUsername(ctx context.Context, username string) ([]model.UserPrize, error) {
	var prizes []mongoPrize

	cursor, err := r.db.Find(ctx, bson.M{"owner_username": username})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, &prizes)
	if err != nil {
		return nil, fmt.Errorf("error images GetAvatars(): %w", err)
	}

	prizesDomain := make([]model.UserPrize, len(prizes))
	for i := range prizes {
		prizesDomain[i] = model.UserPrize(prizes[i])
	}

	return prizesDomain, nil
}

type mongoPrize struct {
	URL           string `bson:"url"`
	OwnerUsername string `bson:"owner_username"`
	Available     bool   `bson:"available"`
}
