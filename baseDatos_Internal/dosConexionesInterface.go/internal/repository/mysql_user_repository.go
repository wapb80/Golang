package repository

import (
	"database/sql"
	"dosConexionesInterface/internal/models"
)

// MySQLUserRepository implementa la interfaz UserRepository para MySQL.
type MySQLUserRepository struct {
	db *sql.DB
}

// NewMySQLUserRepository crea una nueva instancia del repositorio para MySQL.
func NewMySQLUserRepository(db *sql.DB) *MySQLUserRepository {
	return &MySQLUserRepository{db: db}
}

// CreateUser crea un nuevo usuario en MySQL.
func (r *MySQLUserRepository) CreateUser(user *models.User) error {
	_, err := r.db.Exec("INSERT INTO users (name, email) VALUES (?, ?)", user.Name, user.Email)
	return err
}

// GetUserByID obtiene un usuario por ID en MySQL.
func (r *MySQLUserRepository) GetUserByID(id int) (*models.User, error) {
	var user models.User
	err := r.db.QueryRow("SELECT id, name, email FROM users WHERE id = ?", id).
		Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser actualiza un usuario en MySQL.
func (r *MySQLUserRepository) UpdateUser(user *models.User) error {
	_, err := r.db.Exec("UPDATE users SET name = ?, email = ? WHERE id = ?", user.Name, user.Email, user.ID)
	return err
}

// DeleteUser elimina un usuario en MySQL.
func (r *MySQLUserRepository) DeleteUser(id int) error {
	_, err := r.db.Exec("DELETE FROM users WHERE id = ?", id)
	return err
}

// GetAllUsers obtiene todos los usuarios en MySQL.
func (r *MySQLUserRepository) GetAllUsers() ([]*models.User, error) {
	rows, err := r.db.Query("SELECT id, name, email FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}
