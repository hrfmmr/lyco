package timer

import (
	"errors"
	"time"
)

type (
	Duration interface {
		Value() int64
	}
	timerDuration struct {
		value int64
	}
)

func NewDuration(value int64) (Duration, error) {
	if value != value/int64(time.Nanosecond) {
		return nil, errors.New("😕 [InvalidArgumentError] timerDuration value must be nano scale")
	}
	return &timerDuration{value}, nil
}

func (s *timerDuration) Value() int64 {
	return s.value
}
