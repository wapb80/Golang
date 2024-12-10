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
	itemCntLine = 6
	fruits      = []string{"Apple", "Banana", "Peach", "Lemon", "Pear", "Cherry"}
)

func generateLineItems() []opts.LineData {
	items := make([]opts.LineData, 0)
	for i := 0; i < itemCntLine; i++ {
		items = append(items, opts.LineData{Value: rand.Intn(300)})
	}
	return items
}

func generateBarItems() []opts.BarData {
	items := make([]opts.BarData, 0)
	for i := 0; i < itemCntLine; i++ {
		items = append(items, opts.BarData{Value: rand.Intn(300)})
	}
	return items
}

func lineBase() *charts.Line {
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "Basic Line Example", Subtitle: "This is the subtitle."}),
	)

	line.SetXAxis(fruits).
		AddSeries("Category A", generateLineItems())
	return line
}

func barBase() *charts.Bar {
	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "Basic Bar Chart", Subtitle: "A simple bar chart."}),
	)

	bar.SetXAxis(fruits).
		AddSeries("Category B", generateBarItems())
	return bar
}

func main() {
	// Crea una página con varios gráficos
	page := components.NewPage()
	page.AddCharts(
		lineBase(), // Agrega gráfico de líneas
		barBase(),  // Agrega gráfico de barras
	)

	// Crea un archivo HTML y renderiza el contenido en él
	f, err := os.Create("examples/html/line_and_bar.html")
	if err != nil {
		panic(err)
	}
	page.Render(io.MultiWriter(f))
}
