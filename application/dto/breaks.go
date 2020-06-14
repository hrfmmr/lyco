package dto

import (
	"fmt"
	"time"

	"github.com/hrfmmr/lyco/domain/breaks"
)

type (
	BreaksDTO interface {
		RemainsTimerText() string
	}

	breaksDTO struct {
		duration  int64
		startedAt *int64
		endedAt   *int64
	}
)

func ConvertBreaksToDTO(b breaks.Breaks) BreaksDTO {
	var startedAt *int64
	if b.StartedAt() != nil {
		val := b.StartedAt().Value()
		startedAt = &val
	}
	var endedAt *int64
	if b.EndedAt() != nil {
		val := b.EndedAt().Value()
		endedAt = &val
	}
	return &breaksDTO{
		b.Duration().Value(),
		startedAt,
		endedAt,
	}
}

func (b *breaksDTO) RemainsTimerText() string {
	rsec := b.remainsDuration() / 1e9
	return fmt.Sprintf("%02d:%02d", int(rsec/60)%60, rsec%60)
}

func (b *breaksDTO) remainsDuration() int64 {
	to := *b.startedAt + b.duration
	now := time.Now().UnixNano()
	return to - now
}
