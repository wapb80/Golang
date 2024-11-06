// main.go ,, no me funcionaba al cerrarr el modal , debia presionar un click fuera para que volviera a toimar el control
package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
)

type User struct {
	ID    int
	Name  string
	Email sql.NullString
}

type Club struct {
	ID   int
	Name string
}

var db *sql.DB
var templates *template.Template

func main() {
	r := http.NewServeMux()
	db = InitDB()
	defer db.Close()

	templates = template.Must(template.ParseGlob("templates/*.html"))

	r.HandleFunc("/", HomeHandler)
	// http.HandleFunc("/listClubs", listClubesHandler)
	r.HandleFunc("/listUsers", listUsersHandler)
	r.HandleFunc("/createUser", templateCreateUserHandler)
	// http.HandleFunc("/users", ListUsersHandler)
	// http.HandleFunc("/clubs", ListClubsHandler)
	r.HandleFunc("/users/create", CreateUserHandler)
	// http.HandleFunc("/clubs/create", CreateClubHandler)
	// http.HandleFunc("/clubs/insert", insertClubHandler)
	r.HandleFunc("/user/delete/{id}", DeleteUserHandler)
	// http.HandleFunc("/clubs/delete/", DeleteClubHandler)
	r.HandleFunc("GET /users/edit/{id}", EditUserHandler)
	r.HandleFunc("POST /users/edit/{id}", EditUserHandlerPost)
	// http.HandleFunc("/clubs/edit/", EditClubHandler)

	log.Println("Servidor iniciado en :8080")
	http.ListenAndServe(":8080", r)
}

// Listar Usuarios y Clubes
// func listUsersHandler(w http.ResponseWriter, r *http.Request) {
// 	// Trae los usuarios de la BD y los envía a la plantilla
// 	// users := getUsersFromDB() // Función que obtendría los usuarios
// 	// templates.ExecuteTemplate(w, "user_table.html", map[string]interface{}{"Users": users})
// 	templates.ExecuteTemplate(w, "user_table.html", nil)
// }

// func listClubesHandler(w http.ResponseWriter, r *http.Request) {
// 	// Trae los usuarios de la BD y los envía a la plantilla
// 	// users := getUsersFromDB() // Función que obtendría los usuarios
// 	// templates.ExecuteTemplate(w, "user_table.html", map[string]interface{}{"Users": users})
// 	templates.ExecuteTemplate(w, "club_table.html", nil)
// }

func listUsersHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, name, email FROM users")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}
	// println(len(users))
	templates.ExecuteTemplate(w, "user_table.html", users)
}

// func ListClubsHandler(w http.ResponseWriter, r *http.Request) {
// 	rows, err := db.Query("SELECT id, name FROM clubs")
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	defer rows.Close()

// 	var clubs []Club
// 	for rows.Next() {
// 		var club Club
// 		if err := rows.Scan(&club.ID, &club.Name); err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 		clubs = append(clubs, club)
// 	}
// 	templates.ExecuteTemplate(w, "clubs.html", clubs)
// }

// // Crear USER
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {

	name := r.FormValue("name")
	email := r.FormValue("email")
	_, err := db.Exec("INSERT INTO users (name, email) VALUES (?, ?)", name, email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/listUsers", http.StatusSeeOther)
}

func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	_, err := db.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/listUsers", http.StatusSeeOther)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "base.html", nil)
}

func templateCreateUserHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "create_user.html", nil)
}

// func insertClubHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.Method == http.MethodPost {
// 		// name := r.FormValue("name")
// 		// representante := r.FormValue("club-representante")

// 		// _, err := db.Exec("INSERT INTO clubs (name) VALUES (?)", name)
// 		// if err != nil {
// 		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
// 		// 	return
// 		// }

// 		http.Redirect(w, r, "/clubs", http.StatusSeeOther)
// 		return
// 	}

// }

// func CreateClubHandler(w http.ResponseWriter, r *http.Request) {

// 	tmpl, err := template.ParseFiles("templates/create_club.html")

// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	tmpl.Execute(w, nil)
// 	// tmpl.ExecuteTemplate(w)
// 	// name := r.FormValue("name")
// 	// _, err := db.Exec("INSERT INTO clubs (name) VALUES (?)", name)
// 	// if err != nil {
// 	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
// 	// 	return
// 	// }
// 	// http.Redirect(w, r, "/clubs", http.StatusSeeOther)
// 	// return

// }

// func DeleteClubHandler(w http.ResponseWriter, r *http.Request) {
// 	id := r.URL.Path[len("/clubs/delete/"):]
// 	_, err := db.Exec("DELETE FROM clubs WHERE id = ?", id)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	http.Redirect(w, r, "/clubs", http.StatusSeeOther)
// }

func EditUserHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	// if r.Method == http.MethodPost {
	// 	id := r.URL.Path[len("/clubs/edit/"):]
	// Obtener datos actuales del usuario
	row := db.QueryRow("SELECT id, name, email FROM users WHERE id = ?", id)
	var user User
	err := row.Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	templates.ExecuteTemplate(w, "edit_user.html", user)

	// }

}

func EditUserHandlerPost(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	// if r.Method == http.MethodPost {
	name := r.FormValue("name")
	email := r.FormValue("email")
	_, err := db.Exec("UPDATE users SET name = ?,email = ? WHERE id = ?", name, email, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/listUsers", http.StatusSeeOther)

}

// 	// Obtener datos actuales del club
// 	row := db.QueryRow("SELECT id, name FROM clubs WHERE id = ?", id)
// 	var club Club
// 	err := row.Scan(&club.ID, &club.Name)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	templates.ExecuteTemplate(w, "edit_club.html", club)
// }
