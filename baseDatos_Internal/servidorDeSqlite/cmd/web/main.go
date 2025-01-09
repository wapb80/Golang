package main

import (
	"log"
	"net/http"
)

func main() {
	serv := http.Server{
		Addr:    ":8000",
		Handler: routes(),
	}
	log.Println("escuchando on :8000")
	serv.ListenAndServe()

}
