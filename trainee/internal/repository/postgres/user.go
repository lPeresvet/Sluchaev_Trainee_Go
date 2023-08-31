package postgres

import (
	"context"
	"fmt"
	"trainee/internal/core"
	"trainee/internal/core/errors"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	INSERT_USER_SEGMENT = `
	INSERT INTO user_segment (user_id, segment_id) 
	VALUES ($1, $2);
	`
	INSERT_USER = `
	INSERT INTO accounts (id) 
	VALUES ($1);
	`
	CHECK_USER_SEGMENTS = `
	SELECT count(user_id) 
	FROM user_segment 
	WHERE user_id = $1 and segment_id = $2;
	`
	CHECK_USER = `
	SELECT count(id) 
	FROM accounts 
	WHERE id = $1;
	`
	DELETE_USER_SEGMENT = `
	DELETE FROM user_segment 
	WHERE user_id = $1 and segment_id = $2;
	`
	GET_SEGMENT_ID = `
	SELECT id 
	FROM segments 
	WHERE slug = $1 and status = 0;
	`
	GET_USER_SEGMENTS = `
	SELECT s.slug 
	FROM (SELECT * 
		FROM user_segment 
		WHERE user_id = $1) as "us" 
		INNER JOIN segments s 
			ON us.segment_id = s.id 
		WHERE s.status = 0;
	`
	INSERT_LOG = `
	INSERT INTO segment_log (user_id, segment_id, operation, operation_time) 
	VALUES ($1, $2, $3, now());
	`
)

type UserRepository struct {
	conn *pgxpool.Pool
}

func NewUserRepository(connection *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		conn: connection,
	}
}

func (ur *UserRepository) Save(ctx context.Context, userId int64) (*core.User, error) {
	chUser, err := ur.CountUser(ctx, userId)

	if err != nil {
		return nil, err
	}

	if chUser == 0 {
		_, err = ur.conn.Exec(ctx, INSERT_USER, userId)
		if err != nil {
			return nil, err
		}
	}
	return &core.User{
		Id: userId,
	}, nil
}

func (ur *UserRepository) AddSegments(ctx context.Context, userId int64,
	segments []*core.Segment) ([]*core.Segment, error) {
	var segmentId int64
	var segmentsCount int
	for i := 0; i < len(segments); i++ {
		err := ur.conn.QueryRow(ctx, GET_SEGMENT_ID, segments[i].Slug).Scan(&segmentId)

		if err != nil {
			if err.Error() == "no rows in result set" {
				return nil, &errors.NotFoundErrorWithMessage{
					Message: "Can not add segment with slug <" + segments[i].Slug + ">",
				}
			}
			return nil, err
		}

		err = ur.conn.QueryRow(ctx, CHECK_USER_SEGMENTS, userId, segmentId).Scan(&segmentsCount)

		if err != nil {
			return nil, err
		}

		if segmentsCount == 0 {
			if err = ur.insertSegment(ctx, userId, segmentId); err != nil {
				return nil, err
			}
		}
	}
	return segments, nil
}

func (ur *UserRepository) insertSegment(ctx context.Context, userId int64, segmentId int64) error {
	tx, err := ur.conn.Begin(ctx)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, INSERT_USER_SEGMENT, userId, segmentId)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	_, err = tx.Exec(ctx, INSERT_LOG, userId, segmentId, ACTIVE)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}

func (ur *UserRepository) DeleteSegments(ctx context.Context, userId int64,
	segments []*core.Segment) error {

	var segmentId int64

	for i := 0; i < len(segments); i++ {

		err := ur.conn.QueryRow(ctx, GET_SEGMENT_ID, segments[i].Slug).Scan(&segmentId)

		if err != nil {
			if err.Error() == "no rows in result set" {
				return &errors.NotFoundErrorWithMessage{
					Message: "Can not delete segment with slug <" + segments[i].Slug + ">",
				}
			}
			return err
		}

		if err = ur.deleteSegment(ctx, userId, segmentId); err != nil {
			return err
		}

	}
	return nil
}

func (ur *UserRepository) deleteSegment(ctx context.Context, userId int64, segmentId int64) error {
	tx, err := ur.conn.Begin(ctx)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, DELETE_USER_SEGMENT, userId, segmentId)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	_, err = tx.Exec(ctx, INSERT_LOG, userId, segmentId, DELETED)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return nil
}

func (ur *UserRepository) GetById(ctx context.Context, userId int64) (*core.User, error) {
	chUser, err := ur.CountUser(ctx, userId)

	if err != nil {
		return nil, err
	}

	if chUser == 0 {
		return nil, &errors.NotFoundErrorWithMessage{
			Message: fmt.Sprintf("No user with id <%v> found", userId),
		}
	}

	var segments []*core.Segment

	rows, err := ur.conn.Query(ctx, GET_USER_SEGMENTS, userId)

	if err != nil {
		return nil, err
	}

	if err := pgxscan.ScanAll(&segments, rows); err != nil {
		return nil, err
	}

	return &core.User{
		Id:       userId,
		Segments: segments,
	}, nil
}

func (ur *UserRepository) CountUser(ctx context.Context, userId int64) (int, error) {
	chUser := 0

	err := ur.conn.QueryRow(ctx, CHECK_USER, userId).Scan(&chUser)

	if err != nil && err.Error() != "no rows in result set" {
		return 0, err
	}

	return chUser, nil
}
