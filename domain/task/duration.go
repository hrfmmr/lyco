package task

import (
	"errors"
	"strconv"
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
	if len(strconv.FormatInt(value, 10)) != len(strconv.FormatInt(time.Now().UnixNano(), 10)) {
		return nil, errors.New("😕 [InvalidArgumentError] value must be nano scale")
	}
	return &taskDuration{value}, nil
}

func (s *taskDuration) Value() int64 {
	return s.value
}
