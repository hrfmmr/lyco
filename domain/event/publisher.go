package event

type (
	Publisher interface {
		Publish(e Event)
		Subscribe(subscribers ...Subscriber)
	}

	publisher struct {
		subscribers []Subscriber
	}
)

var (
	DefaultPublisher = NewPublisher()
)

func NewPublisher() Publisher {
	return &publisher{
		[]Subscriber{},
	}
}

func (p *publisher) Publish(e Event) {
	for _, s := range p.subscribers {
		if s.EventType()&(EventTypeAny|e.Type()) > 0 {
			s.HandleEvent(e)
		}
	}
}

func (p *publisher) Subscribe(subscribers ...Subscriber) {
	p.subscribers = subscribers
}
