package models

type User struct {
	ID    int    `json:"id"`    // ID del usuario
	Name  string `json:"name"`  // Nombre del usuario
	Email string `json:"email"` // Correo electrónico
}
