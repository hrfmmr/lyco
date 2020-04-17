package main

import (
	"github.com/gcla/gowid"
	"github.com/gcla/gowid/examples"
	"github.com/gcla/gowid/widgets/dialog"
	"github.com/gcla/gowid/widgets/edit"
	"github.com/gcla/gowid/widgets/holder"
	"github.com/gcla/gowid/widgets/text"
)

//======================================================================

var (
	app   *gowid.App
	err   error
	input *dialog.Widget
)

func main() {
	f := examples.RedirectLogger("editor.log")
	defer f.Close()

	editWidget := edit.New()
	txt := text.New("hello, world")
	viewHolder := holder.New(txt)
	input = dialog.New(editWidget)
	input.Open(viewHolder, gowid.RenderWithRatio{R: 0.5}, app)
	app, err = gowid.NewApp(gowid.AppArgs{
		View: viewHolder,
	})
	examples.ExitOnErr(err)
	app.SimpleMainLoop()
}
