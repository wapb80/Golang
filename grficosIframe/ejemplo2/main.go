package main

import (
	"bytes"
	"math/rand"
	"net/http"
	"sync"
	"text/template"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

var templates = template.Must(template.ParseFiles("templates/index.html", "templates/charts.html"))

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/generate", generateChartsHandler)
	http.ListenAndServe(":8080", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func generateChartsHandler(w http.ResponseWriter, r *http.Request) {
	var wg sync.WaitGroup
	var barHTML, lineHTML, pieHTML string

	// Generar gráficos en goroutines
	wg.Add(3)

	go func() {
		defer wg.Done()
		barHTML = renderBarChart()
	}()

	go func() {
		defer wg.Done()
		lineHTML = renderLineChart()
	}()

	go func() {
		defer wg.Done()
		pieHTML = renderPieChart()
	}()

	// Esperar a que todas las goroutines terminen
	wg.Wait()

	// Renderizar el template con los gráficos generados
	data := map[string]string{
		"BarChart":  barHTML,
		"LineChart": lineHTML,
		"PieChart":  pieHTML,
	}

	err := templates.ExecuteTemplate(w, "charts.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func renderBarChart() string {
	bar := charts.NewBar()
	bar.SetGlobalOptions(charts.WithTitleOpts(opts.Title{Title: "Gráfico de Barras"}))

	xAxis := []string{"Grupo 1", "Grupo 2", "Grupo 3"}
	bar.SetXAxis(xAxis).
		AddSeries("Datos", generateRandomData())

	var buf bytes.Buffer
	err := bar.Render(&buf)
	if err != nil {
		return ""
	}
	return buf.String()
}

func renderLineChart() string {
	line := charts.NewLine()
	line.SetGlobalOptions(charts.WithTitleOpts(opts.Title{Title: "Gráfico de Líneas"}))

	xAxis := []string{"Enero", "Febrero", "Marzo"}
	line.SetXAxis(xAxis).
		AddSeries("Datos", generateRandomLineData())

	var buf bytes.Buffer
	err := line.Render(&buf)
	if err != nil {
		return ""
	}
	return buf.String()
}

func renderPieChart() string {
	pie := charts.NewPie()
	pie.SetGlobalOptions(charts.WithTitleOpts(opts.Title{Title: "Gráfico Circular"}))

	pie.AddSeries("Datos", generatePieData())

	var buf bytes.Buffer
	err := pie.Render(&buf)
	if err != nil {
		return ""
	}
	return buf.String()
}

func generateRandomData() []opts.BarData {
	items := make([]opts.BarData, 0)
	for i := 0; i < 3; i++ {
		items = append(items, opts.BarData{Value: rand.Intn(100)})
	}
	return items
}

func generatePieData() []opts.PieData {
	items := make([]opts.PieData, 0)
	for i := 0; i < 3; i++ {
		items = append(items, opts.PieData{Name: "Categoría " + string(i+1), Value: rand.Intn(100)})
	}
	return items
}

func generateRandomLineData() []opts.LineData {
	items := make([]opts.LineData, 0)
	for i := 0; i < 3; i++ {
		items = append(items, opts.LineData{Value: rand.Intn(100)})
	}
	return items
}
