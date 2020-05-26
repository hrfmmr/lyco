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

type keymap struct {
	key  string
	desc string
}

const (
	paletteCyan = "cyan"
	paletteInv  = "inv"
)

var (
	appContainer    *holder.Widget
	taskInputEditor *edit.Widget
	taskInputDialog *dialog.Widget
	palette         = gowid.Palette{
		paletteCyan: gowid.MakePaletteEntry(gowid.ColorCyan, gowid.ColorBlack),
		paletteInv:  gowid.MakePaletteEntry(gowid.ColorWhite, gowid.ColorBlack),
	}
	onStartTask  = make(chan string, 1)
	onPauseTask  = make(chan struct{}, 1)
	onResumeTask = make(chan struct{}, 1)
	onStopTask   = make(chan struct{}, 1)
)

func NewKeymap(key, desc string) *keymap {
	return &keymap{key, desc}
}

func Build() (*gowid.App, error) {
	appView := pile.New([]gowid.IContainerWidget{
		&gowid.ContainerWidget{
			IWidget: NewTaskView(),
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
		View:    appContainer,
		Palette: &palette,
		Log:     logrus.StandardLogger(),
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
		case tcell.KeyCtrlE:
			handled = true
			logrus.Info("⌨ui#UnhandledInput::case tcell.KeyCtrlE")
			var buf bytes.Buffer
			_, err := io.Copy(&buf, taskText)
			if err != nil {
				logrus.Fatal(err)
			}
			s := buf.String()
			onStartTask <- s
		case tcell.KeyCtrlP:
			handled = true
			logrus.Info("⌨ui#UnhandledInput::case tcell.KeyCtrlP")
			onPauseTask <- struct{}{}
		case tcell.KeyCtrlR:
			handled = true
			logrus.Info("⌨ui#UnhandledInput::case tcell.KeyCtrlR")
			onResumeTask <- struct{}{}
		case tcell.KeyCtrlQ:
			handled = true
			logrus.Info("⌨ui#UnhandledInput::case tcell.KeyCtrlQ")
			onStopTask <- struct{}{}
		}
	}
	return handled
}

func OnStartTask() <-chan string {
	return onStartTask
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
				onStartTask <- s
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
	keymaps := []*keymap{}
	for _, action := range task.AvailableActions() {
		keymaps = append(keymaps, convertTaskActionToKeymap(action))
	}
	app.Run(gowid.RunFunction(func(app gowid.IApp) {
		updateTaskText(app, task.Name())
		updateTimerText(app, task.RemainsTimerText())
		updateKeymaps(app, keymaps)
	}))
}

func convertTaskActionToKeymap(action dto.AvailableAction) *keymap {
	switch action {
	case dto.AvailableActionStart:
		return NewKeymap("C-e", "to start")
	case dto.AvailableActionPause:
		return NewKeymap("C-p", "to pause")
	case dto.AvailableActionResume:
		return NewKeymap("C-r", "to resume")
	case dto.AvailableActionAbort:
		return NewKeymap("C-q", "to abort")
	}
	return nil
}
