package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql" // Driver para Mysql
	"github.com/joho/godotenv"         // Carga el archivo .env
	_ "github.com/lib/pq"              // Driver para PostgreSQL
)

// InitDB inicializa la conexión a la base de datos
func InitDBPostgres() (*sql.DB, error) {
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
		return nil, fmt.Errorf("error abriendo la conexión Postgres: %w", err)
	}

	// Verifica la conexión
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("error verificando la conexión Postgres: %w", err)
	}

	fmt.Println("Conectado exitosamente a la base de datos Postgres")
	return db, nil
}

// InitDB inicializa la conexión a la base de datos
func InitDBMysql() (*sql.DB, error) {
	// Carga las variables de entorno desde el archivo .env
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error cargando el archivo .env: %v", err)
	}

	// Lee las variables de entorno
	dbHost := os.Getenv("MYSQL_HOST")
	dbPort := os.Getenv("MYSQL_PORT")
	dbUser := os.Getenv("MYSQL_USER")
	dbPassword := os.Getenv("MYSQL_PASSWORD")
	dbName := os.Getenv("MYSQL_NAME")
	// dbSSLMode := os.Getenv("MYSQL_SSLMODE")

	// Construye el DSN (Data Source Name)
	// dsn := fmt.Sprintf(
	// 	"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
	// 	dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode,
	// )
	dsn := dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName

	// Abre la conexión a Mysql
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("error abriendo la conexión Mysql: %w", err)
	}

	// Verifica la conexión
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("error verificando la conexión MYSQL: %w", err)
	}

	fmt.Println("Conectado exitosamente a la base de datos Mysql")
	return db, nil
}
