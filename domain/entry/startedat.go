package entry

import (
	"errors"
	"strconv"
	"time"
)

type (
	StartedAt interface {
		Value() int64
	}
	entryStartedAt struct {
		value int64
	}
)

func NewStartedAt(value int64) (StartedAt, error) {
	if len(strconv.FormatInt(value, 10)) != len(strconv.FormatInt(time.Now().UnixNano(), 10)) {
		return nil, errors.New("😕 entry.StartedAt value must be nano scale")
	}
	return &entryStartedAt{value}, nil
}

func (s *entryStartedAt) Value() int64 {
	return s.value
}
