// main.go ,, no me funcionaba al cerrarr el modal , debia presionar un click fuera para que volviera a toimar el control
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type Jugador struct {
	ID               int
	Mail             sql.NullString
	Rut              int
	Dv               string
	Nombres          string
	Apellido_paterno string
	Apellido_materno string
	Club_juega       string
	Serie_juega      string
	Foto             string
	Edad             int
}

type Club struct {
	ID     int
	Nombre string
}

var db *sql.DB
var templates *template.Template

const uploadPath = "./uploads/" // Directorio donde se guardarán las imágenes

func main() {
	// Crear el directorio de subida si no existe
	os.MkdirAll(uploadPath, os.ModePerm)

	r := http.NewServeMux()
	db = InitDB()
	defer db.Close()

	templates = template.Must(template.ParseGlob("templates/*.html"))
	// Servir el directorio de imágenes
	r.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir(uploadPath))))

	r.HandleFunc("GET /", HomeHandler)
	// http.HandleFunc("/listClubs", listClubesHandler)
	r.HandleFunc("GET /listUsers", listUsersHandler)
	r.HandleFunc("GET /createUser", templateCreateUserHandler)
	// http.HandleFunc("/users", ListUsersHandler)
	// http.HandleFunc("/clubs", ListClubsHandler)
	r.HandleFunc("POST /users/create", CreateUserHandler)
	// http.HandleFunc("/clubs/create", CreateClubHandler)
	// http.HandleFunc("/clubs/insert", insertClubHandler)
	r.HandleFunc("GET /user/delete/{id}", DeleteUserHandler)
	r.HandleFunc("GET /user/find/{id}", FindUserHandler)
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
	rows, err := db.Query("SELECT id,rut,dv,nombres,apellido_paterno,apellido_materno,club_juega,foto,serie_juega,edad FROM jugador")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var Jugadores []Jugador
	for rows.Next() {
		var jugador Jugador
		if err := rows.Scan(&jugador.ID, &jugador.Rut, &jugador.Dv, &jugador.Nombres, &jugador.Apellido_paterno, &jugador.Apellido_materno, &jugador.Club_juega, &jugador.Foto, &jugador.Serie_juega, &jugador.Edad); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		Jugadores = append(Jugadores, jugador)
	}
	// println(len(users))
	templates.ExecuteTemplate(w, "user_table.html", Jugadores)
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
	// Parsear el formulario multipart

	err := r.ParseMultipartForm(10 << 20) // Limitar a 10 MB
	if err != nil {
		http.Error(w, "Error al parsear el formulario: "+err.Error(), http.StatusBadRequest)
		return
	}
	file, _, err := r.FormFile("foto")
	if err != nil {
		http.Error(w, "Error al obtener el archivo: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	//********************************************************************** Cambiar nombre de la imagen
	// Obtener la fecha y hora actuales
	now := time.Now()
	// Formatear la fecha y hora en el formato deseado: día mes año hora minutos sin separadores
	formattedDateTime := now.Format("020120061504")

	// Definir un entero que deseas concatenar
	number := r.FormValue("rut2")

	// Concatenar la fecha y hora con el entero (convertido a cadena)
	nombreFoto := fmt.Sprintf("%s_%s", number, formattedDateTime)

	//*********************************************************************************** Cambiar nombre de la imagen
	out, err := os.Create(uploadPath + nombreFoto + ".jpg") // Cambia el nombre según sea necesario
	if err != nil {
		http.Error(w, "Error al crear el archivo: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		http.Error(w, "Error al guardar el archivo: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// fmt.Fprintf(w, "Imagen subida exitosamente: %s", uploadPath+"uploaded_image.jpg")
	// esto esta bien sin imagenes
	// name := r.FormValue("name")
	// email := r.FormValue("email")
	// _, err = db.Exec("INSERT INTO users (name, email) VALUES (?, ?)", name, email)

	_, err = db.Exec("INSERT INTO jugador (rut,dv,nombres,apellido_paterno,apellido_materno,mail,edad,fecha_nacimiento,comuna,direccion,club_juega,serie_juega,historial,activo,foto) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ? ,?)", r.FormValue("rut2"), r.FormValue("dv"), r.FormValue("nombres"), r.FormValue("apellido_paterno"), r.FormValue("apellido_materno"), r.FormValue("email"), r.FormValue("edad"), r.FormValue("fecha_nacimiento"), r.FormValue("comuna"), r.FormValue("direccion"), r.FormValue("club_juega"), r.FormValue("serie_juega"), r.FormValue("historial"), 1, nombreFoto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Usuario Ingeresado exitosamente ")
	//templates.ExecuteTemplate(w, "usuarioCreado.html", nil)
	// http.Redirect(w, r, "/", http.StatusSeeOther)
	//listUsersHandler(w, r)
}
func FindUserHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	// println("select rut FROM jugador WHERE activo=1 and  rut = ?", id)
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM jugador WHERE activo=1 and  rut = ?)", id).Scan(&exists)
	// err := db.QueryRow("select rut FROM jugador WHERE activo=1 and  rut = ?", id).Scan(&exists)
	if err != nil {
		http.Error(w, "Error querying the database", http.StatusInternalServerError)
		log.Println("Database query error:", err)
		return
	}

	// Responde con el resultado en formato JSON
	response := map[string]bool{"exists": exists}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
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
	rows, err := db.Query("SELECT nombre FROM ClubDeportivo")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var Clubes []Club
	for rows.Next() {
		var Club Club
		if err := rows.Scan(&Club.Nombre); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		Clubes = append(Clubes, Club)
	}
	// println(len(Clubes))
	templates.ExecuteTemplate(w, "create_user2.html", Clubes)
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
	var user Jugador
	err := row.Scan(&user.ID, &user.Nombres, &user.Apellido_paterno)
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
