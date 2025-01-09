package main

import (
	"log"
	"net/http"
)

func principalHandler(w http.ResponseWriter, r *http.Request) {
	// esto es para que la pagina principal no acepto otras paginas
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Write([]byte("Principal"))
}

func menuHandler(w http.ResponseWriter, r *http.Request) {
	// cuando la url es :: /menu?page=miPagina
	page := r.URL.Query().Get("page")
	println(page)
	w.Write([]byte("/Menu rescatando la query es : " + page))
}

func menusHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ruta subtree"))
}

func usuariosHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ruta Dinamica"))
}
func usuarioHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	w.Write([]byte("/ruta dinamica , el id es : " + id))
}
func main() {

	serv := http.NewServeMux()
	serv.HandleFunc("GET /", principalHandler)
	serv.HandleFunc("GET /menu", menuHandler)
	serv.HandleFunc("GET /menu/", menusHandler)
	serv.HandleFunc("GET /usuario", usuariosHandler)
	serv.HandleFunc("PUT /usuario/{id}", usuarioHandler)

	log.Println("Servidor iniciado en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", serv))
}
