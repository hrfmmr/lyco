package ui

import (
	"fmt"
	"strings"

	"github.com/gcla/gowid"
	"github.com/gcla/gowid/widgets/pile"
	"github.com/gcla/gowid/widgets/styled"
	"github.com/gcla/gowid/widgets/text"
	"github.com/gcla/gowid/widgets/vpadding"
)

const (
	taskFooterStatus = iota
	taskFooterKeymaps
)

var (
	taskText          *text.Widget
	timerText         *text.Widget
	taskFooterContent []text.ContentSegment
	taskFooterText    *styled.Widget
)

func init() {
	taskText = text.New("", text.Options{
		Align: gowid.HAlignMiddle{},
	})
	timerText = text.New("", text.Options{
		Align: gowid.HAlignMiddle{},
	})
	taskFooterContent = []text.ContentSegment{
		text.StyledContent("", gowid.MakePaletteRef(paletteCyan)),
		text.StyledContent("", gowid.MakePaletteRef(paletteInv)),
	}
	taskFooterText = styled.New(
		text.NewFromContent(text.NewContent(taskFooterContent)),
		gowid.MakePaletteRef(paletteInv),
	)
}

func NewCurrentTaskView() gowid.IWidget {
	return vpadding.New(
		pile.NewFlow(taskText, timerText),
		gowid.VAlignMiddle{},
		gowid.RenderFlow{},
	)
}

func NewTaskView() gowid.IWidget {
	currentTaskView := vpadding.New(
		pile.NewFlow(taskText, timerText),
		gowid.VAlignMiddle{},
		gowid.RenderFlow{},
	)
	return pile.New([]gowid.IContainerWidget{
		&gowid.ContainerWidget{
			IWidget: currentTaskView,
			D:       gowid.RenderWithWeight{W: 3},
		},
		&gowid.ContainerWidget{
			IWidget: taskFooterText,
			D:       gowid.RenderFlow{},
		},
	})
}

func updateTaskText(app gowid.IApp, s string) {
	taskText.SetText(s, app)
}

func updateTimerText(app gowid.IApp, remains string) {
	timerText.SetText(remains, app)
}

func updateKeymaps(app gowid.IApp, keymaps []*keymap) {
	var sb strings.Builder
	for i, k := range keymaps {
		sb.WriteString(fmt.Sprintf("%s %s", k.key, k.desc))
		if i < len(keymaps)-1 {
			sb.WriteString(", ")
		}
	}
	taskFooterContent[taskFooterKeymaps].Text = sb.String()
	taskFooterText.SetSubWidget(
		text.NewFromContent(text.NewContent(taskFooterContent)),
		app,
	)
}
