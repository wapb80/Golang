package main

import (
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
	// Obtener el ID del producto
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID de producto inválido", http.StatusBadRequest)
		return
	}

	// Buscar el producto y actualizar cantidades
	for i, p := range products {
		if p.ID == id {
			if p.Quantity > 0 {
				// Disminuir cantidad disponible
				products[i].Quantity--

				// Agregar al carrito
				if item, exists := cart[id]; exists {
					item.Quantity++
					cart[id] = item
				} else {
					cart[id] = CartItem{Name: p.Name, Price: p.Price, Quantity: 1}
				}
			}
			break
		}
	}

	// Crear datos para los templates
	data := PageData{
		Products:  products,
		CartItems: getCartItems(),
		CartTotal: calculateCartTotal(),
	}

	// Respuesta fragmentada para HTMX
	RenderPartialTemplate(w, data)
}

// Renderizar la plantilla completa
func RenderTemplate(w http.ResponseWriter, templateName string, data PageData) {
	tmpl, err := template.ParseFiles("templates/layout.html", "templates/products.html", "templates/cart.html")
	if err != nil {
		http.Error(w, "Error al cargar templates", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Error al renderizar templates", http.StatusInternalServerError)
	}
}

// Renderizar solo los fragmentos actualizados
func RenderPartialTemplate(w http.ResponseWriter, data PageData) {
	tmpl, err := template.ParseFiles("templates/products.html", "templates/cart.html")
	if err != nil {
		http.Error(w, "Error al cargar templates", http.StatusInternalServerError)
		return
	}

	// Configurar encabezados HTMX
	w.Header().Set("HX-Trigger", "update-cart")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Renderizar productos
	err = tmpl.ExecuteTemplate(w, "products", data)
	if err != nil {
		http.Error(w, "Error al renderizar productos", http.StatusInternalServerError)
		return
	}

	// Renderizar carrito
	err = tmpl.ExecuteTemplate(w, "cart", data)
	if err != nil {
		http.Error(w, "Error al renderizar carrito", http.StatusInternalServerError)
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
