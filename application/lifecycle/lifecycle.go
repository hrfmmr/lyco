package lifecycle

import (
	"github.com/hrfmmr/lyco/domain/event"
	"github.com/sirupsen/logrus"
)

type LifecycleEventHub struct{}

func NewLifecycleEventHub() *LifecycleEventHub {
	return &LifecycleEventHub{}
}

func (h *LifecycleEventHub) EventType() event.EventType {
	return event.EventTypeAny
}

func (h *LifecycleEventHub) HandleEvent(e event.Event) {
	logrus.Debugf("👀 occurred event:%v", e)
}
