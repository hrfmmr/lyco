package breaks

import (
	"errors"
	"strconv"
	"time"
)

type (
	EndedAt interface {
		Value() int64
	}

	endedAt struct {
		value int64
	}
)

func NewEndedAt(value int64) (EndedAt, error) {
	if len(strconv.FormatInt(value, 10)) != len(strconv.FormatInt(time.Now().UnixNano(), 10)) {
		return nil, errors.New("😕 [InvalidArgumentError] value must be nano scale")
	}
	return &endedAt{value}, nil
}

func (s *endedAt) Value() int64 {
	return s.value
}
