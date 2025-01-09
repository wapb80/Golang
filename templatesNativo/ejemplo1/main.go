package main

import (
	"html/template"
	"log"
	"net/http"
)

// Función para renderizar plantillas
func renderTemplate(w http.ResponseWriter, name string, data interface{}) {
	tmpl, err := template.ParseFiles(
		"templates/base.html",
		"templates/"+name+".html",
	)
	if err != nil {
		log.Println("Error al cargar plantilla:", err)
		http.Error(w, "Error al cargar plantilla", http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "base.html", data)
	if err != nil {
		log.Println("Error al renderizar plantilla:", err)
		http.Error(w, "Error al renderizar plantilla", http.StatusInternalServerError)
	}
}

func main() {
	// Servir archivos estáticos
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Ruta para /
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "home", nil)
	})

	// Ruta para /home
	// http.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
	// 	renderTemplate(w, "home", nil)
	// })

	// Ruta para /users
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "users", nil)
	})

	log.Println("Servidor iniciado en http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
