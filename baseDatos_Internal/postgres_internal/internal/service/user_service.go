// Service - internal/service/user_service.go
// Este archivo contiene la lógica empresarial y valida los datos.
// El servicio recibe la solicitud desde el controlador:
package service

import (
	"errors"
	"postgres_internal/internal/models"
	"postgres_internal/internal/repository"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUsers() ([]models.User, error) {
	return s.repo.GetAll()
}

func (s *UserService) CreateUser(user models.User) error {
	// Validaciones adicionales
	if user.Name == "" || user.Email == "" {
		return errors.New("nombre y correo electrónico son obligatorios")
	}
	return s.repo.Create(user)
}
