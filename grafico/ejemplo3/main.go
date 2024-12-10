package main

import (
	"io"
	"math/rand"
	"os"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
)

var (
	itemCnt = 7
	weeks   = []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}
)

func generateLineItems() []opts.LineData {
	items := make([]opts.LineData, 0)
	for i := 0; i < itemCnt; i++ {
		items = append(items, opts.LineData{Value: rand.Intn(300)})
	}
	return items
}

func lineSetToolbox() *charts.Line {
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "Toolbox for Line Chart"}),
		charts.WithToolboxOpts(opts.Toolbox{
			Right: "20%",
			Feature: &opts.ToolBoxFeature{
				SaveAsImage: &opts.ToolBoxFeatureSaveAsImage{
					Type:  "jpg",
					Title: "Anything you want",
				},
				DataView: &opts.ToolBoxFeatureDataView{
					Title: "DataView",
					// set the language
					Lang: []string{"data view", "turn off", "refresh"},
				},
			}}),
	)
	line.SetXAxis(weeks).
		AddSeries("Category A", generateLineItems()).
		AddSeries("Category B", generateLineItems())
	return line
}

type Exampler interface {
	Examples()
}

type LineExamples struct{}

func (LineExamples) Examples() {
	page := components.NewPage()
	page.AddCharts(
		lineSetToolbox(),
	)
	f, err := os.Create("examples/html/line.html")
	if err != nil {
		panic(err)
	}
	page.Render(io.MultiWriter(f))
}

func main() {
	examples := LineExamples{}
	examples.Examples()
}
