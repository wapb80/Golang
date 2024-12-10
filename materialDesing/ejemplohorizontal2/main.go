package main

import (
	"html/template"
	"net/http"
)

var tmpl = template.Must(template.ParseGlob("templates/*.html"))

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/menu", menuHandler)
	http.HandleFunc("/contenido", contenidoHandler)

	http.ListenAndServe(":8080", nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "layout", nil)
}

func menuHandler(w http.ResponseWriter, r *http.Request) {
	menu := r.URL.Query().Get("menu")
	tmpl.ExecuteTemplate(w, "menu", menu)
}

func contenidoHandler(w http.ResponseWriter, r *http.Request) {
	contenido := r.URL.Query().Get("contenido")
	tmpl.ExecuteTemplate(w, contenido, nil)
}
