package main

import (
	"log"
	"net/http"
)

func menuHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Menu"))
}
func principalHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Principal"))
}

func main() {

	http.HandleFunc("/", principalHandler)
	http.HandleFunc("/menu/", menuHandler)

	log.Println("Servidor iniciado en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
