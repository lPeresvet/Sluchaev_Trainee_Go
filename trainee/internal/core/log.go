package core

import (
	"fmt"
	"strconv"
	"strings"
)

type Log struct {
	UserId    int64
	Slug      string
	Operation int
	Day       int
	Hour      int
	Minute    int
	Second    float32
}

func (log *Log) ToCSVString() string {
	var operationStr string
	if log.Operation == 0 {
		operationStr = "ADD"
	} else {
		operationStr = "DELETE"
	}
	return fmt.Sprintf("%v,%v,%v,%v %v:%v:%v",
		log.UserId, log.Slug, operationStr, log.Day, log.Hour, log.Minute, log.Second)
}

type LogRequest struct {
	UserId int64
	Month  int
	Year   int
}

type LogRequestDto struct {
	UserId    int64  `json:"id"`
	YearMonth string `json:"period"`
}

func (log *LogRequestDto) FromDto() (*LogRequest, error) {
	data := strings.Split(log.YearMonth, "-")

	y, err := strconv.Atoi(data[0])
	if err != nil {
		return nil, err
	}

	m, err := strconv.Atoi(data[1])
	if err != nil {
		return nil, err
	}

	return &LogRequest{
		UserId: log.UserId,
		Month:  m,
		Year:   y,
	}, nil
}

type LogResponse struct {
	UserId int64  `json:"id"`
	Link   string `json:"cssFile"`
}
