package sql

import (
	"context"
	"fmt"
	"strconv"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/wintermonth2298/xp-loyalty/internal/model"
)

type AdminsRepository struct {
	db *sqlx.DB
}

func NewAdminsRepository(db *sqlx.DB) *AdminsRepository {
	return &AdminsRepository{db: db}
}

func (r *AdminsRepository) Create(ctx context.Context, admin model.Admin) error {
	u := toSQLAdmin(admin)

	query, args, err := psql.Insert("admin").
		Columns("username",
			"role",
		).
		Values(u.Username,
			u.Role,
		).ToSql()
	if err != nil {
		return fmt.Errorf("adminsRepo - Create() - sq: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("adminsRepo - Create() - ExecContext(): %w", err)
	}

	return nil
}

func (r *AdminsRepository) GetByUsername(ctx context.Context, username string) (model.Admin, error) {
	var admin Admin
	result, err := get(ctx, r.db, "admin", sq.Eq{"admin.username": username}, admin)
	if err != nil {
		return model.Admin{}, fmt.Errorf("adminRepo - GetByUsername(): %w", err)
	}
	return toModelAdmin(result), nil
}

type Admin struct {
	ID       int    `db:"id"`
	Username string `db:"username"`
	Role     int    `db:"role"`
}

func toSQLAdmin(u model.Admin) Admin {
	id, _ := strconv.Atoi(u.ID)
	return Admin{
		ID:       id,
		Username: u.Username,
		Role:     int(u.Role),
	}
}

func toModelAdmin(u Admin) model.Admin {
	return model.Admin{
		CredentialsSecure: model.CredentialsSecure{
			ID:       strconv.Itoa(u.ID),
			Username: u.Username,
			Role:     model.Role(u.Role),
		},
	}
}
