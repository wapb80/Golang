package main

import (
	"bytes"
	"math/rand"
	"net/http"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"

	"html/template"
	"sync"
)

var templates = template.Must(template.ParseFiles("templates/index.html", "templates/charts.html"))

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/generate", generateChartsHandler)
	http.HandleFunc("/chart/bar", barChartHandler)
	http.HandleFunc("/chart/line", lineChartHandler)
	http.HandleFunc("/chart/pie", pieChartHandler)
	http.ListenAndServe(":8080", nil)
}

var (
	itemCntLine = 6
	fruits      = []string{"Apple", "Banana", "Peach", "Lemon", "Pear", "Cherry"}
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "index.html", nil)
}

func generateChartsHandler(w http.ResponseWriter, r *http.Request) {
	// Obtener valores seleccionados
	age := r.FormValue("age")
	gender := r.FormValue("gender")
	nationality := r.FormValue("nationality")

	var wg sync.WaitGroup
	chartHTMLs := make(map[string]string)

	// Generar gráficos en goroutines
	wg.Add(3)

	go func() {
		defer wg.Done()
		chartHTMLs["bar"] = renderBarChart(age, gender, nationality)
	}()

	go func() {
		defer wg.Done()
		chartHTMLs["line"] = renderLineChart(age, gender, nationality)
	}()

	go func() {
		defer wg.Done()
		chartHTMLs["pie"] = renderPieChart(age, gender, nationality)
	}()

	// Esperar a que todas las goroutines terminen
	wg.Wait()

	// Almacenar en Memcached (opcional)

	// Renderizar todos los gráficos en un solo iframe usando un template
	err := templates.ExecuteTemplate(w, "charts.html", chartHTMLs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func renderBarChart(age, gender, nationality string) string {
	bar := charts.NewBar()
	bar.SetGlobalOptions(charts.WithTitleOpts(opts.Title{Title: "Gráfico de Barras"}))

	xAxis := []string{"Grupo 1", "Grupo 2", "Grupo 3"}
	bar.SetXAxis(xAxis).
		AddSeries("Datos", generateRandomData())

	// Crear un buffer para almacenar el HTML generado
	var buf bytes.Buffer
	err := bar.Render(&buf)
	if err != nil {
		return "" // Maneja el error adecuadamente
	}
	return buf.String()

}

func renderLineChart(age, gender, nationality string) string {
	line := charts.NewLine()
	line.SetGlobalOptions(charts.WithTitleOpts(opts.Title{Title: "Gráfico de Líneas"}))

	xAxis := []string{"Enero", "Febrero", "Marzo"}
	line.SetXAxis(xAxis).
		AddSeries("Datos", generateRandomLineData())

	var buf bytes.Buffer
	err := line.Render(&buf)
	if err != nil {
		return "" // Maneja el error adecuadamente
	}
	return buf.String()
}

func renderPieChart(age, gender, nationality string) string {
	pie := charts.NewPie()
	pie.SetGlobalOptions(charts.WithTitleOpts(opts.Title{Title: "Gráfico Circular"}))

	pie.AddSeries("Datos", generatePieData())

	var buf bytes.Buffer
	err := pie.Render(&buf)
	if err != nil {
		return "" // Maneja el error adecuadamente
	}
	return buf.String()
}

// Función para generar datos aleatorios para gráficos
func generateRandomData() []opts.BarData {
	items := make([]opts.BarData, 0)
	for i := 0; i < 3; i++ {
		items = append(items, opts.BarData{Value: rand.Intn(100)})
	}
	return items
}

// Función para generar datos para el gráfico circular
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

func barChartHandler(w http.ResponseWriter, r *http.Request) {
	chartHTML := renderBarChart("", "", "")
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(chartHTML))
}

func lineChartHandler(w http.ResponseWriter, r *http.Request) {
	chartHTML := renderLineChart("", "", "")
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(chartHTML))
}

func pieChartHandler(w http.ResponseWriter, r *http.Request) {
	chartHTML := renderPieChart("", "", "")
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(chartHTML))
}
