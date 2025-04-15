package user

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/kurochkinivan/user_service/internal/entity"
	"github.com/kurochkinivan/user_service/internal/lib/sl"
)

type User struct {
	log         *slog.Logger
	userUpdator UserUpdator
}

func New(log *slog.Logger, userUpdator UserUpdator) *User {
	return &User{
		log:         log,
		userUpdator: userUpdator,
	}
}

type UserUpdator interface {
	UpdateProfile(ctx context.Context, userID string, user *entity.User) (*entity.User, error)
}

func (u *User) UpdateProfile(ctx context.Context, userID string, user *entity.User) (*entity.User, error) {
	const op = "user.UpdateProfile"
	log := u.log.With(
		slog.String("op", op),
		slog.String("user_id", userID),
	)

	log.Info("updating user")

	user, err := u.userUpdator.UpdateProfile(ctx, userID, user)
	if err != nil {
		log.Warn("failed to update user", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}
