package timer

import (
	"context"
	"time"

	"github.com/hrfmmr/lyco/domain/task"
)

type TaskTimer interface {
	Start(task.Task)
	Stop()
	Ticker() <-chan task.Task
	OnFinished() <-chan struct{}
}

type timer struct {
	tickCh chan task.Task
	finCh  chan struct{}
	ctx    context.Context
	cancel func()
}

func NewTaskTimer() TaskTimer {
	ctx, cancel := context.WithCancel(context.Background())
	return &timer{
		tickCh: make(chan task.Task, 1),
		finCh:  make(chan struct{}, 1),
		ctx:    ctx,
		cancel: cancel,
	}
}

func (t *timer) Start(m task.Task) {
	go func(d time.Duration) {
		tick := time.Second
		ticker := time.NewTicker(tick)
		for r := d; r >= 0; r -= tick {
			select {
			case <-t.ctx.Done():
				return
			default:
				t.tickCh <- m
				<-ticker.C
			}
		}
		t.finCh <- struct{}{}
		t.cancel()
	}(m.Duration())
}

func (t *timer) Stop() {
	t.cancel()
}

func (t *timer) Ticker() <-chan task.Task {
	return t.tickCh
}

func (t *timer) OnFinished() <-chan struct{} {
	return t.finCh
}
