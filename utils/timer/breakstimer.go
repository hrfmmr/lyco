package timer

import (
	"context"
	"time"

	"github.com/hrfmmr/lyco/domain/breaks"
)

type (
	BreaksTimer interface {
		Start(breaks.Breaks)
		Breaks() breaks.Breaks
		Stop()
		Ticker() <-chan breaks.Breaks
		OnFinished() <-chan struct{}
	}

	breaksTimer struct {
		b      breaks.Breaks
		tickCh chan breaks.Breaks
		finCh  chan struct{}
		ctx    context.Context
		cancel func()
	}
)

func NewBreaksTimer() BreaksTimer {
	ctx, cancel := context.WithCancel(context.Background())
	return &breaksTimer{
		tickCh: make(chan breaks.Breaks, 1),
		finCh:  make(chan struct{}, 1),
		ctx:    ctx,
		cancel: cancel,
	}
}

func (t *breaksTimer) Start(b breaks.Breaks) {
	t.b = b
	go func(d time.Duration) {
		tick := time.Second
		ticker := time.NewTicker(tick)
		for r := d; r >= 0; r -= tick {
			select {
			case <-t.ctx.Done():
				return
			default:
				t.tickCh <- b
				<-ticker.C
			}
		}
		t.finCh <- struct{}{}
		t.cancel()
	}(b.Duration())
}

func (t *breaksTimer) Breaks() breaks.Breaks {
	return t.b
}

func (t *breaksTimer) Stop() {
	t.cancel()
}

func (t *breaksTimer) Ticker() <-chan breaks.Breaks {
	return t.tickCh
}

func (t *breaksTimer) OnFinished() <-chan struct{} {
	return t.finCh
}
