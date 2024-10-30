package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var secretKey = []byte("your_secret_key")

// Estructura para el usuario
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Generar JWT
func generateJWT(username string) (string, error) {
	claims := &jwt.StandardClaims{
		Subject:   username,
		ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// Validar JWT
func validateJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})
}

// Manejar la página principal
func homeHandler(w http.ResponseWriter, r *http.Request) {
	html := `
    <!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Autenticación con HTMX</title>
    <script src="https://unpkg.com/htmx.org"></script>
</head>
<body>
    <h1>Inicio de Sesión</h1>
    <form hx-post="/login" hx-target="#response" hx-ext="json-enc">
        <input type="text" name="username" placeholder="Usuario" required>
        <input type="password" name="password" placeholder="Contraseña" required>
        <button type="submit">Iniciar Sesión</button>
    </form>
    
    <div id="response"></div>

    <script src="https://unpkg.com/htmx.org/dist/ext/json-enc.js"></script>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

// Manejar el inicio de sesión
func loginHandler(w http.ResponseWriter, r *http.Request) {
	var user User

	err := json.NewDecoder(r.Body).Decode(&user)
	log.Printf("Error reading body: %v", err)
	log.Printf("Usuario: %s, Contraseña: %s\n", user.Username, user.Password)
	if err != nil || user.Username == "" || user.Password == "" {
		http.Error(w, "Credenciales inválidas", http.StatusBadRequest)
		return
	}

	// Aquí puedes validar las credenciales del usuario (ejemplo simplificado)
	if user.Username == "admin" && user.Password == "password" {
		tokenString, err := generateJWT(user.Username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Establecer el token en una cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    tokenString,
			Path:     "/",
			HttpOnly: true,
			Secure:   false,
			MaxAge:   3600,
		})

		w.Write([]byte("Inicio de sesión exitoso."))
		return
	}

	http.Error(w, "Credenciales incorrectas", http.StatusUnauthorized)
}

// Manejar el cierre de sesión
func logoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})

	w.Write([]byte("Has cerrado sesión exitosamente."))
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)

	fmt.Println("Servidor escuchando en :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
