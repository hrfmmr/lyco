package ui

import (
	"bytes"
	"io"

	"github.com/gcla/gowid"
	"github.com/gcla/gowid/widgets/dialog"
	"github.com/gcla/gowid/widgets/divider"
	"github.com/gcla/gowid/widgets/edit"
	"github.com/gcla/gowid/widgets/holder"
	"github.com/gcla/gowid/widgets/pile"
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

func OnStartTask() <-chan string {
	return onSubmitTask
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
		updateTimerText(app, task.Duration(), task.StartedAt())
	}))
}
