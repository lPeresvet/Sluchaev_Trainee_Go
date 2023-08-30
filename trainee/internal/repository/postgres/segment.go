package postgres

import (
	"context"
	"trainee/internal/core"

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
	_, err := sr.conn.Exec(ctx, DELETE_SEGMENT_BY_ID, segment.Slug)
	if err != nil {
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
