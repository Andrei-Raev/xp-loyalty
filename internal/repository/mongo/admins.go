package mongo

import (
	"context"
	"errors"
	"fmt"

	"github.com/Andrei-Raev/xp-loyalty/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AdminsRepository struct {
	db *mongo.Collection
}

func NewAdminsRepository(db *mongo.Database) *AdminsRepository {
	return &AdminsRepository{db: db.Collection("admins")}
}

func (r *AdminsRepository) Create(ctx context.Context, admin model.Admin) error {
	if _, err := r.db.InsertOne(ctx, toMongoAdmin(admin)); err != nil {
		return fmt.Errorf("error admins Create(): %w", err)
	}
	return nil
}

func (r *AdminsRepository) GetByUsername(ctx context.Context, username string) (model.Admin, error) {
	var admin mongoAdmin

	err := r.db.FindOne(ctx, bson.M{"username": username}).Decode(&admin)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.Admin{}, model.ErrUserNotFound
		}
		return model.Admin{}, fmt.Errorf("erorr admins GetByUsername(): %w", err)
	}

	return toModelAdmin(admin), nil
}

type mongoAdmin struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Username string             `bson:"username"`
	Role     int                `bson:"role"`
}

func toModelAdmin(a mongoAdmin) model.Admin {
	admin := model.Admin{
		CredentialsSecure: model.CredentialsSecure{
			ID:       a.ID.Hex(),
			Username: a.Username,
			Role:     model.Role(a.Role),
		},
	}
	return admin
}

func toMongoAdmin(a model.Admin) mongoAdmin {
	id, _ := primitive.ObjectIDFromHex(a.ID)
	admin := mongoAdmin{
		ID:       id,
		Username: a.Username,
		Role:     int(a.Role),
	}
	return admin
}
