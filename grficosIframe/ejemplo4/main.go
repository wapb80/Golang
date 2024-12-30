package main

import (
	"io"
	"math/rand"
	"net/http"
	"os"
	"sync"

	"html/template"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
)

var templates = template.Must(template.ParseFiles("templates/index.html", "templates/charts.html"))

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/generate", generateChartsHandler)
	http.ListenAndServe(":8080", nil)
}

var (
	itemCntLine = 6
	fruits      = []string{"Apple", "Banana", "Peach", "Lemon", "Pear", "Cherry"}
	tmpl        = template.Must(template.ParseGlob("templates/*.html"))
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "index.html", nil)
}

func renderTemplate(w http.ResponseWriter, tmplName string, data interface{}) {
	err := tmpl.ExecuteTemplate(w, tmplName, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func generateChartsHandler(w http.ResponseWriter, r *http.Request) {
	var wg sync.WaitGroup
	var chartsList []components.Charter
	var mu sync.Mutex // Mutex para proteger el acceso al slice

	// Iniciar las gorutinas para generar los gráficos
	wg.Add(2)

	// Gráfico de barras
	go func() {
		//  lineBase() // Enviar al canal en el orden deseado
		defer wg.Done()
		mu.Lock() // Bloqueamos el acceso al slice
		chartsList = append(chartsList, lineBase())
		mu.Unlock() // Liberamos el acceso al slice
	}()

	// Gráfico de líneas
	go func() {
		//  lineBase() // Enviar al canal en el orden deseado
		defer wg.Done()
		mu.Lock() // Bloqueamos el acceso al slice
		chartsList = append(chartsList, barBase())
		mu.Unlock() // Liberamos el acceso al slice
	}()

	// Esperar que todas las gorutinas terminen
	wg.Wait()

	// Crear una página y agregar los gráficos en el orden correcto
	page := components.NewPage()
	for _, chart := range chartsList {
		page.AddCharts(chart)
	}

	// Crea un archivo HTML y renderiza el contenido en él
	f, err := os.Create("templates/charts.html")
	if err != nil {
		panic(err)
	}
	page.Render(io.MultiWriter(f))
	tmpl.ExecuteTemplate(w, "charts.html", nil)
	//renderTemplate(w, "charts.html", nil)
}

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
