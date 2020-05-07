package ui

import (
	"github.com/gcla/gowid"
	"github.com/gcla/gowid/widgets/divider"
	"github.com/gcla/gowid/widgets/pile"
)

var (
	SwitchTaskChan = make(chan struct{}, 1)
)

func NewAppView() gowid.IWidget {
	view := pile.New([]gowid.IContainerWidget{
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
	return view
}
