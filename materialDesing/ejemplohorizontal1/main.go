package main

import (
	"html/template"
	"log"
	"net/http"
)

var tmpl = template.Must(template.ParseGlob("templates/*.html"))

func renderTemplate(w http.ResponseWriter, tmplName string, data interface{}) {
	err := tmpl.ExecuteTemplate(w, tmplName, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "base.html", map[string]interface{}{
			"Title": "Bienvenida",
		})
	})

	http.HandleFunc("/menu/reportes", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "menu_reportes.html", nil)
	})

	http.HandleFunc("/menu/comparativas", func(w http.ResponseWriter, r *http.Request) {
		// Similar para otros menús
		renderTemplate(w, "menu_comparativas.html", nil)
	})

	http.HandleFunc("/menu/georreferenciacion", func(w http.ResponseWriter, r *http.Request) {
		// Similar para otros menús
		renderTemplate(w, "menu_georreferenciacion.html", nil)
	})

	log.Println("Servidor iniciado en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
