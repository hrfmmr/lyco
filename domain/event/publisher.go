package event

type (
	Publisher interface {
		Publish(e Event)
		Subscribe(s Subscriber)
	}

	publisher struct {
		subscribers []Subscriber
	}
)

func NewPublisher() Publisher {
	return &publisher{
		[]Subscriber{},
	}
}

func (p *publisher) Publish(e Event) {
	for _, s := range p.subscribers {
		if e.Type()&(EventTypeAny|s.EventType()) > 0 {
			s.HandleEvent(e)
		}
	}
}

func (p *publisher) Subscribe(s Subscriber) {
	p.subscribers = append(p.subscribers, s)
}
