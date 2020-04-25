package main

import (
	"bytes"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/gcla/gowid"
	"github.com/gcla/gowid/examples"
	"github.com/gcla/gowid/widgets/dialog"
	"github.com/gcla/gowid/widgets/edit"
	"github.com/gcla/gowid/widgets/holder"
	"github.com/gcla/gowid/widgets/pile"
	"github.com/gcla/gowid/widgets/text"
	"github.com/gcla/gowid/widgets/vpadding"
	"github.com/gdamore/tcell"
	"github.com/hrfmmr/lyco/appkeys"
	"github.com/hrfmmr/lyco/timer"

	"github.com/sirupsen/logrus"
)

//======================================================================

const (
	pomodoroDuration = 25 * time.Minute
)

var (
	wg            sync.WaitGroup
	app           *gowid.App
	err           error
	inpute        *edit.Widget
	inputd        *dialog.Widget
	tasktxt       *text.Widget
	timertxt      *text.Widget
	finChan       = make(chan struct{}, 1)
	startTaskChan = make(chan struct{}, 1)
)

func init() {
	tasktxt = text.New("", text.Options{
		Align: gowid.HAlignMiddle{},
	})
	timertxt = text.New("timer", text.Options{
		Align: gowid.HAlignMiddle{},
	})
}

func handleEnter() appkeys.KeyInputFn {
	return func(ev *tcell.EventKey, app gowid.IApp) bool {
		handled := false
		switch ev.Key() {
		case tcell.KeyEnter:
			logrus.Infof("🐛 main#handleEnter case tcell.KeyEnter")
			var buf bytes.Buffer
			_, err := io.Copy(&buf, inpute)
			if err != nil {
				logrus.Fatal(err)
			}
			s := buf.String()
			logrus.Infof("🐛 main#handleEnter case tcell.KeyEnter - 📝editor buf:%v", s)
			inputd.Close(app)
			tasktxt.SetText(s, app)
			startTaskChan <- struct{}{}
			handled = true
		}
		return handled
	}
}

func startTask() {
	logrus.Infof("🐛 main#startTask")
	t := timer.New()
	for d := range t.Start(time.Second, pomodoroDuration) {
		seconds := d / 1e9
		remaining := fmt.Sprintf("%02d:%02d", int(seconds/60)%60, seconds%60)
		logrus.Infof("🐛 main#startTask - remaining:%v", remaining)
		app.Run(gowid.RunFunction(func(app gowid.IApp) {
			timertxt.SetText(remaining, app)
		}))
	}
}

func main() {
	f := examples.RedirectLogger("lyco.log")
	defer f.Close()

	flow := gowid.RenderFlow{}
	txts := pile.New([]gowid.IContainerWidget{
		&gowid.ContainerWidget{
			IWidget: tasktxt,
			D:       flow,
		},
		&gowid.ContainerWidget{
			IWidget: timertxt,
			D:       flow,
		},
	})
	status := vpadding.New(txts, gowid.VAlignMiddle{}, flow)
	viewHolder := holder.New(status)
	inpute = edit.New()
	onelineEd := appkeys.New(
		inpute,
		handleEnter(),
		appkeys.Options{
			ApplyBefore: true,
		})
	inputd = dialog.New(onelineEd)
	inputd.Open(viewHolder, gowid.RenderWithRatio{R: 0.5}, app)

	wg.Add(1)
	go func() {
		defer wg.Done()
	Loop:
		for {
			select {
			case <-finChan:
				logrus.Infof("🔴 #main case <-finChan")
				break Loop
			case <-startTaskChan:
				logrus.Infof("🐛 #main case <-startTaskChan")
				go startTask()
			}
		}
	}()

	app, err = gowid.NewApp(gowid.AppArgs{
		View: viewHolder,
		Log:  logrus.StandardLogger(),
	})
	examples.ExitOnErr(err)
	app.SimpleMainLoop()
	finChan <- struct{}{}
	wg.Wait()
}
