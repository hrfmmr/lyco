package ui

import (
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

func updateTimerText(app gowid.IApp, remains string) {
	timerText.SetText(remains, app)
}
