package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
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
		// Genera datos dinámicos para 20 selects
		menuData := []map[string]interface{}{
			{"Label": "Nacionalidad", "Options": []string{"Chilena", "Argentina", "Peruana", "Colombiana"}},
			{"Label": "Género", "Options": []string{"Masculino", "Femenino", "Otro"}},
			{"Label": "Comida Favorita", "Options": []string{"Pizza", "Empanadas", "Ceviche", "Tacos"}},
		}

		// Agrega más opciones
		for i := 4; i <= 20; i++ {
			menuData = append(menuData, map[string]interface{}{
				"Label":   "Opción " + strconv.Itoa(i),
				"Options": []string{"Opción A", "Opción B", "Opción C"},
			})
		}

		renderTemplate(w, "menu_reportes.html", map[string]interface{}{
			"Options": menuData,
		})
	})

	log.Println("Servidor iniciado en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
