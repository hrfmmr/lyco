package ui

import (
	"bytes"
	"fmt"
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

type mode int

const (
	modeTask mode = iota
	modeBreaks
)

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
	onStartTask   = make(chan string, 1)
	onPauseTask   = make(chan struct{}, 1)
	onResumeTask  = make(chan struct{}, 1)
	onStopTask    = make(chan struct{}, 1)
	onSwitchTask  = make(chan string, 1)
	onAbortBreaks = make(chan struct{}, 1)
	currentMode   = modeTask
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
		if evk.Key() == tcell.KeyCtrlC {
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
		}
		switch currentMode {
		case modeTask:
			return handleTaskKeyInput(app, evk)
		case modeBreaks:
			return handleBreaksKeyInput(app, evk)
		}
	}
	return handled
}

func handleTaskKeyInput(app gowid.IApp, k *tcell.EventKey) (handled bool) {
	handled = false
	switch k.Key() {
	case tcell.KeyCtrlE:
		handled = true
		logrus.Info("⌨ui#handleTaskKeyInput::case tcell.KeyCtrlE")
		var buf bytes.Buffer
		_, err := io.Copy(&buf, taskText)
		if err != nil {
			logrus.Fatal(err)
		}
		s := buf.String()
		onStartTask <- s
	case tcell.KeyCtrlP:
		handled = true
		logrus.Info("⌨ui#handleTaskKeyInput::case tcell.KeyCtrlP")
		onPauseTask <- struct{}{}
	case tcell.KeyCtrlR:
		handled = true
		logrus.Info("⌨ui#handleTaskKeyInput::case tcell.KeyCtrlR")
		onResumeTask <- struct{}{}
	case tcell.KeyCtrlQ:
		handled = true
		logrus.Info("⌨ui#handleTaskKeyInput::case tcell.KeyCtrlQ")
		onStopTask <- struct{}{}
	case tcell.KeyCtrlS:
		handled = true
		logrus.Info("⌨ui#handleTaskKeyInput::case tcell.KeyCtrlS")
		switchTask(app)
	}
	return
}

func handleBreaksKeyInput(app gowid.IApp, k *tcell.EventKey) (handled bool) {
	handled = false
	switch k.Key() {
	case tcell.KeyCtrlQ:
		handled = true
		logrus.Info("⌨ui#handleBreaksKeyInput::case tcell.KeyCtrlQ")
		onAbortBreaks <- struct{}{}
	}
	return
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

func OnSwitchTask() <-chan string {
	return onSwitchTask
}

func OnAbortBreaks() <-chan struct{} {
	return onAbortBreaks
}

func StartTask(app gowid.IApp) {
	showTaskInputDialog(app, true)
}

func switchTask(app gowid.IApp) {
	showTaskInputDialog(app, false)
}

func showTaskInputDialog(app gowid.IApp, bootstrap bool) {
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
				switch bootstrap {
				case true:
					onStartTask <- s
				case false:
					onSwitchTask <- s
				}
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

func UpdateTask2(app gowid.IApp, task dto.TaskState) {
	currentMode = modeTask
	logrus.Infof("🔃ui#UpdateTask2 task:%v", task)
	keymaps := []*keymap{}
	for _, action := range task.AvailableActions() {
		keymaps = append(keymaps, convertTaskActionToKeymap(action))
	}
	app.Run(gowid.RunFunction(func(app gowid.IApp) {
		updateTaskText(app, task.TaskName())
		updateTimerText(app, task.RemainsTimerText())
		updateKeymaps(app, keymaps)
	}))
}

//TODO: remove
func UpdateTask(app gowid.IApp, task dto.TaskDTO) {
	currentMode = modeTask
	logrus.Infof("🔃ui#UpdateTask task:%v", task)
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

//TODO: remove
func UpdateBreaks(app gowid.IApp, b dto.BreaksDTO) {
	currentMode = modeBreaks
	logrus.Infof("🔃ui#UpdateBreaks breaks:%v", b)
	keymaps := []*keymap{
		NewKeymap("C-q", "to abort breaks"),
	}
	app.Run(gowid.RunFunction(func(app gowid.IApp) {
		updateTimerText(app, fmt.Sprintf("☕ %s", b.RemainsTimerText()))
		updateKeymaps(app, keymaps)
	}))
}

func convertTaskActionToKeymap(action dto.AvailableTaskAction) *keymap {
	switch action {
	case dto.AvailableTaskActionStart:
		return NewKeymap("C-e", "to start")
	case dto.AvailableTaskActionPause:
		return NewKeymap("C-p", "to pause")
	case dto.AvailableTaskActionResume:
		return NewKeymap("C-r", "to resume")
	case dto.AvailableTaskActionStop:
		return NewKeymap("C-q", "to stop")
	case dto.AvailableTaskActionSwitch:
		return NewKeymap("C-s", "to switch")
	}
	return nil
}
