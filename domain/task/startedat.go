package task

import (
	"errors"
	"strconv"
	"time"
)

type (
	StartedAt interface {
		Value() int64
	}
	taskStartedAt struct {
		value int64
	}
)

func NewStartedAt(value int64) (StartedAt, error) {
	if len(strconv.FormatInt(value, 10)) != len(strconv.FormatInt(time.Now().UnixNano(), 10)) {
		return nil, errors.New("😕 [InvalidArgumentError] task.StartedAt value must be nano scale")
	}
	return &taskStartedAt{value}, nil
}

func (s *taskStartedAt) Value() int64 {
	return s.value
}
