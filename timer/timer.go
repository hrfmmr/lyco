package timer

import "time"

type Timer interface {
	Start(tick, d time.Duration) <-chan time.Duration
}

type timer struct{}

func New() Timer {
	return &timer{}
}

func (t *timer) Start(tick, d time.Duration) <-chan time.Duration {
	ch := make(chan time.Duration, 1)
	go func(ch chan time.Duration) {
		ticker := time.NewTicker(tick)
		for r := d; r >= 0; r -= tick {
			ch <- r
			<-ticker.C
		}
		close(ch)
	}(ch)
	return ch
}
