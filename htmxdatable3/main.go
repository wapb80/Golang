package main

import (
	"net/http"
	"text/template"
)

var templates = template.Must(template.ParseGlob("templates/*.html"))

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/listUsers", listUsersHandler)
	http.HandleFunc("/createUser", createUserHandler)
	http.HandleFunc("/listClubs", listUsersHandler)

	http.ListenAndServe(":8080", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "base.html", nil)

}

func listUsersHandler(w http.ResponseWriter, r *http.Request) {
	// Trae los usuarios de la BD y los envía a la plantilla
	// users := getUsersFromDB() // Función que obtendría los usuarios
	// templates.ExecuteTemplate(w, "user_table.html", map[string]interface{}{"Users": users})
	templates.ExecuteTemplate(w, "user_table.html", nil)
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "modal.html", nil)
}

// Similar para listClubs, createUser, y createClub
