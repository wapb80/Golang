package repository

import (
	"postgres_internalInterface/internal/models"
)

// UserRepository define el contrato para operaciones relacionadas con usuarios.
type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByID(id int) (*models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(id int) error
	GetAllUsers() ([]*models.User, error)
}
