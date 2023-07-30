package sql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/wintermonth2298/xp-loyalty/internal/model"
)

type CredentialsRepository struct {
	db *sqlx.DB
}

func NewCredentialsRepository(db *sqlx.DB) *CredentialsRepository {
	return &CredentialsRepository{db: db}
}

func (r *CredentialsRepository) Create(ctx context.Context, creds model.Credentials) error {
	c := toSQLCreds(creds)

	query, args, err := psql.Insert("creds").Columns("username", "role", "password").
		Values(c.Username, c.Role, c.Password).ToSql()
	if err != nil {
		return fmt.Errorf("credsRepo - Create() - sq: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("credsRepo - Create() - ExecContext(): %w", err)
	}

	return nil
}

func (r *CredentialsRepository) GetByUsername(ctx context.Context, username string) (model.Credentials, error) {
	return r.get(ctx, sq.Eq{"creds.username": username})
}

func (r *CredentialsRepository) GetByCredentials(ctx context.Context, username, password string) (model.Credentials, error) {
	return r.get(ctx, sq.Eq{"creds.username": username, "creds.password": password})
}

func (r *CredentialsRepository) get(ctx context.Context, eq sq.Eq) (model.Credentials, error) {
	var creds Credentials

	result, err := get(ctx, r.db, "creds", eq, creds)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = model.ErrUserNotFound
		}
		return model.Credentials{}, fmt.Errorf("credsRepo - get(): %w", err)
	}
	return toModelCreds(result), nil
}

type Credentials struct {
	ID       int    `db:"id"`
	Username string `db:"username"`
	Role     int    `db:"role"`
	Password string `db:"password"`
}

func toSQLCreds(c model.Credentials) Credentials {
	id, _ := strconv.Atoi(c.ID)
	return Credentials{
		ID:       id,
		Username: c.Username,
		Role:     int(c.Role),
		Password: c.Password,
	}
}

func toModelCreds(c Credentials) model.Credentials {
	return model.Credentials{
		CredentialsSecure: model.CredentialsSecure{
			ID:       strconv.Itoa(c.ID),
			Username: c.Username,
			Role:     model.Role(c.Role),
		},
		Password: c.Password,
	}
}
