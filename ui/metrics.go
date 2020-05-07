package ui

import (
	"github.com/gcla/gowid"
	"github.com/gcla/gowid/widgets/divider"
	"github.com/gcla/gowid/widgets/fill"
	"github.com/gcla/gowid/widgets/pile"
	"github.com/gcla/gowid/widgets/table"
	"github.com/gcla/gowid/widgets/text"
	"github.com/gcla/gowid/widgets/vpadding"
)

func NewMetricsView() gowid.IWidget {
	// TODO:temporary metrics model
	metricsModel := table.NewSimpleModel(
		[]string{"title", "elapsed", "poms"},
		[][]string{
			{"task1", "1h25m19s", "5"},
			{"task2", "1h20m19s", "4"},
			{"task3", "1h2m19s", "3"},
			{"task4", "1h2m19s", "2"},
			{"task5", "1h2m19s", "1"},
		},
		table.SimpleOptions{
			Style: table.StyleOptions{
				HorizontalSeparator: divider.NewAscii(),
				TableSeparator:      divider.NewUnicode(),
				VerticalSeparator:   fill.New('|'),
			},
		},
	)
	metricsView := vpadding.New(
		pile.NewFlow(
			text.New("Total ⏰:1h24m59s 🍅:6poms"),
			table.New(metricsModel),
		),
		gowid.VAlignTop{},
		gowid.RenderFlow{},
	)
	return metricsView
}
