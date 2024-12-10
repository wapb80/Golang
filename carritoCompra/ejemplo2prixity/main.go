package main

import (
	"html/template"
	"net/http"
	"strconv"
)

// Estructuras para productos y carrito
type Product struct {
	ID    int
	Name  string
	Price float64
}

type CartItem struct {
	Name  string
	Price float64
}

type PageVariables struct {
	Products []Product
	Cart     []CartItem
	Total    float64
}

// Variables globales para simular el carrito y los productos disponibles
var cart []CartItem
var products = []Product{
	{ID: 1, Name: "Producto A", Price: 10.0},
	{ID: 2, Name: "Producto B", Price: 15.0},
	{ID: 3, Name: "Producto C", Price: 20.0},
}

func main() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/add-to-cart/", addToCart)

	// Iniciar el servidor en el puerto 8080
	http.ListenAndServe(":8080", nil)
}

func homePage(w http.ResponseWriter, r *http.Request) {
	renderPage(w)
}

func addToCart(w http.ResponseWriter, r *http.Request) {
	// Obtener el ID del producto desde la URL
	productIDStr := r.URL.Path[len("/add-to-cart/"):]

	// Convertir productIDStr a un entero
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	// Buscar el producto correspondiente y agregarlo al carrito
	for _, product := range products {
		if product.ID == productID {
			cart = append(cart, CartItem{Name: product.Name, Price: product.Price})
			break
		}
	}

	// Renderizar la página después de agregar el producto al carrito
	renderPage(w)
}

// Función para renderizar la página con los productos y el carrito
func renderPage(w http.ResponseWriter) {
	vars := PageVariables{
		Products: products,
		Cart:     cart,
		Total:    calculateTotal(cart),
	}

	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/content1.html", "templates/content2.html"))

	w.Header().Set("Content-Type", "text/html")

	// Renderizar la plantilla base con el contenido inicial o actualizado
	tmpl.ExecuteTemplate(w, "layout.html", vars)
}

func calculateTotal(cart []CartItem) float64 {
	total := 0.0
	for _, item := range cart {
		total += item.Price
	}
	return total
}
