package ui

import (
	"bytes"
	"io"

	"github.com/gcla/gowid"
	"github.com/gcla/gowid/widgets/dialog"
	"github.com/gcla/gowid/widgets/divider"
	"github.com/gcla/gowid/widgets/edit"
	"github.com/gcla/gowid/widgets/framed"
	"github.com/gcla/gowid/widgets/holder"
	"github.com/gcla/gowid/widgets/hpadding"
	"github.com/gcla/gowid/widgets/pile"
	"github.com/gcla/gowid/widgets/text"
	"github.com/gdamore/tcell"
	"github.com/hrfmmr/lyco/appkeys"
	"github.com/hrfmmr/lyco/application/dto"
	"github.com/sirupsen/logrus"
)

var (
	appContainer    *holder.Widget
	taskInputEditor *edit.Widget
	taskInputDialog *dialog.Widget
	onSubmitTask    = make(chan string, 1)
	onPauseTask     = make(chan struct{}, 1)
	onResumeTask    = make(chan struct{}, 1)
	onStopTask      = make(chan struct{}, 1)
)

func Build() (*gowid.App, error) {
	appView := pile.New([]gowid.IContainerWidget{
		&gowid.ContainerWidget{
			IWidget: NewCurrentTaskView(),
			D:       gowid.RenderWithRatio{R: 0.2},
		},
		&gowid.ContainerWidget{
			IWidget: divider.NewAscii(),
			D:       gowid.RenderFlow{},
		},
		&gowid.ContainerWidget{
			IWidget: NewMetricsView(),
			D:       gowid.RenderFlow{},
		},
	})
	appContainer = holder.New(appView)
	return gowid.NewApp(gowid.AppArgs{
		View: appContainer,
		Log:  logrus.StandardLogger(),
	})
}

func UnhandledInput(app gowid.IApp, event interface{}) bool {
	handled := false
	if evk, ok := event.(*tcell.EventKey); ok {
		switch evk.Key() {
		case tcell.KeyCtrlC:
			logrus.Info("⌨ui#UnhandledInput::case tcell.KeyCtrlC")
			handled = true
			msg := text.New("Do you want to quit?")
			yesno := dialog.New(
				framed.NewSpace(hpadding.New(msg, gowid.HAlignMiddle{}, gowid.RenderFixed{})),
				dialog.Options{
					Buttons: dialog.OkCancel,
				},
			)
			yesno.Open(appContainer, gowid.RenderWithRatio{R: 0.5}, app)
		case tcell.KeyCtrlP:
			logrus.Info("⌨ui#UnhandledInput::case tcell.KeyCtrlP")
			onPauseTask <- struct{}{}
		case tcell.KeyCtrlR:
			logrus.Info("⌨ui#UnhandledInput::case tcell.KeyCtrlR")
			onResumeTask <- struct{}{}
		case tcell.KeyCtrlQ:
			logrus.Info("⌨ui#UnhandledInput::case tcell.KeyCtrlQ")
			onStopTask <- struct{}{}
		}
	}
	return handled
}

func OnStartTask() <-chan string {
	return onSubmitTask
}

func OnPauseTask() <-chan struct{} {
	return onPauseTask
}

func OnResumeTask() <-chan struct{} {
	return onResumeTask
}

func OnStopTask() <-chan struct{} {
	return onStopTask
}

func SwitchTask(app gowid.IApp) {
	taskInputEditor = edit.New()
	onelineEditor := appkeys.New(
		taskInputEditor,
		func(ev *tcell.EventKey, app gowid.IApp) bool {
			handled := false
			switch ev.Key() {
			case tcell.KeyEnter:
				handled = true
				logrus.Infof("🐛 SwitchTask::case tcell.KeyEnter")
				var buf bytes.Buffer
				_, err := io.Copy(&buf, taskInputEditor)
				if err != nil {
					logrus.Fatal(err)
				}
				s := buf.String()
				logrus.Infof("🐛 SwitchTask::case tcell.KeyEnter - 📝editor buf:%v", s)
				taskInputDialog.Close(app)
				onSubmitTask <- s
			}
			return handled
		},
		appkeys.Options{
			ApplyBefore: true,
		},
	)
	taskInputDialog = dialog.New(onelineEditor)
	taskInputDialog.Open(appContainer, gowid.RenderWithRatio{R: 0.5}, app)
}

func Update(app gowid.IApp, task dto.TaskDTO) {
	logrus.Infof("🔃ui#Update task:%v", task)
	app.Run(gowid.RunFunction(func(app gowid.IApp) {
		updateTaskText(app, task.Name())
		updateTimerText(app, task.RemainsTimerText())
	}))
}
