package event

type Subscriber interface {
	EventType() EventType
	HandleEvent(e Event)
}
