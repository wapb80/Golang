// aca tengo una funcion para traer leer todas las funciones del directorio y subdirectorios
package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func loadTemplates(dir string) *template.Template {
	var files []string

	// Buscar archivos .html en todos los subdirectorios
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Filtrar solo archivos con extensión .html
		if !info.IsDir() && filepath.Ext(path) == ".html" {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	// Parsear las plantillas encontradas
	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		panic(err)
	}
	return tmpl
}

var (
	tmpl *template.Template
)

// renderTemplate: Función para renderizar plantillas con encabezado, menú y contenido
func renderTemplate(w http.ResponseWriter, name string, data interface{}) {
	tmpl = loadTemplates("templates")
	println(name)
	err := tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Println("Error al renderizar plantilla:", err)
		http.Error(w, "Error al renderizar plantilla", http.StatusInternalServerError)
	}
}

func main() {
	// Servir archivos estáticos (CSS, JS, imágenes)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Rutas
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "home", nil)
	})

	http.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "about", nil)
	})

	http.HandleFunc("/contact", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "contact", nil)
	})

	http.HandleFunc("/services", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "services", nil)
	})

	http.HandleFunc("/portfolio", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "portafolio", nil)
	})

	// Iniciar servidor
	log.Println("Servidor iniciado en http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
