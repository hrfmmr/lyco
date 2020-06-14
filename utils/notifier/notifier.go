package notifier

import (
	"fmt"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/hrfmmr/lyco/domain/breaks"
	"github.com/hrfmmr/lyco/domain/task"
	"github.com/sirupsen/logrus"
)

type (
	Notifier interface {
		Notify(title, message string) error
	}

	beeepNotifier struct{}
)

func (n *beeepNotifier) Notify(title, message string) error {
	if err := beeep.Notify(title, message, ""); err != nil {
		return err
	}
	return nil
}

func New() Notifier {
	return &beeepNotifier{}
}

func NotifyForBreaksStart(n Notifier, b breaks.Breaks) {
	till := time.Unix(0, b.StartedAt().Value()+b.Duration().Value())
	d := time.Duration(b.Duration().Value())
	n.Notify(
		"lyco",
		fmt.Sprintf("☕Take a %v break till %v", d, till.Format("15:04")),
	)
}

func NotifyForBreaksEnd(n Notifier, t task.Task) {
	if t.StartedAt() == nil {
		logrus.Errorf("❗[NotifyForBreaksEnd] startedAt is nil for task:%v", t)
		return
	}
	till := time.Unix(0, t.StartedAt().Value()+t.Duration().Value())
	d := time.Duration(t.Duration().Value())
	n.Notify(
		"lyco",
		fmt.Sprintf("🔨Work for %v till %v on %s", d, till.Format("15:04"), t.Name().Value()),
	)
}
