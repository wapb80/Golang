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

	//-- Crear tabla Club Deportivo
	createTablesSQL := `
	CREATE TABLE IF NOT EXISTS ClubDeportivo (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		nombre TEXT,
		comuna TEXT,
		direccion TEXT,
		representante TEXT,
		activo INTEGER
	);

	--- Crear tabla Jugador
	CREATE TABLE IF NOT EXISTS Jugador (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		rut INTEGER ,
		dv TEXT,
		nombres TEXT,
		apellido_paterno TEXT ,
		apellido_materno TEXT ,
		edad INTEGER NOT NULL,
		fecha_nacimiento DATE,
		comuna TEXT,
		direccion TEXT,
		serie_juega TEXT,
		historial TEXT,
		foto TEXT,
	    mail TEXT,
		club_juega, -- Agregamos la referencia al club deportivo
		activo INTEGER
		 -- FOREIGN KEY (club_deportivo_id) REFERENCES ClubDeportivo(id) ON DELETE SET NULL -- o ON DELETE CASCADE dependiendo del comportamiento deseado
	);

	CREATE TABLE IF NOT EXISTS Comuna (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			nombre TEXT NOT NULL
		);


	CREATE TABLE IF NOT EXISTS Series (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			nombre TEXT NOT NULL
	);



	`
	// delete from 'Comuna';
	// update sqlite_sequence set seq=0 where name='comuna';

	// INSERT INTO Series (nombre) VALUES
	// 	('Infantil'),
	// 	('Adulta'),
	// 	('Senior');

	// 	INSERT INTO Comuna (nombre) VALUES
	// 	('Cerrillos'),
	// 	('Cerro Navia'),
	// 	('Conchalí'),
	// 	('El Bosque'),
	// 	('Estación Central'),
	// 	('Huechuraba'),
	// 	('Independencia'),
	// 	('La Cisterna'),
	// 	('La Florida'),
	// 	('La Granja'),
	// 	('La Pintana'),
	// 	('La Reina'),
	// 	('Las Condes'),
	// 	('Lo Barnechea'),
	// 	('Lo Espejo'),
	// 	('Lo Prado'),
	// 	('Macul'),
	// 	('Maipú'),
	// 	('Ñuñoa'),
	// 	('Pedro Aguirre Cerda'),
	// 	('Peñalolén'),
	// 	('Providencia'),
	// 	('Pudahuel'),
	// 	('Quilicura'),
	// 	('Quinta Normal'),
	// 	('Recoleta'),
	// 	('Renca'),
	// 	('San Joaquín'),
	// 	('San Miguel'),
	// 	('San Ramón'),
	// 	('Santiago'),
	// 	('Vitacura');
	// INSERT INTO ClubDeportivo (nombre, comuna, direccion, representante) VALUES
	// ('Deportivo Los Andes', 'Providencia', 'Calle Los Leones 456', 'María López'),
	// ('Club Futuro Estrella', 'La Florida', 'Av. La Florida 789', 'Carlos Rojas'),
	// ('Academia de Fútbol La Cisterna', 'La Cisterna', 'Calle Central 1011', 'Ana González'),
	// ('Unión Maipú', 'Maipú', 'Av. Pajaritos 1213', 'Luis Fernández');
	// // Creación de tablas
	// createTablesSQL := `
	// CREATE TABLE IF NOT EXISTS users (
	// 	id INTEGER PRIMARY KEY AUTOINCREMENT,
	// 	name TEXT,
	// 	email TEXT
	// );
	// CREATE TABLE IF NOT EXISTS clubs (
	// 	id INTEGER PRIMARY KEY AUTOINCREMENT,
	// 	name TEXT
	// );
	// CREATE TABLE IF NOT EXISTS user_club (
	// 	user_id INTEGER,
	// 	club_id INTEGER,
	// 	FOREIGN KEY (user_id) REFERENCES users (id),
	// 	FOREIGN KEY (club_id) REFERENCES clubs (id),
	// 	PRIMARY KEY (user_id, club_id)
	// );`

	_, err = db.Exec(createTablesSQL)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

// GetDB devuelve la instancia de la base de datos
// func GetDB() *sql.DB {
// 	return db
// }
