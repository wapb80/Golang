package main

import (
	"html/template"
	"log"
	"net/http"
)

// Compilamos todas las plantillas
var tmpl = template.Must(template.ParseGlob("templates/*.html"))

// Función para renderizar plantillas
func renderTemplate(w http.ResponseWriter, tmplName string, data interface{}) {
	err := tmpl.ExecuteTemplate(w, tmplName, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	// Ruta principal
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "base.html", map[string]interface{}{
			"Title": "Bienvenida",
		})
	})

	// Menús dinámicos
	http.HandleFunc("/menu/reportes", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "menu_reportes.html", nil)
	})
	http.HandleFunc("/menu/comparativas", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "menu_comparativas.html", nil)
	})
	http.HandleFunc("/menu/georreferenciacion", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "menu_georreferenciacion.html", nil)
	})

	// Servidor
	log.Println("Servidor iniciado en http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Error al iniciar el servidor: %s", err)
	}
}
