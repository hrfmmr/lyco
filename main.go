package main

import (
	"bytes"
	"io"

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

	"github.com/sirupsen/logrus"
)

//======================================================================

var (
	app      *gowid.App
	err      error
	inpute   *edit.Widget
	inputd   *dialog.Widget
	tasktxt  *text.Widget
	timertxt *text.Widget
)

func init() {
	tasktxt = text.New("", text.Options{
		Align: gowid.HAlignMiddle{},
	})
	timertxt = text.New("timer", text.Options{
		Align: gowid.HAlignMiddle{},
	})
}

type (
	handler struct{}
)

func (h handler) UnhandledInput(app gowid.IApp, ev interface{}) bool {
	handled := false
	if evk, ok := ev.(*tcell.EventKey); ok {
		logrus.Infof("🐛 UnhandledInput evk:%v", evk.Key())
		switch evk.Key() {
		case tcell.KeyEsc, tcell.KeyCtrlC:
			logrus.Infof("🐛 case tcell.KeyEnter, tcell.KeyCtrlC ev:%v", ev)
			handled = true
			app.Quit()
		case tcell.KeyCtrlX:
			logrus.Infof("🐛 case tcell.KeyCtrlX:")
		case tcell.KeyEnter:
			//TODO: dialogに対するEnterキー押下で入力値をログ出力しようとしたが、ダメ。(dialog表示時にはこのUnhandledInput自体フックされず。)
			// dialog自体に対するイベントフックAPIがないか調査する🔍
			logrus.Infof("🐛 case tcell.KeyEnter:")
			handled = true
			var buf bytes.Buffer
			_, err := io.Copy(&buf, inpute)
			if err != nil {
				logrus.Fatal(err)
			}
			logrus.Infof("🐛 editor buf:%v", buf.String())
		}
	}
	return handled
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
			tasktxt.SetText(s, app)
			inputd.Close(app)
			handled = true
		}
		return handled
	}
}

func main() {
	f := examples.RedirectLogger("editor.log")
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
	app, err = gowid.NewApp(gowid.AppArgs{
		View: viewHolder,
		Log:  logrus.StandardLogger(),
	})
	examples.ExitOnErr(err)
	app.MainLoop(handler{})
}
