package main

import (
	"html/template"
	"net/http"
)

func getHome(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./assets/template/home.page.html")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	t.Execute(w, nil)
}
