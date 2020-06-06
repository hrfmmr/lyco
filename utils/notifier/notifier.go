package notifier

import (
	"fmt"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/hrfmmr/lyco/domain/breaks"
	"github.com/hrfmmr/lyco/domain/task"
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
	till := time.Unix(0, b.StartedAt().Value()+int64(b.Duration()))
	n.Notify(
		"[lyco] ☕",
		fmt.Sprintf("Take a %v break till %v", b.Duration(), till.Format("15:04")),
	)
}

func NotifyForBreaksEnd(n Notifier, t task.Task) {
	till := time.Unix(0, t.StartedAt().Value()+t.Duration().Value())
	n.Notify(
		"[lyco] 🔨",
		fmt.Sprintf("Work for %v till %v on %s", t.Duration(), till.Format("15:04"), t.Name().Value()),
	)
}
