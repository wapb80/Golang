package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

// Datos para las páginas
type Product struct {
	ID       int
	Name     string
	Price    float64
	Quantity int
}

type CartItem struct {
	Name     string
	Price    float64
	Quantity int
}

type PageData struct {
	Products  []Product
	CartItems []CartItem
	CartTotal float64
}

var products = []Product{
	{ID: 1, Name: "Producto A", Price: 10.0, Quantity: 10},
	{ID: 2, Name: "Producto B", Price: 20.0, Quantity: 5},
}

var cart = map[int]CartItem{}

func main() {
	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/add-to-cart", AddToCartHandler)

	log.Println("Servidor corriendo en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Manejador para la página principal
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Products:  products,
		CartItems: getCartItems(),
		CartTotal: calculateCartTotal(),
	}
	RenderTemplate(w, "layout.html", data)
}

// Manejador para agregar productos al carrito
func AddToCartHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID de producto inválido", http.StatusBadRequest)
		return
	}

	// Actualizar las cantidades de productos y el carrito
	for i, p := range products {
		if p.ID == id && p.Quantity > 0 {
			products[i].Quantity-- // Disminuir cantidad en el producto
			if item, exists := cart[id]; exists {
				item.Quantity++ // Aumentar cantidad en el carrito
				cart[id] = item
			} else {
				cart[id] = CartItem{Name: p.Name, Price: p.Price, Quantity: 1}
			}
			break
		}
	}

	// Preparar los datos para los templates
	data := PageData{
		Products:  products,
		CartItems: getCartItems(),
		CartTotal: calculateCartTotal(),
	}

	// Solo renderizamos las secciones relevantes
	if r.Header.Get("HX-Request") == "true" {
		// Si es una petición HTMX, solo renderizamos las partes necesarias
		renderPartialTemplate(w, data, "products")
		renderPartialTemplate(w, data, "cart")
	} else {
		// Si no es una petición HTMX, renderizamos todo
		RenderTemplate(w, "layout.html", data)
	}
}

// Renderizar la plantilla completa
// Renderizar la plantilla completa
func RenderTemplate(w http.ResponseWriter, templateName string, data PageData) {
	tmpl, err := template.ParseFiles("templates/layout.html", "templates/products.html", "templates/cart.html")
	if err != nil {
		log.Printf("Error al cargar los templates: %v", err) // Imprimir el error
		http.Error(w, fmt.Sprintf("Error al cargar templates: %v", err), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Error al renderizar el template: %v", err) // Imprimir el error
		http.Error(w, fmt.Sprintf("Error al renderizar templates: %v", err), http.StatusInternalServerError)
	}
}

// Renderizar solo los fragmentos actualizados
func renderPartialTemplate(w http.ResponseWriter, data PageData, target string) {
	tmpl, err := template.ParseFiles("templates/" + target + ".html")
	if err != nil {
		log.Printf("Error al cargar el template de fragmento: %v", err) // Imprimir el error
		http.Error(w, fmt.Sprintf("Error al cargar template de fragmento: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Error al renderizar el fragmento: %v", err) // Imprimir el error
		http.Error(w, fmt.Sprintf("Error al renderizar el fragmento: %v", err), http.StatusInternalServerError)
	}
}

func getCartItems() []CartItem {
	var items []CartItem
	for _, item := range cart {
		items = append(items, item)
	}
	return items
}

func calculateCartTotal() float64 {
	var total float64
	for _, item := range cart {
		total += item.Price * float64(item.Quantity)
	}
	return total
}
