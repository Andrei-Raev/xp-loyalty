package sql

import (
	"context"
	"sort"
	"strconv"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/wintermonth2298/xp-loyalty/internal/model"
)

type CardsRepository struct {
	db *sqlx.DB
}

func NewCardsRepository(db *sqlx.DB) *CardsRepository {
	return &CardsRepository{db: db}
}

func (r *CardsRepository) Create(ctx context.Context, admin model.Admin) error {
	return nil
}

func (r *CardsRepository) GetStatic(ctx context.Context, ids []string) ([]model.CardStatic, error) {
	if len(ids) != 0 && ids[0] == "*" {
		return r.getStatic(ctx, sq.Eq{})
	}
	return nil, nil
}

func (r *CardsRepository) GetStaticByPool(ctx context.Context, pool string) ([]model.CardStatic, error) {
	return r.getStatic(ctx, sq.Eq{"pool": pool})
}

func (r *CardsRepository) getStatic(ctx context.Context, eq sq.Eq) ([]model.CardStatic, error) {
	var rows []CardStaticRow

	query, args, err := psql.
		Select("card_static.*, award.xpoints, award.prize, award.prize_img_url, award.opt").
		From("card_static").LeftJoin("award ON award.card_static_id = card_static.id").
		Where(eq).
		OrderBy("card_static.id, award.opt").
		ToSql()
	if err != nil {
		return nil, err
	}

	if err := r.db.SelectContext(ctx, &rows, query, args...); err != nil {
		return nil, err
	}

	cards := ToModelCardStatic(rows)
	sort.Slice(cards, func(i, j int) bool {
		if cards[i].ChainName != cards[j].ChainName {
			return cards[i].ChainName < cards[j].ChainName
		}
		return cards[i].ChainOrder < cards[j].ChainOrder
	})

	return cards, nil
}

type CardStaticRow struct {
	ID               int       `db:"id"`
	Title            string    `db:"title"`
	ShortDescription string    `db:"short_description"`
	LongDescription  string    `db:"long_description"`
	Goal             string    `db:"goal"`
	CreatedAt        time.Time `db:"created_at"`
	Type             string    `db:"type"`
	Pool             string    `db:"pool"`
	SettingID        int       `db:"setting_id"`
	XPoints          int       `db:"xpoints"`
	Prize            string    `db:"prize"`
	PrizeImageURL    string    `db:"prize_img_url"`
	ChainName        *string   `db:"chain_name"`
	ChainOrder       *int      `db:"chain_order"`
	MaxProgress      *int      `db:"max_progress"`
	Opt              *float32  `db:"opt"`
}

func ToModelCardStatic(rows []CardStaticRow) []model.CardStatic {
	var cards []model.CardStatic

	for i := 0; i < len(rows); i++ {
		row := rows[i]

		card := model.CardStatic{
			ID:               strconv.Itoa(row.ID),
			Title:            row.Title,
			ShortDescription: row.ShortDescription,
			LongDescription:  row.LongDescription,
			Goal:             row.Goal,
			CreatedAt:        row.CreatedAt,
			Type:             row.Type,
			Pool:             row.Pool,
			ChainName:        "",
			ChainOrder:       0,
		}

		if row.Pool == model.PoolConst {
			card.ChainName = *row.ChainName
			card.ChainOrder = *row.ChainOrder
		}

		switch row.Type {
		case model.TypeOptions:
			awards := make([]model.Award, 0, 5)
			opts := make([]float32, 0, 5)
			var j int
			for j = i; j < len(rows); j++ {
				if rows[i].ID != rows[j].ID {
					break
				}
				a := model.Award{
					XPoints:       rows[j].XPoints,
					Prize:         rows[j].Prize,
					PrizeImageURL: rows[j].PrizeImageURL,
				}
				awards = append(awards, a)
				opts = append(opts, *rows[j].Opt)
			}
			i = j - 1
			card.OptSettings = &model.OptSettings{
				Awards:  awards,
				Options: opts,
			}
		case model.TypeProgress:
			card.PrgSettings = &model.PrgSettings{
				Award: model.Award{
					XPoints:       row.XPoints,
					Prize:         row.Prize,
					PrizeImageURL: row.PrizeImageURL,
				},
				MaxProgress: *row.MaxProgress,
			}
		case model.TypeOrdinary:
			card.OrdSettings = &model.OrdSettings{
				Award: model.Award{
					XPoints:       row.XPoints,
					Prize:         row.Prize,
					PrizeImageURL: row.PrizeImageURL,
				},
			}
		}

		cards = append(cards, card)
	}

	return cards
}
