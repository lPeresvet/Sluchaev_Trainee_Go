package postgres

import (
	"context"
	"trainee/internal/core"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	INSERT_SEGMENT = `
	INSERT INTO segments (slug) 
	VALUES ($1);
	`
	DELETE_SEGMENT_BY_ID = `
	UPDATE segments 
	SET status = 1 
	WHERE slug = $1;
	`
	CHECK_SEGMENT = `
	SELECT count(id) 
	FROM segments 
	WHERE slug = $1 and status = 0;
	`
	GET_INFO_FOR_LOG = `
	SELECT us.user_id, us.segment_id
	FROM (SELECT id 
		FROM segments
		WHERE slug = $1) as "s"
		INNER JOIN user_segment us 
			ON us.segment_id = s.id;
	`
	ACTIVE  = 0
	DELETED = 1
)

type SegmentRepository struct {
	conn *pgxpool.Pool
}

func NewSegmentRepository(connection *pgxpool.Pool) *SegmentRepository {
	return &SegmentRepository{
		conn: connection,
	}
}

func (sr *SegmentRepository) Save(ctx context.Context, segment *core.Segment) (*core.Segment, error) {
	_, err := sr.conn.Exec(ctx, INSERT_SEGMENT, segment.Slug)
	if err != nil {
		return nil, err
	}
	return segment, nil
}

func (sr *SegmentRepository) Remove(ctx context.Context, segment *core.Segment) error {
	var info []*core.UserSegmentId

	rows, err := sr.conn.Query(ctx, GET_INFO_FOR_LOG, segment.Slug)

	if err != nil {
		return err
	}

	if err := pgxscan.ScanAll(&info, rows); err != nil {
		return err
	}

	tx, err := sr.conn.Begin(ctx)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, DELETE_SEGMENT_BY_ID, segment.Slug)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}
	for _, el := range info {
		_, err = tx.Exec(ctx, INSERT_LOG, el.UserId, el.SegmentId, DELETED)
		if err != nil {
			tx.Rollback(ctx)
			return err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}
	return nil
}

func (sr *SegmentRepository) GetCountOfActiveSegments(ctx context.Context, segment *core.Segment) (int, error) {
	var segmentsCount int
	err := sr.conn.QueryRow(ctx, CHECK_SEGMENT, segment.Slug).Scan(&segmentsCount)
	if err != nil {
		return 0, err
	}
	return segmentsCount, nil
}
