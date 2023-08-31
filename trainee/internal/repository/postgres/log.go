package postgres

import (
	"context"
	"trainee/internal/core"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	GET_LOG_SET = `
	SELECT op.user_id, s.slug, op.operation, 
		EXTRACT(DAY FROM op.operation_time) as "day", 
		EXTRACT(HOUR FROM op.operation_time) as "hour", 
		EXTRACT(MINUTE FROM op.operation_time) as "minute", 
		EXTRACT(SECOND FROM op.operation_time) as "second" 
	FROM (SELECT user_id, segment_id, operation, operation_time 
		FROM segment_log 
		WHERE EXTRACT(MONTH FROM operation_time) = $1 
			and EXTRACT(YEAR FROM operation_time) = $2 
			and user_id = $3) as "op" 
		INNER JOIN segments s 
			ON s.id = op.segment_id;
	`
)

type LogRepository struct {
	conn *pgxpool.Pool
}

func NewLogRepository(connection *pgxpool.Pool) *LogRepository {
	return &LogRepository{
		conn: connection,
	}
}

func (repository *LogRepository) GetByUserIdAndMonth(ctx context.Context, request *core.LogRequest) ([]*core.Log, error) {
	var logs []*core.Log

	rows, err := repository.conn.Query(ctx, GET_LOG_SET, request.Month, request.Year, request.UserId)

	if err != nil {
		return nil, err
	}

	if err := pgxscan.ScanAll(&logs, rows); err != nil {
		return nil, err
	}
	return logs, nil
}
