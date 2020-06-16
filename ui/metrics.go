package ui

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gcla/gowid"
	"github.com/gcla/gowid/widgets/pile"
	"github.com/gcla/gowid/widgets/table"
	"github.com/gcla/gowid/widgets/text"
	"github.com/gcla/gowid/widgets/vpadding"
	"github.com/hrfmmr/lyco/application/dto"
)

var (
	metricsHeaderText  *text.Widget
	metricsTable       *table.Widget
	metricsTableHeader = []string{"title", "elapsed", "poms"}
)

func NewMetricsView() gowid.IWidget {
	metricsModel := table.NewSimpleModel(
		metricsTableHeader,
		[][]string{},
	)
	metricsHeaderText = text.New("")
	metricsTable = table.New(metricsModel)
	metricsView := vpadding.New(
		pile.NewFlow(
			metricsHeaderText,
			metricsTable,
		),
		gowid.VAlignTop{},
		gowid.RenderFlow{},
	)
	return metricsView
}

func updateMetricsHeaderText(app gowid.IApp, totalElapsed time.Duration, totalPomsCount uint64) {
	metricsHeaderText.SetText(
		fmt.Sprintf("Total ⏰%v 🍅%dpoms", totalElapsed, totalPomsCount),
		app,
	)
}

func updateMetricsModel(app gowid.IApp, entries []dto.MetricsEntry) {
	data := make([][]string, len(entries))
	for i, e := range entries {
		data[i] = []string{
			e.Name(),
			e.Elapsed().String(),
			strconv.FormatUint(e.PomsCount(), 10),
		}
	}
	model := table.NewSimpleModel(
		metricsTableHeader,
		data,
	)
	metricsTable.SetModel(model, app)
}
