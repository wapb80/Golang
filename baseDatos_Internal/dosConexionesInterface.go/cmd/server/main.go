package main

import (
	"log"
	"net/http"

	"dosConexionesInterface/config"
	"dosConexionesInterface/internal/handlers"
	"dosConexionesInterface/internal/repository"
	"dosConexionesInterface/internal/service"

	_ "github.com/lib/pq"
)

// const (
// 	dbHost     = "localhost"
// 	dbPort     = 5432
// 	dbUser     = "user"
// 	dbPassword = "password"
// 	dbName     = "mydb"
// )

func main() {
	// Database connection
	// Inicializa la conexión a la base de datos
	dbPostgres, err := config.InitDBPostgres()
	if err != nil {
		log.Fatalf("No se pudo conectar a la base de datos Postgres: %v", err)
	}
	defer dbPostgres.Close()

	// Inicializa la conexión a la base de datos
	dbMysql, err := config.InitDBMysql()
	if err != nil {
		log.Fatalf("No se pudo conectar a la base de datos Mysql: %v", err)
	}
	defer dbMysql.Close()

	// Dependency Injection
	//userRepo := repository.NewUserRepository(db)
	userRepopostgres := repository.NewPostgresUserRepository(dbPostgres)
	userServicepostgres := service.NewUserService(userRepopostgres)
	userHandlerpostgres := handlers.NewUserHandler(userServicepostgres)

	// Dependency Injection
	//userRepo := repository.NewUserRepository(db)
	userRepomysql := repository.NewPostgresUserRepository(dbMysql)
	userServicemysql := service.NewUserService(userRepomysql)
	userHandlermysql := handlers.NewUserHandler(userServicemysql)

	// Router setup Para postgres
	r := http.NewServeMux()
	r.HandleFunc("GET /usersPostgres", userHandlerpostgres.GetAllUsers)
	r.HandleFunc("POST /usersPostgres", userHandlerpostgres.CreateUser)

	// Router setup para mysql

	r.HandleFunc("GET /usersMysql", userHandlermysql.GetAllUsers)
	r.HandleFunc("POST /usersMysql", userHandlermysql.CreateUser)

	// r := mux.NewRouter()
	// r.HandleFunc("/users", userHandler.GetUsers).Methods("GET")
	// r.HandleFunc("/users", userHandler.CreateUser).Methods("POST")

	log.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", r)
}
