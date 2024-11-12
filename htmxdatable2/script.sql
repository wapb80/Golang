--- script para crear base de datos
-- Crear tabla Ciudad
CREATE TABLE Ciudad (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    nombre TEXT NOT NULL
);

-- Crear tabla Comuna
CREATE TABLE Comuna (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    nombre TEXT NOT NULL,
    ciudad_id INTEGER,
    FOREIGN KEY (ciudad_id) REFERENCES Ciudad(id) ON DELETE CASCADE
);

-- Crear tabla Club Deportivo
CREATE TABLE ClubDeportivo (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    nombre TEXT NOT NULL,
    ciudad_id INTEGER,
    comuna_id INTEGER,
    direccion TEXT,
    representante TEXT,
    FOREIGN KEY (ciudad_id) REFERENCES Ciudad(id) ON DELETE CASCADE,
    FOREIGN KEY (comuna_id) REFERENCES Comuna(id) ON DELETE CASCADE
);

-- Crear tabla Jugador
CREATE TABLE Jugador (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    rut INTEGER NOT NULL,
    nombres TEXT NOT NULL,
    apellido_paterno TEXT NOT NULL,
    apellido_materno TEXT NOT NULL,
    edad INTEGER NOT NULL,
    fecha_nacimiento DATE NOT NULL,
    ciudad_id INTEGER,
    comuna_id INTEGER,
    direccion TEXT,
    serie_juega TEXT,
    historial_deportivo TEXT,
    foto TEXT,
    club_deportivo_id INTEGER, -- Agregamos la referencia al club deportivo
    FOREIGN KEY (ciudad_id) REFERENCES Ciudad(id) ON DELETE CASCADE,
    FOREIGN KEY (comuna_id) REFERENCES Comuna(id) ON DELETE CASCADE,
    FOREIGN KEY (club_deportivo_id) REFERENCES ClubDeportivo(id) ON DELETE SET NULL -- o ON DELETE CASCADE dependiendo del comportamiento deseado
);


--- script para crear base de datos con relacion de un jugador a muchos clubs

-- Crear tabla Ciudad
CREATE TABLE Ciudad (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    nombre TEXT NOT NULL
);

-- Crear tabla Comuna
CREATE TABLE Comuna (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    nombre TEXT NOT NULL,
    ciudad_id INTEGER,
    FOREIGN KEY (ciudad_id) REFERENCES Ciudad(id) ON DELETE CASCADE
);

-- Crear tabla Club Deportivo
CREATE TABLE ClubDeportivo (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    nombre TEXT NOT NULL,
    ciudad_id INTEGER,
    comuna_id INTEGER,
    direccion TEXT,
    representante TEXT,
    FOREIGN KEY (ciudad_id) REFERENCES Ciudad(id) ON DELETE CASCADE,
    FOREIGN KEY (comuna_id) REFERENCES Comuna(id) ON DELETE CASCADE
);

-- Crear tabla Jugador
CREATE TABLE Jugador (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    rut INTEGER NOT NULL,
    nombres TEXT NOT NULL,
    apellido_paterno TEXT NOT NULL,
    apellido_materno TEXT NOT NULL,
    edad INTEGER NOT NULL,
    fecha_nacimiento DATE NOT NULL,
    ciudad_id INTEGER,
    comuna_id INTEGER,
    direccion TEXT,
    serie_juega TEXT,
    historial_deportivo TEXT,
    foto TEXT,
    FOREIGN KEY (ciudad_id) REFERENCES Ciudad(id) ON DELETE CASCADE,
    FOREIGN KEY (comuna_id) REFERENCES Comuna(id) ON DELETE CASCADE
);

-- Crear tabla para la relación entre Jugadores y Clubes Deportivos
CREATE TABLE JugadorClub (
    jugador_id INTEGER,
    club_deportivo_id INTEGER,
    PRIMARY KEY (jugador_id, club_deportivo_id),
    FOREIGN KEY (jugador_id) REFERENCES Jugador(id) ON DELETE CASCADE,
    FOREIGN KEY (club_deportivo_id) REFERENCES ClubDeportivo(id) ON DELETE CASCADE
);