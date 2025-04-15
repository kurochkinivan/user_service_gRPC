package pg

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/kurochkinivan/user_service/internal/entity"
	"github.com/kurochkinivan/user_service/pkg/pgerr"
)

func (s *Storage) DeleteInterests(ctx context.Context, userID string, tx pgx.Tx) error {
	const op = "storage.pg.DeleteInterests"

	sql, args, err := s.qb.
		Delete(TableUserInterests).
		Where(sq.Eq{"user_id": userID}).
		ToSql()
	if err != nil {
		return pgerr.ErrCreateQuery(op, err)
	}

	execFn := s.pool.Exec
	if tx != nil {
		execFn = tx.Exec
	}

	_, err = execFn(ctx, sql, args...)
	if err != nil {
		return pgerr.ErrExec(op, err)
	}

	return nil
}

func (s *Storage) CreateInterests(ctx context.Context, userID string, interests []*entity.Interest, tx pgx.Tx) error {
	const op = "storage.pg.UpdateInterests"

	copyFromFn := s.pool.CopyFrom
	if tx != nil {
		copyFromFn = tx.CopyFrom
	}

	_, err := copyFromFn(
		ctx,
		pgx.Identifier{TableUserInterests},
		[]string{"user_id", "interest_id"},
		pgx.CopyFromSlice(len(interests), func(i int) ([]any, error) {
			return []any{userID, interests[i].ID}, nil
		}),
	)
	if err != nil {
		return pgerr.ErrInsertMultipleRows(op, err)
	}

	return nil
}
