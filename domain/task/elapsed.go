package task

import (
	"errors"
	"time"
)

type (
	Elapsed interface {
		Value() int64
	}
	taskElapsed struct {
		value int64
	}
)

func NewElapsed(value int64) (Elapsed, error) {
	if value != value/int64(time.Nanosecond) {
		return nil, errors.New("😕 [InvalidArgumentError] task.Elapsed value must be nano scale")
	}
	return &taskElapsed{value}, nil
}

func (s *taskElapsed) Value() int64 {
	return s.value
}
