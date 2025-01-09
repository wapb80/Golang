package service

import (
	"dosConexionesInterface/internal/models"
	"dosConexionesInterface/internal/repository"
	"errors"
)

// UserService define la estructura del servicio.
type UserService struct {
	repo repository.UserRepository // Usamos la interfaz , para definir la estructura del servico
}

// NewUserService crea un nuevo servicio.
// actúa como un constructor para crear e inicializar una instancia del tipo UserService
func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// Métodos definidos
// CreateUser valida y pasa los datos al repositorio.
func (s *UserService) CreateUser(user *models.User) error {
	if user.Name == "" || user.Email == "" {
		return errors.New("nombre y correo electrónico son obligatorios")
	}
	return s.repo.CreateUser(user)
}

func (s *UserService) GetUserByID(id int) (*models.User, error) {
	return s.repo.GetUserByID(id)
}

func (s *UserService) UpdateUser(user *models.User) error {
	return s.repo.UpdateUser(user)
}

func (s *UserService) DeleteUser(id int) error {
	return s.repo.DeleteUser(id)
}

func (s *UserService) GetAllUsers() ([]*models.User, error) {
	return s.repo.GetAllUsers()
}
