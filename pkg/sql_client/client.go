package sql_client

import (
	"context"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/wintermonth2298/xp-loyalty/pkg/config"
)

func New(ctx context.Context, cfg config.SQL) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database)

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	db, err := sqlx.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	return db, nil
}
