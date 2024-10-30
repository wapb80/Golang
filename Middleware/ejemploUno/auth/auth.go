package auth

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte("secret_key") // Cambia esto a una clave segura

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginHandler sirve el formulario de login y maneja la autenticación
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl := template.Must(template.ParseFiles("templates/login2.html"))
		tmpl.Execute(w, nil)
		return
	}

	if r.Method == http.MethodPost {

		var creds Credentials
		if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
			//log.Printf("Error decoding JSON: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		fmt.Println("Hello, World!")
		if creds.Username == "user" && creds.Password == "pass" {
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"username": creds.Username,
				"exp":      time.Now().Add(time.Minute * 2).Unix(),
			})
			tokenString, err := token.SignedString(jwtKey)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			http.SetCookie(w, &http.Cookie{
				Name:     "token",
				Value:    tokenString,
				Expires:  time.Now().Add(time.Hour * 24),
				HttpOnly: true,
			})

			w.Write([]byte("Autenticación exitosa"))
			w.Write([]byte(tokenString))
			return
		}

		w.WriteHeader(http.StatusUnauthorized)
	}
}

// HomePage muestra la página protegida
func HomePage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/home.html"))
	tmpl.Execute(w, nil)
}
