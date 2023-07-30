package mongo

import (
	"context"
	"errors"
	"fmt"

	"github.com/wintermonth2298/xp-loyalty/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CredentialsRepository struct {
	db *mongo.Collection
}

func NewCredentialsRepository(db *mongo.Database) *CredentialsRepository {
	return &CredentialsRepository{db: db.Collection("credentials")}
}

func (r *CredentialsRepository) Create(ctx context.Context, credentials model.Credentials) (string, error) {
	doc, err := r.db.InsertOne(ctx, toMongoCredentials(credentials))
	id, ok := doc.InsertedID.(primitive.ObjectID)
	if !ok {
		err = model.ErrInterfaceCast
	}
	if err != nil {
		return "", fmt.Errorf("error credentials Create(): %w", err)
	}
	return id.Hex(), nil
}

func (r *CredentialsRepository) GetByUsername(ctx context.Context, username string) (model.Credentials, error) {
	var credentials mongoCredentials

	err := r.db.FindOne(ctx, bson.M{"username": username}).Decode(&credentials)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.Credentials{}, fmt.Errorf("error credentials GetByUsername(): %w", model.ErrUserNotFound)
		}
		return model.Credentials{}, fmt.Errorf("error credentials GetByUsername(): %w", err)
	}

	return toModelCredentials(credentials), nil
}

func (r *CredentialsRepository) GetByCredentials(ctx context.Context, username, password string) (model.Credentials, error) {
	var credentials mongoCredentials

	err := r.db.FindOne(ctx, bson.M{"username": username, "password": password}).Decode(&credentials)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.Credentials{}, model.ErrUserNotFound
		}
		return model.Credentials{}, fmt.Errorf("error credentials GetByCredentials(): %w", err)
	}

	return toModelCredentials(credentials), nil
}

type mongoCredentials struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Username string             `bson:"username"`
	Role     int                `bson:"role"`
	Password string             `bson:"password"`
}

func toMongoCredentials(c model.Credentials) mongoCredentials {
	id, _ := primitive.ObjectIDFromHex(c.ID)
	creds := mongoCredentials{
		ID:       id,
		Username: c.Username,
		Role:     int(c.Role),
		Password: c.Password,
	}
	return creds
}

func toModelCredentials(c mongoCredentials) model.Credentials {
	creds := model.Credentials{
		CredentialsSecure: model.CredentialsSecure{
			ID:       c.ID.Hex(),
			Username: c.Username,
			Role:     model.Role(c.Role),
		},
		Password: c.Password,
	}
	return creds
}
