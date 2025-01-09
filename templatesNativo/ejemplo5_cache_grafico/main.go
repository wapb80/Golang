package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
)

// TemplateCache almacena las plantillas en memoria
type TemplateCache map[string]*template.Template

// TemplateData almacena datos dinámicos para las plantillas
type TemplateData struct {
	Title string
	Body  template.HTML
}

// newTemplateCache carga todas las plantillas del disco en el cache
func newTemplateCache(dir string) (TemplateCache, error) {
	cache := TemplateCache{}

	// Recorrer el directorio en busca de archivos .html
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".html" {
			// Obtener el nombre del archivo (por ejemplo, "home.html")
			name := filepath.Base(path)

			// Parsear el archivo y almacenarlo en el cache
			tmpl, err := template.ParseFiles(path)
			if err != nil {
				return err
			}
			cache[name] = tmpl
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return cache, nil
}

// addTemplateToCache agrega una nueva plantilla al cache en tiempo de ejecución
func addTemplateToCache(cache TemplateCache, name string, content string) error {
	// Parsear la nueva plantilla desde el contenido proporcionado
	tmpl, err := template.New(name).Parse(content)
	if err != nil {
		return err
	}

	// Agregarla al cache
	cache[name] = tmpl
	return nil
}

// renderTemplate renderiza una plantilla desde el cache
func renderTemplate(w http.ResponseWriter, cache TemplateCache, name string, data TemplateData) {
	tmpl, ok := cache[name]
	if !ok {
		http.Error(w, "Plantilla no encontrada", http.StatusInternalServerError)
		return
	}
	err := tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Error al renderizar la plantilla", http.StatusInternalServerError)
	}
}

func main() {
	// Inicializar el cache de plantillas
	cache, err := newTemplateCache("templates")
	if err != nil {
		panic(err)
	}

	// Handler para renderizar gráficos dinámicos
	http.HandleFunc("/chart", func(w http.ResponseWriter, r *http.Request) {
		// Crear un gráfico con go-echarts
		bar := charts.NewBar()
		bar.SetGlobalOptions(charts.WithTitleOpts(opts.Title{Title: "Gráfico Dinámico"}))
		bar.AddSeries("Ventas", []opts.BarData{
			{Value: 10}, {Value: 20}, {Value: 30},
		})

		// Renderizar el gráfico en un string HTML
		page := components.NewPage()
		page.AddCharts(bar)

		htmlContent := ""
		err := page.Render(w)
		if err != nil {
			fmt.Println("Error renderizando el gráfico:", err)
			return
		}

		// Crear una plantilla dinámica para el gráfico
		templateContent := `
		<!DOCTYPE html>
		<html lang="es">
		<head>
		    <title>{{.Title}}</title>
		</head>
		<body>
		    <h1>{{.Title}}</h1>
		    {{.Body}}
		</body>
		</html>
		`

		// Agregar la plantilla al cache
		addTemplateToCache(cache, "chart.html", templateContent)

		// Renderizar la plantilla con el gráfico
		renderTemplate(w, cache, "chart.html", TemplateData{
			Title: "Gráfico Dinámico",
			Body:  template.HTML(htmlContent),
		})
	})

	// Otros Handlers
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, cache, "home.html", TemplateData{Title: "Página de Inicio"})
	})

	http.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, cache, "about.html", TemplateData{Title: "Acerca de Nosotros"})
	})

	// Levantar el servidor
	fmt.Println("Servidor corriendo en http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
