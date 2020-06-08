package task

import (
	"errors"
	"time"
)

type (
	Duration interface {
		Value() int64
	}
	taskDuration struct {
		value int64
	}
)

func NewDuration(value int64) (Duration, error) {
	if value != value/int64(time.Nanosecond) {
		return nil, errors.New("😕 [InvalidArgumentError] task.Duration value must be nano scale")
	}
	return &taskDuration{value}, nil
}

func (s *taskDuration) Value() int64 {
	return s.value
}
