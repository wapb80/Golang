// database.go
package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB() *sql.DB {
	db, err := sql.Open("sqlite3", "./sports_club.db")
	if err != nil {
		log.Fatal(err)
	}

	// Creación de tablas
	createTablesSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		email TEXT
	);
	CREATE TABLE IF NOT EXISTS clubs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT
	);
	CREATE TABLE IF NOT EXISTS user_club (
		user_id INTEGER,
		club_id INTEGER,
		FOREIGN KEY (user_id) REFERENCES users (id),
		FOREIGN KEY (club_id) REFERENCES clubs (id),
		PRIMARY KEY (user_id, club_id)
	);`

	_, err = db.Exec(createTablesSQL)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

// GetDB devuelve la instancia de la base de datos
func GetDB() *sql.DB {
	return db
}
