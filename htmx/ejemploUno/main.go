package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

// Datos básicos para simular una API
type Item struct {
	ID    int
	Title string
}

var items = []Item{
	{ID: 1, Title: "Item 1"},
	{ID: 2, Title: "Item 2"},
	{ID: 3, Title: "Item 3"},
}

func main() {
	http.HandleFunc("/", renderHomePage)
	http.HandleFunc("/api/items", getItemsAPI)
	http.HandleFunc("/item", getItemPage)

	// Servir el directorio estático para HTMX
	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// renderHomePage maneja la página inicial
func renderHomePage(w http.ResponseWriter, r *http.Request) {
	// tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/content.html"))

	tmpl, err := template.ParseFiles("templates/layout.html", "templates/content.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	tmpl.ExecuteTemplate(w, "layout", items)
	// tmpl.Execute(w, "layout")
}

// getItemsAPI responde con datos JSON
func getItemsAPI(w http.ResponseWriter, r *http.Request) {

	idParam := r.URL.Query().Get("id")
	if idParam != "" {
		fmt.Println("Hello, World!")
		id, _ := strconv.Atoi(idParam)
		for _, item := range items {
			if item.ID == id {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(item)
				return
			}
		}
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

// getItemPage responde con HTML parcial para un solo item
func getItemPage(w http.ResponseWriter, r *http.Request) {

	idParam := r.URL.Query().Get("id")

	id, _ := strconv.Atoi(idParam)
	var selectedItem *Item
	for _, item := range items {
		if item.ID == id {
			selectedItem = &item
			break
		}
	}
	if selectedItem == nil {
		http.NotFound(w, r)
		return
	}

	tmpl, err := template.ParseFiles("templates/content.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// tmpl.Execute(w, selectedItem)
	tmpl.ExecuteTemplate(w, "content", selectedItem)
}
