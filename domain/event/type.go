package event

type EventType int

const (
	EventTypeAny EventType = 1 << (iota + 1)
)
