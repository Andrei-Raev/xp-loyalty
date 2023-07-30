package sql

import (
	"context"
	"fmt"
	"strconv"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/wintermonth2298/xp-loyalty/internal/model"
)

type UsersRepository struct {
	db *sqlx.DB
}

func NewUsersRepository(db *sqlx.DB) *UsersRepository {
	return &UsersRepository{db: db}
}

func (r *UsersRepository) Create(ctx context.Context, user model.User) error {
	u := toSQLUser(user)

	query, args, err := psql.Insert("usr").
		Columns("username",
			"role",
			"nickname",
			"avatar_url",
			"xpoints",
			"registration_time",
			"last_daily_cards_update",
		).
		Values(u.Username,
			u.Role,
			u.Nickname,
			u.AvatarURL,
			u.XPoints,
			u.RegistrationTime,
			u.LastDailyCardsUpdate,
		).ToSql()
	if err != nil {
		return fmt.Errorf("error users Create(): %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("error users Create(): %w", err)
	}

	return nil
}

func (r *UsersRepository) GetByUsername(ctx context.Context, username string) (model.User, error) {
	var user User
	result, err := get(ctx, r.db, "usr", sq.Eq{"usr.username": username}, user)
	if err != nil {
		return model.User{}, fmt.Errorf("usersRepo - GetByUsername(): %w", err)
	}
	return toModelUser(result), nil
}

func (r *UsersRepository) GetAll(ctx context.Context) ([]model.User, error) {
	var users []User
	result, err := getAll(ctx, r.db, "usr", sq.Eq{}, users)
	if err != nil {
		return nil, fmt.Errorf("usersRepo - GetAll(): %w", err)
	}
	return toModelUsers(result), nil
}

func (r *UsersRepository) Update(ctx context.Context, user model.User) error {
	u := toSQLUser(user)

	query, args, err := psql.Update("usr").
		SetMap(map[string]interface{}{
			"username":                u.Username,
			"role":                    u.Role,
			"nickname":                u.Nickname,
			"avatar_url":              u.AvatarURL,
			"xpoints":                 u.XPoints,
			"registration_time":       u.RegistrationTime,
			"last_daily_cards_update": u.LastDailyCardsUpdate,
		}).Where(sq.Eq{"id": u.ID}).ToSql()
	if err != nil {
		return fmt.Errorf("usersRepo - GetAll() - sq: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("usersRepo - Create() - ExecContext(): %w", err)
	}

	return nil
}

type User struct {
	ID                   int       `db:"id"`
	Username             string    `db:"username"`
	Role                 int       `db:"role"`
	Nickname             string    `db:"nickname"`
	AvatarURL            string    `db:"avatar_url"`
	XPoints              int       `db:"xpoints"`
	RegistrationTime     time.Time `db:"registration_time"`
	LastDailyCardsUpdate time.Time `db:"last_daily_cards_update"`
}

func toSQLUser(u model.User) User {
	id, _ := strconv.Atoi(u.ID)
	return User{
		ID:                   id,
		Username:             u.Username,
		Role:                 int(u.Role),
		Nickname:             u.Nickname,
		AvatarURL:            u.AvatarURL,
		XPoints:              u.XPoints,
		RegistrationTime:     u.RegistrationTime,
		LastDailyCardsUpdate: u.LastDailyCardsUpdate,
	}
}

func toModelUser(u User) model.User {
	return model.User{
		CredentialsSecure: model.CredentialsSecure{
			ID:       strconv.Itoa(u.ID),
			Username: u.Username,
			Role:     model.Role(u.Role),
		},
		Nickname:             u.Nickname,
		AvatarURL:            u.AvatarURL,
		XPoints:              u.XPoints,
		RegistrationTime:     u.RegistrationTime,
		LastDailyCardsUpdate: u.LastDailyCardsUpdate,
	}
}

func toModelUsers(u []User) []model.User {
	users := make([]model.User, len(u))
	for i := range users {
		users[i] = toModelUser(u[i])
	}
	return users
}
