package main

import (
	"net/http"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func lineChartHandler(w http.ResponseWriter, _ *http.Request) {
	// Crear un gráfico de líneas
	line := charts.NewLine()

	// Datos para el eje X e Y
	xAxis := []string{"Lunes", "Martes", "Miércoles", "Jueves", "Viernes"}
	yAxis := []opts.LineData{
		{Value: 10}, {Value: 20}, {Value: 15}, {Value: 25}, {Value: 30},
	}

	// Configurar el gráfico
	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    "Gráfico de Líneas",
			Subtitle: "Ejemplo con go-echarts",
		}),
	)

	line.SetXAxis(xAxis).
		AddSeries("Ventas", yAxis)

	// Renderizar el gráfico en la respuesta HTTP
	line.Render(w)
}

func main() {
	http.HandleFunc("/", lineChartHandler)
	http.ListenAndServe(":8080", nil)
}
