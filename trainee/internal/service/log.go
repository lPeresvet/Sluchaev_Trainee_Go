package service

import (
	"context"
	"fmt"
	"os"
	"time"
	"trainee/internal/core"
	"trainee/internal/core/errors"
)

type LogRepository interface {
	GetByUserIdAndMonth(ctx context.Context, request *core.LogRequest) ([]*core.Log, error)
}

type LogService struct {
	logRepository  LogRepository
	userRepository UserRepository
}

func NewLogService(logRepositry LogRepository, userRepository UserRepository) *LogService {
	return &LogService{
		logRepository:  logRepositry,
		userRepository: userRepository,
	}
}

func (service *LogService) GetByUserIdAndMonth(ctx context.Context, request *core.LogRequestDto) (*core.LogResponse, error) {
	count, err := service.userRepository.CountUser(ctx, request.UserId)
	if err != nil {
		return nil, err
	}

	if count == 0 {
		return nil, &errors.NotFoundError{
			Type: "User",
			Id:   fmt.Sprint(request.UserId),
		}
	}

	logRequest, err := request.FromDto()
	if err != nil {
		return nil, err
	}

	logs, err := service.logRepository.GetByUserIdAndMonth(ctx, logRequest)

	if err != nil {
		return nil, err
	}

	link, err := generateCSV(logs)

	if err != nil {
		return nil, err
	}
	return &core.LogResponse{
		UserId: request.UserId,
		Link:   link,
	}, nil
}

func generateCSV(logs []*core.Log) (string, error) {
	now := time.Now()
	name := fmt.Sprintf("/%v:%v:%v.csv", now.Hour(), now.Minute(), now.Second())
	file, err := os.Create("./csv" + name)
	defer file.Close()
	if err != nil {
		return "", err
	}

	for _, log := range logs {
		_, err := file.WriteString(log.ToCSVString() + "\n")
		if err != nil {
			return "", nil
		}
	}

	return name, nil
}
