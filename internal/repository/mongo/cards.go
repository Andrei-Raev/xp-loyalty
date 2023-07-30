package mongo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/wintermonth2298/xp-loyalty/internal/model"
)

type CardsRepository struct {
	staticCardsDB *mongo.Collection
	cardsDB       *mongo.Collection
}

func NewCardsRepository(db *mongo.Database) *CardsRepository {
	return &CardsRepository{
		staticCardsDB: db.Collection("cards_static"),
		cardsDB:       db.Collection("cards"),
	}
}

func (r *CardsRepository) ViewCard(ctx context.Context, cardID string) error {
	_id, _ := primitive.ObjectIDFromHex(cardID)
	filter := bson.M{"_id": _id}
	update := bson.M{"$set": bson.M{"is_viewed": true}}

	if _, err := r.cardsDB.UpdateOne(ctx, filter, update); err != nil {
		return err
	}

	return nil
}

func (r *CardsRepository) Create(ctx context.Context, card model.Card) error {
	if _, err := r.cardsDB.InsertOne(ctx, toMongoCard(card)); err != nil {
		return fmt.Errorf("error cards Create(): %w", err)
	}
	return nil
}

func (r *CardsRepository) Update(ctx context.Context, card model.Card) error {
	_id, _ := primitive.ObjectIDFromHex(card.ID)

	match := bson.M{"_id": _id}
	if _, err := r.cardsDB.ReplaceOne(ctx, match, toMongoCard(card)); err != nil {
		return fmt.Errorf("error cards Update(): %w", err)
	}
	return nil
}

func (r *CardsRepository) DeleteUsersPendingDailyCards(ctx context.Context, username string) error {
	if _, err := r.cardsDB.DeleteMany(ctx, bson.M{"type": model.PoolDaily, "done": 0, "owner_username": username}); err != nil {
		return fmt.Errorf("error cards DeletePendingDailyCards(): %w", err)
	}
	return nil
}

func (r *CardsRepository) GetCardsByOwnerPool(ctx context.Context, ownerUsername, pool string) (model.Cards, error) {
	var cards mongoCards

	query := bson.M{"owner_username": ownerUsername, "static.pool": pool}

	queryOptions := options.Find()
	queryOptions.SetSort(bson.D{bson.E{Key: "chain_name", Value: 1}, {Key: "chain_order", Value: 1}})

	cursor, err := r.cardsDB.Find(ctx, query, queryOptions)
	if err != nil {
		return model.Cards{}, fmt.Errorf("error cards GetCardsByOwnerPool(): %w", err)
	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, &cards)
	if err != nil {
		return model.Cards{}, fmt.Errorf("error cards GetCardsByOwnerPool(): %w", err)
	}

	return toModelCards(cards), nil
}

func (r *CardsRepository) GetCardsByOwner(ctx context.Context, ownerUsername string) (model.Cards, error) {
	var cards mongoCards

	query := bson.M{"owner_username": ownerUsername}
	queryOptions := options.Find()
	queryOptions.SetSort(bson.D{bson.E{Key: "chain_name", Value: 1}, {Key: "chain_order", Value: 1}})

	cursor, err := r.cardsDB.Find(ctx, query, queryOptions)
	if err != nil {
		return model.Cards{}, fmt.Errorf("error cards GetCardsByOwner(): %w", err)
	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, &cards)
	if err != nil {
		return model.Cards{}, fmt.Errorf("error cards GetCardsByOwner(): %w", err)
	}

	return toModelCards(cards), nil
}

func (r *CardsRepository) Get(ctx context.Context, id string) (model.Card, error) {
	var card mongoCard

	_id, _ := primitive.ObjectIDFromHex(id)

	err := r.cardsDB.FindOne(ctx, bson.M{"_id": _id}).Decode(&card)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return model.Card{}, fmt.Errorf("error cards Get(): %w", model.ErrNoSuchCard)
	}

	return toModelCard(card), nil
}

func (r *CardsRepository) CreateStatic(ctx context.Context, card model.CardStatic) error {
	if _, err := r.staticCardsDB.InsertOne(ctx, toMongoCardStatic(card)); err != nil {
		return fmt.Errorf("error cards CreateStatic(): %w", err)
	}
	return nil
}

func (r *CardsRepository) DeleteStatic(ctx context.Context, ids []string) error {
	ids_ := toPrimitives(ids)

	if _, err := r.staticCardsDB.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": ids_}}); err != nil {
		return fmt.Errorf("error cards DeleteStatic(): %w", err)
	}
	return nil
}

func (r *CardsRepository) GetStaticByPool(ctx context.Context, pool string) (model.CardsStatic, error) {
	var cards mongoCardsStatic

	query := bson.M{"pool": pool}
	queryOptions := options.Find()
	queryOptions.SetSort(bson.D{bson.E{Key: "chain_name", Value: 1}, {Key: "chain_order", Value: 1}})

	cursor, err := r.staticCardsDB.Find(ctx, query, queryOptions)
	if err != nil {
		return model.CardsStatic{}, fmt.Errorf("error cards GetStaticByPool(): %w", err)
	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, &cards)
	if err != nil {
		return model.CardsStatic{}, fmt.Errorf("error cards GetStaticByPool(): %w", err)
	}

	return toModelCardsStatic(cards), nil
}

func (r *CardsRepository) GetStatic(ctx context.Context, ids []string) (model.CardsStatic, error) {
	var cards mongoCardsStatic

	var query bson.M
	if len(ids) != 0 && ids[0] == "*" {
		query = bson.M{}
	} else {

		ids_ := toPrimitives(ids)
		query = bson.M{"_id": bson.M{"$in": ids_}}
	}

	queryOptions := options.Find()
	queryOptions.SetSort(bson.D{bson.E{Key: "chain_name", Value: 1}, {Key: "chain_order", Value: 1}})

	cursor, err := r.staticCardsDB.Find(ctx, query, queryOptions)
	if err != nil {
		return model.CardsStatic{}, fmt.Errorf("error cards GetStatic(): %w", err)
	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, &cards)
	if err != nil {
		return model.CardsStatic{}, fmt.Errorf("error cards GetStatic(): %w", err)
	}

	return toModelCardsStatic(cards), nil
}

type mongoCard struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	OwnerUsername string             `bson:"owner_username"`
	Static        mongoCardStatic    `bson:"static"`
	Done          int                `bson:"done"`
	Progress      *int               `bson:"progress,omitempty"`
	IsViewed      bool               `bson:"is_viewed"`
	History       []int              `bson:"history,omitempty"`
	OptDoneNum    int                `bson:"opt_done_num"`
}

type mongoCardStatic struct {
	ID               primitive.ObjectID `bson:"_id,omitempty"`
	Title            string             `bson:"title"`
	ShortDescription string             `bson:"short_description"`
	LongDescription  string             `bson:"long_description"`
	Goal             string             `bson:"goal"`
	CreatedAt        time.Time          `bson:"created_at"`
	Type             string             `bson:"type"`
	Pool             string             `bson:"pool"`
	BackgroundURL    string             `bson:"background_url"`
	ChainName        *string            `bson:"chain_name,omitempty"`
	ChainOrder       *int               `bson:"chain_order,omitempty"`
	OrdSettings      *mongoOrdSettings  `bson:"ordinary_settings,omitempty"`
	PrgSettings      *mongoPrgSettings  `bson:"progress_settings,omitempty"`
	OptSettings      *mongoOptSettings  `bson:"options_settings,omitempty"`
}

type mongoOrdSettings struct {
	Award mongoAward `bson:"award"`
}

type mongoPrgSettings struct {
	Award       mongoAward `bson:"award"`
	MaxProgress int        `bson:"max_progress"`
}

type mongoOptSettings struct {
	Awards  []mongoAward `bson:"awards"`
	Options []float32    `bson:"options"`
}

type mongoAward struct {
	XPoints       int    `bson:"XPoints"`
	Prize         string `bson:"prize"`
	PrizeImageURL string `bson:"prize_image_url"`
}

type mongoCards []mongoCard
type mongoCardsStatic []mongoCardStatic

func toMongoAwards(a []model.Award) []mongoAward {
	awards := make([]mongoAward, len(a))
	for i := range a {
		awards[i] = mongoAward(a[i])
	}
	return awards
}

func toModelAwards(a []mongoAward) []model.Award {
	awards := make([]model.Award, len(a))
	for i := range a {
		awards[i] = model.Award(a[i])
	}
	return awards
}

func toMongoCard(c model.Card) mongoCard {
	var progress *int
	if c.Static.Type == model.TypeProgress {
		progress = &c.Progress
	}

	id, _ := primitive.ObjectIDFromHex(c.ID)
	mongoCard := mongoCard{
		ID:            id,
		OwnerUsername: c.OwnerUsername,
		Static:        toMongoCardStatic(c.Static),
		Done:          c.Done,
		IsViewed:      c.IsViewed,
		Progress:      progress,
		History:       c.History,
		OptDoneNum:    c.OptDoneNum,
	}
	return mongoCard
}

func toModelCard(c mongoCard) model.Card {
	var progress int
	if c.Static.Type == model.TypeProgress {
		progress = *c.Progress
	}

	mongoCard := model.Card{
		ID:            c.ID.Hex(),
		OwnerUsername: c.OwnerUsername,
		Static:        toModelCardStatic(c.Static),
		Done:          c.Done,
		IsViewed:      c.IsViewed,
		Progress:      progress,
		History:       c.History,
		OptDoneNum:    c.OptDoneNum,
	}
	return mongoCard
}

func toMongoCardStatic(c model.CardStatic) mongoCardStatic {
	var chainName *string
	var chainOrder *int
	if c.Pool == model.PoolConst {
		chainName = &c.ChainName
		chainOrder = &c.ChainOrder
	}

	var ordsettings *mongoOrdSettings
	var prgsettings *mongoPrgSettings
	var optsettings *mongoOptSettings
	if c.Type == model.TypeOrdinary {
		ordsettings = &mongoOrdSettings{
			Award: mongoAward(c.OrdSettings.Award),
		}
	} else if c.Type == model.TypeProgress {
		prgsettings = &mongoPrgSettings{
			Award:       mongoAward(c.PrgSettings.Award),
			MaxProgress: c.PrgSettings.MaxProgress,
		}
	} else if c.Type == model.TypeOptions {
		optsettings = &mongoOptSettings{
			Awards:  toMongoAwards(c.OptSettings.Awards),
			Options: c.OptSettings.Options,
		}
	}

	id, _ := primitive.ObjectIDFromHex(c.ID)
	mongoStaticCard := mongoCardStatic{
		ID:               id,
		Title:            c.Title,
		ShortDescription: c.ShortDescription,
		LongDescription:  c.LongDescription,
		Goal:             c.Goal,
		CreatedAt:        c.CreatedAt,
		Type:             c.Type,
		Pool:             c.Pool,
		ChainName:        chainName,
		BackgroundURL:    c.BackgroundURL,
		ChainOrder:       chainOrder,
		OrdSettings:      ordsettings,
		PrgSettings:      prgsettings,
		OptSettings:      optsettings,
	}
	return mongoStaticCard
}

func toModelCardStatic(c mongoCardStatic) model.CardStatic {
	var chainName string
	var chainOrder int
	if c.Pool == model.PoolConst {
		chainName = *c.ChainName
		chainOrder = *c.ChainOrder
	}

	var ordsettings *model.OrdSettings
	var prgsettings *model.PrgSettings
	var optsettings *model.OptSettings
	if c.Type == model.TypeOrdinary {
		ordsettings = &model.OrdSettings{
			Award: model.Award(c.OrdSettings.Award),
		}
	} else if c.Type == model.TypeProgress {
		prgsettings = &model.PrgSettings{
			Award:       model.Award(c.PrgSettings.Award),
			MaxProgress: c.PrgSettings.MaxProgress,
		}
	} else if c.Type == model.TypeOptions {
		optsettings = &model.OptSettings{
			Awards:  toModelAwards(c.OptSettings.Awards),
			Options: c.OptSettings.Options,
		}
	}

	mongoStaticCard := model.CardStatic{
		ID:               c.ID.Hex(),
		Title:            c.Title,
		ShortDescription: c.ShortDescription,
		LongDescription:  c.LongDescription,
		Goal:             c.Goal,
		CreatedAt:        c.CreatedAt,
		Type:             c.Type,
		Pool:             c.Pool,
		ChainName:        chainName,
		BackgroundURL:    c.BackgroundURL,
		ChainOrder:       chainOrder,
		OrdSettings:      ordsettings,
		PrgSettings:      prgsettings,
		OptSettings:      optsettings,
	}
	return mongoStaticCard
}

func toModelCards(c mongoCards) model.Cards {
	cards := make(model.Cards, len(c))
	for i := range c {
		cards[i] = toModelCard(c[i])
	}
	return cards
}

func toModelCardsStatic(c mongoCardsStatic) model.CardsStatic {
	cards := make(model.CardsStatic, len(c))
	for i := range c {
		cards[i] = toModelCardStatic(c[i])
	}
	return cards
}
