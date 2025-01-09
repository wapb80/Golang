package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv" // Carga el archivo .env
	_ "github.com/lib/pq"      // Driver para PostgreSQL
)

// InitDB inicializa la conexión a la base de datos
func InitDB() (*sql.DB, error) {
	// Carga las variables de entorno desde el archivo .env
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error cargando el archivo .env: %v", err)
	}

	// Lee las variables de entorno
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbSSLMode := os.Getenv("DB_SSLMODE")

	// Construye el DSN (Data Source Name)
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode,
	)

	// Abre la conexión
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("error abriendo la conexión: %w", err)
	}

	// Verifica la conexión
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("error verificando la conexión: %w", err)
	}

	fmt.Println("Conectado exitosamente a la base de datos")
	return db, nil
}
