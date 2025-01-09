package main

import (
	"log"
	"net/http"

	"postgres_internalInterface/config"
	"postgres_internalInterface/internal/handlers"
	"postgres_internalInterface/internal/repository"
	"postgres_internalInterface/internal/service"

	_ "github.com/lib/pq"
)

const (
	dbHost     = "localhost"
	dbPort     = 5432
	dbUser     = "user"
	dbPassword = "password"
	dbName     = "mydb"
)

func main() {
	// Database connection
	// Inicializa la conexión a la base de datos
	db, err := config.InitDB()
	if err != nil {
		log.Fatalf("No se pudo conectar a la base de datos: %v", err)
	}
	defer db.Close()

	// Dependency Injection
	//userRepo := repository.NewUserRepository(db)
	userRepo := repository.NewPostgresUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	// Router setup
	r := http.NewServeMux()
	r.HandleFunc("GET /users", userHandler.GetAllUsers)
	r.HandleFunc("POST /users", userHandler.CreateUser)

	// r := mux.NewRouter()
	// r.HandleFunc("/users", userHandler.GetUsers).Methods("GET")
	// r.HandleFunc("/users", userHandler.CreateUser).Methods("POST")

	log.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", r)
}
