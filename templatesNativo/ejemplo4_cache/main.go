package main

import (
	"html/template"
	"net/http"
	"os"
	"path/filepath"
)

// TemplateCache es un mapa para almacenar las plantillas compiladas
type TemplateCache map[string]*template.Template

// Función para inicializar el cache de plantillas
func newTemplateCache(dir string) (TemplateCache, error) {
	cache := TemplateCache{}

	// Recorrer el directorio para encontrar todos los archivos .html
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Solo procesar archivos con extensión .html
		if !info.IsDir() && filepath.Ext(path) == ".html" {
			// Obtener el nombre base del archivo (por ejemplo, "home.html")
			name := filepath.Base(path)

			// Parsear el archivo como plantilla
			tmpl, err := template.ParseFiles(path)
			if err != nil {
				return err
			}

			// Almacenar en el cache
			cache[name] = tmpl
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return cache, nil
}

// Estructura para almacenar datos dinámicos en las plantillas
type TemplateData struct {
	Title string
}

func main() {
	// Crear el cache de plantillas al iniciar el servidor
	cache, err := newTemplateCache("templates")
	if err != nil {
		panic(err)
	}

	// Handlers
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, cache, "home.html", TemplateData{Title: "Página de Inicio"})
	})

	http.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, cache, "about.html", TemplateData{Title: "Acerca de Nosotros"})
	})

	// Levantar el servidor
	http.ListenAndServe(":8080", nil)
}

// renderTemplate utiliza el cache para renderizar la plantilla
func renderTemplate(w http.ResponseWriter, cache TemplateCache, name string, data TemplateData) {
	// Buscar la plantilla en el cache
	tmpl, ok := cache[name]
	if !ok {
		http.Error(w, "Plantilla no encontrada", http.StatusInternalServerError)
		return
	}

	// Ejecutar la plantilla con los datos
	err := tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Error al renderizar la plantilla", http.StatusInternalServerError)
	}
}
