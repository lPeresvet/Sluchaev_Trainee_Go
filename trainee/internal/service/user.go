package service

import (
	"context"
	"trainee/internal/core"
)

type UserRepository interface {
	Save(ctx context.Context, userId int64) (*core.User, error)
	AddSegments(ctx context.Context, userId int64, segments []*core.Segment) ([]*core.Segment, error)
	DeleteSegments(ctx context.Context, userId int64, segments []*core.Segment) (error)
	GetById(ctx context.Context, userId int64) (*core.User, error)
	CountUser(ctx context.Context, userId int64) (int, error)
}

type UserService struct {
	userRepository UserRepository
}

func NewUserService(userRepository UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (us *UserService) GetById(ctx context.Context, userId int64) (*core.User, error) {
	return us.userRepository.GetById(ctx, userId)
}

func (us *UserService) ProcessUserData(ctx context.Context, userRequest *core.UserRequest) (*core.User, error) {
	_, err := us.userRepository.Save(ctx, userRequest.Id)
	if err != nil {
		return nil, err
	}
	_, err = us.userRepository.AddSegments(ctx, userRequest.Id, userRequest.SegmentsToAdd)
	if err != nil {
		return nil, err
	}
	err = us.userRepository.DeleteSegments(ctx, userRequest.Id, userRequest.SegmentsToDelete)
	if err != nil {
		return nil, err
	}
	return us.userRepository.GetById(ctx, userRequest.Id)
}
