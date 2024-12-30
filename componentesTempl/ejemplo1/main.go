package main

import (
	"context"
	"ejemplo1/models"    // Importa el modelo
	"ejemplo1/templates" // Importa las plantillas generadas
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Crear el objeto usuario usando el modelo compartido
		user := models.User{
			Name:  "Juan Pérez",
			Email: "juan.perez@example.com",
			Age:   30,
		}

		// Crear el componente del usuario
		userComponent := templates.UserComponent(user)

		// Renderizar el componente
		err := userComponent.Render(context.Background(), w)
		if err != nil {
			http.Error(w, "Error al renderizar la plantilla", http.StatusInternalServerError)
			log.Println("Error:", err)
		}
	})

	log.Println("Servidor iniciado en http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
