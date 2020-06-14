package timer

//go:generate stringer -type=TimerMode
type TimerMode uint

const (
	TimerModeTask TimerMode = iota
	TimerModeBreaks
)
