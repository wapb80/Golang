// Handlers - internal/handlers/user_handler.go
// Este archivo maneja las solicitudes HTTP y responde con datos en formato JSON u otros
package handlers

import (
	"encoding/json"
	"net/http"
	"postgres_internal/internal/models"
	"postgres_internal/internal/service"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(users)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	// Decodifica el cuerpo de la solicitud en un objeto User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Llama al servicio para crear el usuario
	if err := h.service.CreateUser(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Responde con éxito
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}
