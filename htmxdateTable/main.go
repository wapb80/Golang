package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Conectar a SQLite
	db, err := sql.Open("sqlite3", "./example.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Crear una tabla
	createTableSQL := `CREATE TABLE IF NOT EXISTS users (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"name" TEXT
	);`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Error al crear la tabla: %v\n", err)
	}

	// Insertar datos
	insertUserSQL := `INSERT INTO users (name) VALUES (?)`
	_, err = db.Exec(insertUserSQL, "Juan Pérez")
	if err != nil {
		log.Fatalf("Error al insertar datos: %v\n", err)
	}

	// Consultar datos
	rows, err := db.Query("SELECT id, name FROM users")
	if err != nil {
		log.Fatalf("Error al consultar datos: %v\n", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		err := rows.Scan(&id, &name)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("ID: %d, Nombre: %s\n", id, name)
	}

	// Verificar si hubo errores en las filas
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}
