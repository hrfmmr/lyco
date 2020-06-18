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
		Start(mode TimerMode, duration Duration)
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

func (t *timer) Start(mode TimerMode, duration Duration) {
	t.ensureContextInitialized()
	go func(d time.Duration) {
		ticker := time.NewTicker(t.tickinterval)
		remains := d
		event.DefaultPublisher.Publish(NewTimerTicked(mode))
	Loop:
		for {
			select {
			case <-t.ctx.Done():
				return
			case <-ticker.C:
				remains -= t.tickinterval
				if remains <= 0 {
					break Loop
				}
				event.DefaultPublisher.Publish(NewTimerTicked(mode))
			}
		}
		t.cancel()
		event.DefaultPublisher.Publish(NewTimerFinished(mode))
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
