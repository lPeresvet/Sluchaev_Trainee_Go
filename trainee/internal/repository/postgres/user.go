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
	chUser := int64(-1)
	err := ur.conn.QueryRow(ctx, CHECK_USER, userId).Scan(&chUser)
	if err != nil && err.Error() != "no rows in result set" {
		fmt.Println("CHECK_USER")
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
			_, err = ur.conn.Exec(ctx, INSERT_USER_SEGMENT, userId, segmentId)
			if err != nil {
				return nil, err
			}
		}
	}
	return segments, nil
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
		_, err = ur.conn.Exec(ctx, DELETE_USER_SEGMENT, userId, segmentId)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ur *UserRepository) GetById(ctx context.Context, userId int64) (*core.User, error) {
	chUser := int64(-1)
	err := ur.conn.QueryRow(ctx, CHECK_USER, userId).Scan(&chUser)
	if err != nil && err.Error() != "no rows in result set" {
		fmt.Println("CHECK_USER")
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
