package entry

import (
	"errors"
	"time"
)

type (
	Elapsed interface {
		Value() int64
		Add(elapsed Elapsed)
	}
	entryElapsed struct {
		value int64
	}
)

func NewElapsed(value int64) (Elapsed, error) {
	if value != value/int64(time.Nanosecond) {
		return nil, errors.New("😕 entry.Elapsed value must be nano scale")
	}
	return &entryElapsed{value}, nil
}

func (e *entryElapsed) Value() int64 {
	return e.value
}

func (e *entryElapsed) Add(elapsed Elapsed) {
	e.value += elapsed.Value()
}
