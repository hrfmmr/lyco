package breaks

import (
	"errors"
	"strconv"
	"time"
)

type (
	StartedAt interface {
		Value() int64
	}

	startedAt struct {
		value int64
	}
)

func NewStartedAt(value int64) (StartedAt, error) {
	if len(strconv.FormatInt(value, 10)) != len(strconv.FormatInt(time.Now().UnixNano(), 10)) {
		return nil, errors.New("😕 [InvalidArgumentError] value must be nano scale")
	}
	return &startedAt{value}, nil
}

func (s *startedAt) Value() int64 {
	return s.value
}
