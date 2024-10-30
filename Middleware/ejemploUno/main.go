package main

import (
	"ejemploUno/auth"
	"ejemploUno/middleware"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/login", auth.LoginHandler)
	http.HandleFunc("/home", middleware.JWTMiddleware(auth.HomePage))

	log.Println("Servidor iniciado en :8080")
	http.ListenAndServe(":8080", nil)
}
