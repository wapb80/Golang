// se cargan todas las paginas , se separo el contenido en una template diferente
package main

import (
	"html/template"
	"log"
	"net/http"
)

// renderTemplate: Función para renderizar plantillas con encabezado, menú y contenido
func renderTemplate(w http.ResponseWriter, name string, data interface{}) {
	tmpl, err := template.ParseFiles(
		"templates/base.html",             // Plantilla base
		"templates/header.html",           // Encabezado
		"templates/sidebar.html",          // Menú lateral
		"templates/content/"+name+".html", // Contenido dinámico
	)
	if err != nil {
		log.Println("Error al cargar plantilla:", err)
		http.Error(w, "Error al cargar plantilla", http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "base", data)
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
