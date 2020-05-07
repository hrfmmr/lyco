package ui

import (
	"fmt"
	"time"

	"github.com/gcla/gowid"
	"github.com/gcla/gowid/widgets/pile"
	"github.com/gcla/gowid/widgets/text"
	"github.com/gcla/gowid/widgets/vpadding"
	"github.com/sirupsen/logrus"
)

var (
	taskText  *text.Widget
	timerText *text.Widget
)

func init() {
	taskText = text.New("", text.Options{
		Align: gowid.HAlignMiddle{},
	})
	timerText = text.New("", text.Options{
		Align: gowid.HAlignMiddle{},
	})
}

func NewCurrentTaskView() gowid.IWidget {
	return vpadding.New(
		pile.NewFlow(taskText, timerText),
		gowid.VAlignMiddle{},
		gowid.RenderFlow{},
	)
}

func updateTaskText(app gowid.IApp, s string) {
	logrus.Infof("🔃ui#updtaeTaskText s:%s", s)
	taskText.SetText(s, app)
}

func updateTimerText(app gowid.IApp, duration time.Duration, startedAt int64) {
	to := startedAt + int64(duration)
	now := time.Now().UnixNano()
	logrus.Infof("🔃ui#updtaeTimerText duration:%v startedAt:%v now:%v", duration, startedAt, now)
	r := to - now
	rsec := r / 1e9
	s := fmt.Sprintf("%02d:%02d", int(rsec/60)%60, rsec%60)
	timerText.SetText(s, app)
}
