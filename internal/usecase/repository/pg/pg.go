package pg

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kurochkinivan/user_service/internal/entity"
	"github.com/kurochkinivan/user_service/pkg/pgerr"
)

type Repository struct {
	pool *pgxpool.Pool
	qb   sq.StatementBuilderType
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{
		pool: pool,
		qb:   sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

const (
	TableUsers = "users"
)

func (r *Repository) Update(ctx context.Context, userID string, user *entity.User) error {
	const op = "repository.pg.Update"

	sql, args, err := r.qb.
		Insert(TableUsers).
		Columns(
			"id",
			"name",
			"age",
			"gender",
			"about",
		).
		Values(
			userID,
			user.Name,
			user.Age,
			user.Gender,
			user.About,
		).
		Suffix(`
		ON CONFLICT (id) DO UPDATE SET
			name   = EXCLUDED.name,
			age    = EXCLUDED.age,
			gender = EXCLUDED.gender,
			about  = EXCLUDED.about
	`).ToSql()
	if err != nil {
		return pgerr.ErrCreateQuery(op, err)
	}

	_, err = r.pool.Exec(ctx, sql, args...)
	if err != nil {
		return pgerr.ErrExec(op, err)
	}

	return nil
}
