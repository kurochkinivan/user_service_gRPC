package pg

import (
	"context"
	"fmt"

	"github.com/kurochkinivan/user_service/internal/entity"
	"github.com/kurochkinivan/user_service/pkg/pgerr"
)

func (r *Storage) UpdateProfile(ctx context.Context, userID string, user *entity.User) (*entity.User, error) {
	const op = "storage.pg.UpdateProfile"

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, pgerr.ErrCreateTx(op, err)
	}
	defer tx.Rollback(ctx)

	err = r.UpdateUser(ctx, userID, user, tx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = r.DeleteInterests(ctx, userID, tx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = r.CreateInterests(ctx, userID, user.Interests, tx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	user, err = r.User(ctx, userID, tx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, pgerr.ErrCommit(op, err)
	}

	return user, nil
}
