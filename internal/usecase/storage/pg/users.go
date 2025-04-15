package pg

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kurochkinivan/user_service/internal/entity"
	"github.com/kurochkinivan/user_service/pkg/pgerr"
)

type Storage struct {
	pool *pgxpool.Pool
	qb   sq.StatementBuilderType
}

func New(pool *pgxpool.Pool) *Storage {
	return &Storage{
		pool: pool,
		qb:   sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *Storage) User(ctx context.Context, userID string, tx pgx.Tx) (*entity.User, error) {
	const op = "storage.pg.User"

	SQL, args, err := r.qb.
		Select(
			"u.id",
			"u.name",
			"u.age",
			"u.gender",
			"u.about",
			"p.id AS photo_id",
			"p.photo_url AS photo_url",
			"p.created_at AS photo_created_at",
			"i.id AS interest_id",
			"i.name AS interest_name",
		).
		From(TableUsers + " AS u").
		LeftJoin(TablePhotos + " AS p ON u.id = p.user_id").
		LeftJoin(TableUserInterests + " AS ui ON u.id = ui.user_id").
		LeftJoin(TableInterests + " AS i ON ui.interest_id = i.id").
		Where(
			sq.Eq{"u.id": userID},
		).
		ToSql()
	if err != nil {
		return nil, pgerr.ErrCreateQuery(op, err)
	}

	queryFn := r.pool.Query
	if tx != nil {
		queryFn = tx.Query
	}

	rows, err := queryFn(ctx, SQL, args...)
	if err != nil {
		return nil, pgerr.ErrDoQuery(op, err)
	}
	defer rows.Close()

	photoSeen := make(map[int64]bool)
	interestSeen := make(map[int64]bool)

	user := new(entity.User)
	for rows.Next() {
		var photoID, interestID sql.NullInt64
		var photoUrl, interestName sql.NullString
		var photoCreatedAt sql.NullTime

		rows.Scan(
			&user.ID,
			&user.Name,
			&user.Age,
			&user.Gender,
			&user.About,
			&photoID,
			&photoUrl,
			&photoCreatedAt,
			&interestID,
			&interestName,
		)

		if photoID.Valid && !photoSeen[photoID.Int64] {
			photoSeen[photoID.Int64] = true
			user.Photos = append(user.Photos, &entity.Photo{
				ID:        photoID.Int64,
				Url:       photoUrl.String,
				CreatedAt: photoCreatedAt.Time,
			})
		}

		if interestID.Valid && !interestSeen[interestID.Int64] {
			interestSeen[interestID.Int64] = true
			user.Interests = append(user.Interests, &entity.Interest{
				ID:   interestID.Int64,
				Name: interestName.String,
			})
		}
	}

	if err := rows.Err(); err != nil {
		return nil, pgerr.ErrScan(op, err)
	}

	return user, nil
}

func (r *Storage) UpdateUser(ctx context.Context, userID string, user *entity.User, tx pgx.Tx) error {
	const op = "Storage.pg.UpdateUser"

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

	execFn := r.pool.Exec
	if tx != nil {
		execFn = tx.Exec
	}

	_, err = execFn(ctx, sql, args...)
	if err != nil {
		return pgerr.ErrExec(op, err)
	}

	return nil
}
