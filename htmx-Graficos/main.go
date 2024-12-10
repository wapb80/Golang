package main

import (
	"bytes"
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"golang.org/x/net/html"
)

// Estructura para pasar datos a la plantilla
type PageData struct {
	LineChartHTML template.HTML
	BarChartHTML  template.HTML
}

// Función para generar el gráfico de línea
func generateLineChart() (template.HTML, error) {
	// Crear un gráfico de línea
	line := charts.NewLine()

	// Configurar las opciones del gráfico
	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "Gráfico de Línea con go-echarts"}),
	)

	// Crear una nueva fuente de aleatoriedad utilizando el tiempo actual como semilla
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)

	// Datos del eje X (días de la semana)
	xAxisData := []string{"Lunes", "Martes", "Miércoles", "Jueves", "Viernes"}

	// Datos aleatorios para el eje Y
	var yAxisData []opts.LineData
	for i := 0; i < len(xAxisData); i++ {
		yAxisData = append(yAxisData, opts.LineData{Value: rng.Intn(50) + 10}) // Valores entre 10 y 59
	}

	// Añadir los datos al gráfico
	line.SetXAxis(xAxisData).AddSeries("Serie 1", yAxisData)

	// Renderizar el gráfico en un buffer
	var buffer bytes.Buffer
	if err := line.Render(&buffer); err != nil {
		return "", fmt.Errorf("error al renderizar el gráfico: %v", err)
	}

	return template.HTML(buffer.String()), nil
}

// Función para generar el gráfico de barras
func generateBarChart() (template.HTML, error) {
	// Crear un gráfico de barras
	bar := charts.NewBar()

	// Configurar las opciones del gráfico
	bar.SetGlobalOptions(
	// charts.WithTitleOpts(opts.Title{Title: "Gráfico de Barras con go-echarts"}),

	)

	// Datos del eje X (días de la semana)
	xAxisData := []string{"Lunes", "Martes", "Miércoles", "Jueves", "Viernes"}

	// Datos aleatorios para el eje Y
	var yAxisData []opts.BarData
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)
	for i := 0; i < len(xAxisData); i++ {
		yAxisData = append(yAxisData, opts.BarData{Value: rng.Intn(50) + 10}) // Valores entre 10 y 59
	}

	// Añadir los datos al gráfico
	bar.SetXAxis(xAxisData).AddSeries("Serie 1", yAxisData)

	// Renderizar el gráfico en un buffer
	var buffer bytes.Buffer
	if err := bar.Render(&buffer); err != nil {
		return "", fmt.Errorf("error al renderizar el gráfico: %v", err)
	}

	return template.HTML(buffer.String()), nil
}

// Función para renderizar la página principal
func mainPage(w http.ResponseWriter, r *http.Request) {
	// Generar los gráficos
	lineChartHTML, err := generateLineChart()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al generar el gráfico de líneas: %v", err), http.StatusInternalServerError)
		return
	}

	barChartHTML, err := generateBarChart()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al generar el gráfico de barras: %v", err), http.StatusInternalServerError)
		return
	}

	// Cargar la plantilla
	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al cargar la plantilla: %v", err), http.StatusInternalServerError)
		return
	}

	// Renderizar la plantilla
	pageData := PageData{
		LineChartHTML: lineChartHTML,
		BarChartHTML:  barChartHTML,
	}
	if err := tmpl.Execute(w, pageData); err != nil {
		// Evitar escribir después de un fallo
		fmt.Printf("Error al renderizar la plantilla: %v\n", err)
		return
	}
}

// Función para actualizar el gráfico de líneas
func updateLineChart(w http.ResponseWriter, r *http.Request) {
	// Devolver solo el gráfico de líneas actualizado con datos aleatorios
	chartHTML, err := generateLineChart()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al generar el gráfico de líneas: %v", err), http.StatusInternalServerError)
		return
	}

	bodyContent, err := extractBodyContentFromTemplateHTML(chartHTML)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Fprint(w, bodyContent)
	println(bodyContent)
}

/***********************************/
// Función para extraer solo el contenido dentro de <body> de una cadena de tipo template.HTML
// Función para extraer el contenido dentro de <body> incluyendo las etiquetas <body>
func extractBodyContentFromTemplateHTML(htmlContent template.HTML) (template.HTML, error) {
	// Convertir template.HTML a string
	htmlString := string(htmlContent)

	// Parsear el HTML y buscar el nodo <body>
	doc, err := html.Parse(strings.NewReader(htmlString))
	if err != nil {
		return "", fmt.Errorf("error al parsear el HTML: %v", err)
	}

	// Buscar el nodo <body> y devolver su HTML completo
	var bodyContent string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "body" {
			// Renderizar el nodo <body> y almacenarlo
			var sb strings.Builder
			html.Render(&sb, n)
			bodyContent = sb.String()
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	// Iniciar la búsqueda
	f(doc)

	// Verificar si encontramos el contenido dentro de <body>
	if bodyContent == "" {
		return "", fmt.Errorf("no se encontró contenido dentro de las etiquetas <body>")
	}

	// Devolver el contenido de <body> como template.HTML
	return template.HTML(bodyContent), nil
}

/***********************************/
func main() {
	// Ruta para la página principal
	http.HandleFunc("/", mainPage)

	// Ruta para actualizar solo el gráfico de líneas
	http.HandleFunc("/update-line-chart", updateLineChart)

	// Iniciar el servidor
	fmt.Println("Servidor ejecutándose en http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error al iniciar el servidor: %v\n", err)
	}
}
