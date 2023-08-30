package service

import (
	"context"
	"trainee/internal/core"
	"trainee/internal/core/errors"
)

type SegmentRepository interface {
	Save(ctx context.Context, segment *core.Segment) (*core.Segment, error)
	Remove(ctx context.Context, segment *core.Segment) error
	GetCountOfActiveSegments(ctx context.Context, segment *core.Segment) (int, error)
}

type SegmentService struct {
	segmentRepository SegmentRepository
}

func NewSegmentService(segmentRepository SegmentRepository) *SegmentService {
	return &SegmentService{
		segmentRepository: segmentRepository,
	}
}

func (service *SegmentService) Create(ctx context.Context, segment *core.Segment) (*core.Segment, error) {
	return service.segmentRepository.Save(ctx, segment)
}

func (service *SegmentService) Delete(ctx context.Context, segment *core.Segment) error {
	count, err := service.segmentRepository.GetCountOfActiveSegments(ctx, segment)
	if err != nil {
		return err
	}
	if count == 0 {
		return &errors.NotFoundError{
			Type: "Segment",
			Id:   segment.Slug,
		}
	}
	return service.segmentRepository.Remove(ctx, segment)
}
