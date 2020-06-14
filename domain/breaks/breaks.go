package breaks

import (
	"time"

	"github.com/hrfmmr/lyco/domain/event"
)

const (
	DefaultShortDuration = 5 * time.Minute
	DefaultLongDuration  = 15 * time.Minute
)

type (
	Breaks interface {
		// props
		Duration() Duration
		StartedAt() StartedAt
		EndedAt() EndedAt
		// behaviors
		Start() error
		Stop() error
		// utils
		IsStopped() bool
	}

	breaks struct {
		duration  Duration
		startedAt StartedAt
		endedAt   EndedAt
	}
)

func NewBreaks(duration Duration) Breaks {
	return &breaks{duration: duration}
}

func ShortDefault() Breaks {
	d, _ := NewDuration(int64(DefaultShortDuration))
	return NewBreaks(d)
}

func LongDefault() Breaks {
	d, _ := NewDuration(int64(DefaultLongDuration))
	return NewBreaks(d)
}

func (b *breaks) Duration() Duration {
	return b.duration
}

func (b *breaks) StartedAt() StartedAt {
	return b.startedAt
}

func (b *breaks) EndedAt() EndedAt {
	return b.endedAt
}

func (b *breaks) Start() error {
	now := time.Now()
	startedAt, err := NewStartedAt(now.UnixNano())
	if err != nil {
		return err
	}
	b.startedAt = startedAt
	event.DefaultPublisher.Publish(NewBreaksStarted(
		b.startedAt,
		b.duration,
	))
	return nil
}

func (b *breaks) Stop() error {
	now := time.Now()
	endedAt, err := NewEndedAt(now.UnixNano())
	if err != nil {
		return err
	}
	b.endedAt = endedAt
	event.DefaultPublisher.Publish(NewBreaksEnded(
		b.startedAt,
		b.duration,
		b.endedAt,
	))
	return nil
}

func (b *breaks) IsStopped() bool {
	return b.endedAt != nil
}
