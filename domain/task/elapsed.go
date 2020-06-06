package task

import (
	"errors"
	"strconv"
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
	if len(strconv.FormatInt(value, 10)) != len(strconv.FormatInt(time.Now().UnixNano(), 10)) {
		return nil, errors.New("😕 [InvalidArgumentError] value must be nano scale")
	}
	return &taskElapsed{value}, nil
}

func (s *taskElapsed) Value() int64 {
	return s.value
}
