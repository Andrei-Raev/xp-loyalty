package mongo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/wintermonth2298/xp-loyalty/internal/model"
)

type UsersRepository struct {
	db *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *UsersRepository {
	return &UsersRepository{db: db.Collection("users")}
}

func (r *UsersRepository) Create(ctx context.Context, user model.User) error {
	_, err := r.db.InsertOne(ctx, toMongoUser(user))
	if err != nil {
		return fmt.Errorf("error users Create(): %w", err)
	}
	return nil
}

func (r *UsersRepository) GetByUsername(ctx context.Context, username string) (model.User, error) {
	var user mongoUser

	err := r.db.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return model.User{}, fmt.Errorf("error users GetByUsername(): %w", model.ErrUserNotFound)
	}

	if err != nil {
		return model.User{}, fmt.Errorf("error users GetByUsername(): %w", err)
	}

	return toModelUser(user), nil
}

func (r *UsersRepository) GetAll(ctx context.Context) ([]model.User, error) {
	var users []mongoUser

	cursor, err := r.db.Find(ctx, bson.M{})
	if err != nil {
		return []model.User{}, fmt.Errorf("error users GetAll(): %w", err)
	}
	defer cursor.Close(ctx)
	err = cursor.All(ctx, &users)
	if err != nil {
		return nil, fmt.Errorf("error users GetAll(): %w", err)
	}

	return toModelUsers(users), nil
}

func (r *UsersRepository) Update(ctx context.Context, user model.User) error {
	match := bson.M{"username": user.Username}
	_, err := r.db.ReplaceOne(ctx, match, toMongoUser(user))
	if err != nil {
		return fmt.Errorf("error users Update(): %w", err)
	}
	return nil

}

type mongoUser struct {
	ID                   primitive.ObjectID `bson:"_id,omitempty"`
	Username             string             `bson:"username"`
	Role                 model.Role         `bson:"role"`
	Nickname             string             `bson:"nickname"`
	AvatarURL            string             `bson:"avatar_url"`
	XPoints              int                `bson:"XPoints"`
	RegistrationTime     time.Time          `bson:"registration_time"`
	LastDailyCardsUpdate time.Time          `bson:"last_daily_cards_update"`
	Prizes               []model.UserPrize  `bson:"prizes"`
}

func toMongoUser(u model.User) mongoUser {
	id, _ := primitive.ObjectIDFromHex(u.ID)
	return mongoUser{
		ID:                   id,
		Username:             u.Username,
		Role:                 u.Role,
		Nickname:             u.Nickname,
		AvatarURL:            u.AvatarURL,
		XPoints:              u.XPoints,
		RegistrationTime:     u.RegistrationTime,
		LastDailyCardsUpdate: u.LastDailyCardsUpdate,
		Prizes:               u.Prizes,
	}
}

func toModelUser(u mongoUser) model.User {
	return model.User{
		CredentialsSecure: model.CredentialsSecure{
			ID:       u.ID.Hex(),
			Username: u.Username,
			Role:     u.Role,
		},
		Nickname:             u.Nickname,
		AvatarURL:            u.AvatarURL,
		XPoints:              u.XPoints,
		RegistrationTime:     u.RegistrationTime,
		LastDailyCardsUpdate: u.LastDailyCardsUpdate,
		Prizes:               u.Prizes,
	}
}

func toModelUsers(u []mongoUser) []model.User {
	users := make([]model.User, len(u))
	for i := range u {
		users[i] = toModelUser(u[i])
	}
	return users
}
