package timer

import (
	"context"
	"time"

	"github.com/hrfmmr/lyco/domain/event"
)

const (
	defaultTickInterval = time.Second
)

type (
	Timer interface {
		Start(duration Duration)
		Stop()
	}

	timer struct {
		tickinterval time.Duration
		ctx          context.Context
		cancel       func()
	}
)

func NewTimer() Timer {
	ctx, cancel := context.WithCancel(context.Background())
	return &timer{
		defaultTickInterval,
		ctx,
		cancel,
	}
}

func (t *timer) Start(duration Duration) {
	t.ensureContextInitialized()
	go func(d time.Duration) {
		ticker := time.NewTicker(t.tickinterval)
		for r := d; r > 0; r -= t.tickinterval {
			select {
			case <-t.ctx.Done():
				return
			default:
				event.DefaultPublisher.Publish(NewTimerTicked())
				<-ticker.C
			}
		}
		t.cancel()
		event.DefaultPublisher.Publish(NewTimerFinished())
	}(time.Duration(duration.Value()))
}

func (t *timer) Stop() {
	t.cancel()
}

func (t *timer) ensureContextInitialized() {
	if t.ctx.Err() == nil {
		t.cancel()
	}
	ctx, cancel := context.WithCancel(context.Background())
	t.ctx = ctx
	t.cancel = cancel
}
