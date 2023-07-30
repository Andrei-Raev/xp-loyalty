package sql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/wintermonth2298/xp-loyalty/internal/model"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

func getAll[T any](ctx context.Context, db *sqlx.DB, table string, eq sq.Eq, dest []T) ([]T, error) {
	query, args, err := psql.Select("*").From(table).ToSql()
	if err != nil {
		return nil, fmt.Errorf("common - getAll() - sq: %w", err)
	}

	err = db.SelectContext(ctx, &dest, query, args...)
	if err != nil {
		return nil, fmt.Errorf("common - getAll() - SelectContext(): %w", err)
	}
	return dest, nil
}

func get[T any](ctx context.Context, db *sqlx.DB, table string, eq sq.Eq, dest T) (T, error) {
	var empty T
	query, args, err := psql.Select("*").From(table).
		Where(eq).ToSql()
	if err != nil {
		return empty, fmt.Errorf("common - getAll() - sq: %w", err)
	}

	err = db.QueryRowxContext(ctx, query, args...).StructScan(&dest)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = model.ErrUserNotFound
		}
		return empty, fmt.Errorf("common - getAll() - QueryRowxContext(): %w", err)
	}
	return dest, nil
}
