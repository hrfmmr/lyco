package breaks

import "time"

const (
	DefaultShortDuration = 5 * time.Minute
	DefaultLongDuration  = 15 * time.Minute
)

type (
	Breaks interface {
		// props
		Duration() time.Duration
		StartedAt() StartedAt
		EndedAt() EndedAt
		// behaviors
		Start() error
		Stop() error
		// utils
		IsStopped() bool
	}

	breaks struct {
		duration  time.Duration
		startedAt StartedAt
		endedAt   EndedAt
	}
)

func NewBreaks(duration time.Duration) Breaks {
	return &breaks{duration: duration}
}

func ShortDefault() Breaks {
	return NewBreaks(DefaultShortDuration)
}

func LongDefault() Breaks {
	return NewBreaks(DefaultLongDuration)
}

func (b *breaks) Duration() time.Duration {
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
	return nil
}

func (b *breaks) Stop() error {
	now := time.Now()
	endedAt, err := NewEndedAt(now.UnixNano())
	if err != nil {
		return err
	}
	b.endedAt = endedAt
	return nil
}

func (b *breaks) IsStopped() bool {
	return b.endedAt != nil
}
