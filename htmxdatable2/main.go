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
	Fecha_nacimiento time.Time
	Comuna           string
	Direccion        string
	Historial        string
	Activo           int
}

type Club struct {
	ID            int
	Nombre        string
	Representante string
	Comuna        string
	Direccion     string
}
type Comuna struct {
	ID     int
	Nombre string
}
type Serie struct {
	ID     int
	Nombre string
}

// type PageData struct {
// 	Clubs        []Club
// 	Comunas      []Comuna
// 	Series       []Serie
// 	SelectClub   string
// 	SelectComuna string
// 	SelectSerie  string
// 	Jugadores    []Jugador
// 	//  UserName  string
// 	// IsAdmin   bool
// 	// Agrega más campos si es necesario
// }

type PageData struct {
	Clubs        []Club
	Comunas      []Comuna
	Series       []Serie
	SelectClub   string
	SelectComuna string
	SelectSerie  string
	Jugadores    Jugador
	Club         Club
	//  UserName  string
	// IsAdmin   bool
	// Agrega más campos si es necesario
}

// Función para formatear fechas
func formatFecha(t time.Time) string {
	// return t.Format("2006/01/02")
	return t.UTC().Format("2006-01-02")

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

	// Define el FuncMap con las funciones personalizadas
	funcMap := template.FuncMap{
		"formatFecha": formatFecha,
	}

	// Usa New y Funcs para aplicar el FuncMap antes de cargar las plantillas
	templates = template.Must(template.New("").Funcs(funcMap).ParseGlob("templates/*.html"))

	//templates = template.Must(template.ParseGlob("templates/*.html"))
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
	r.HandleFunc("DELETE /user/delete/{id}", DeleteUserHandler)
	r.HandleFunc("GET /user/find/{id}", FindUserHandler)
	// http.HandleFunc("/clubs/delete/", DeleteClubHandler)
	r.HandleFunc("GET /users/edit/{id}", EditUserHandler)
	r.HandleFunc("POST /users/edit/", EditUserHandlerPost)
	/************************************************************* club ****************/
	r.HandleFunc("GET /clubs/edit/{id}", EditClubHandler)
	r.HandleFunc("POST /club/edit/", EditClubHandlerPost)
	r.HandleFunc("GET /createClub", templateCreateClubHandler)
	r.HandleFunc("POST /club/create", CreateClubHandler)
	r.HandleFunc("DELETE /club/delete/{id}", DeleteClubHandler)
	r.HandleFunc("GET /listClubs", listClubHandler)

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
	rows, err := db.Query("SELECT id,rut,dv,nombres,apellido_paterno,apellido_materno,club_juega,foto,serie_juega,edad FROM jugador where activo  in (1) ")
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
	//http.Redirect(w, r, "/listClubs", http.StatusSeeOther)
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
	_, err := db.Exec("Update jugador set activo=0 WHERE id = ?", id)
	// _, err := db.Exec("DELETE FROM users WHERE id = ?", id)
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
	//////////////////////////////////////////////////////////////////////////////
	rowcomuna, err := db.Query("SELECT nombre FROM Comuna")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rowcomuna.Close()

	/*itera sobre los resultados de una consulta a la base de datos y extrae información de cada fila.
	for rows.Next():

	Esta es la estructura básica de un bucle for en Go, en el que se llama al método Next() de rows para iterar sobre los resultados de la consulta.
	rows.Next() devuelve true mientras haya más filas en el conjunto de resultados. Cuando se alcanza el final, devuelve false, y el bucle termina.
	*/
	var Comunas []Comuna
	for rowcomuna.Next() {
		var Comuna Comuna //crea una instancia de Comuna para almacenar los datos de cada fila.
		if err := rowcomuna.Scan(&Comuna.Nombre); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		Comunas = append(Comunas, Comuna)
	}
	//////////////////////////////////////////////////////////////////////////////
	rowseries, err := db.Query("SELECT nombre FROM Series")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rowseries.Close()

	var Series []Serie
	for rowseries.Next() {
		var Serie Serie //crea una instancia de Comuna para almacenar los datos de cada fila.
		if err := rowseries.Scan(&Serie.Nombre); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		Series = append(Series, Serie)
	}
	//////////////////////////////////////////////////////////////////////////////
	pageData := PageData{
		Clubs:   Clubes,
		Comunas: Comunas,
		Series:  Series,
		// UserName: "Juan Perez",
		// IsAdmin:  true,
	}
	// println(len(Clubes))

	templates.ExecuteTemplate(w, "create_user2.html", pageData)
}

func EditUserHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var clubNombre string
	var serienombre string
	var comunanombre string
	//var Users []Jugador
	row := db.QueryRow("SELECT id,rut,dv,nombres,apellido_paterno,apellido_materno,mail,edad,fecha_nacimiento,comuna,direccion,club_juega,serie_juega,historial,activo,foto FROM jugador WHERE id = ?", id)
	var user Jugador
	err := row.Scan(&user.ID, &user.Rut, &user.Dv, &user.Nombres, &user.Apellido_paterno, &user.Apellido_materno, &user.Mail, &user.Edad, &user.Fecha_nacimiento, &user.Comuna, &user.Direccion, &user.Club_juega, &user.Serie_juega, &user.Historial, &user.Activo, &user.Foto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//	Users = append(Users, user)
	// fmt.Printf("Fecha de Nacimiento: %s\n", user.Fecha_nacimiento)

	rows, err := db.Query("SELECT nombre FROM ClubDeportivo")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	clubNombre = user.Club_juega
	serienombre = user.Serie_juega
	comunanombre = user.Comuna
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
	//////////////////////////////////////////////////////////////////////////////
	rowcomuna, err := db.Query("SELECT nombre FROM Comuna")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rowcomuna.Close()

	var Comunas []Comuna
	for rowcomuna.Next() {
		var Comuna Comuna //crea una instancia de Comuna para almacenar los datos de cada fila.
		if err := rowcomuna.Scan(&Comuna.Nombre); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		Comunas = append(Comunas, Comuna)

	}
	//////////////////////////////////////////////////////////////////////////////
	rowseries, err := db.Query("SELECT nombre FROM Series")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rowseries.Close()

	var Series []Serie
	for rowseries.Next() {
		var Serie Serie //crea una instancia de Comuna para almacenar los datos de cada fila.
		if err := rowseries.Scan(&Serie.Nombre); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		Series = append(Series, Serie)

	}
	//////////////////////////////////////////////////////////////////////////////
	pageData := PageData{
		Clubs:        Clubes,
		Comunas:      Comunas,
		Series:       Series,
		SelectClub:   clubNombre,
		SelectComuna: comunanombre,
		SelectSerie:  serienombre,
		Jugadores:    user,
		// IsAdmin:  true,
	}
	// println(len(Clubes))

	templates.ExecuteTemplate(w, "edit_user.html", pageData)

	// templates.ExecuteTemplate(w, "edit_user.html", user)

	// }

}

func EditUserHandlerPost(w http.ResponseWriter, r *http.Request) {
	// var q string
	file, _, err := r.FormFile("foto")
	if err != nil {
		_, err := db.Exec("UPDATE jugador SET nombres = ?,apellido_paterno= ?,apellido_materno= ?,mail= ?,edad= ?,fecha_nacimiento= ?,comuna= ?,direccion= ?,club_juega= ?,serie_juega= ?,historial= ? WHERE id = ?", r.FormValue("nombres"), r.FormValue("apellido_paterno"), r.FormValue("apellido_materno"), r.FormValue("email"), r.FormValue("edad"), r.FormValue("fecha_nacimiento"), r.FormValue("comuna"), r.FormValue("direccion"), r.FormValue("club_juega"), r.FormValue("serie_juega"), r.FormValue("historial"), r.FormValue("id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
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

		_, err = db.Exec("UPDATE jugador SET nombres = ?,apellido_paterno= ?,apellido_materno= ?,mail= ?,edad= ?,fecha_nacimiento= ?,comuna= ?,direccion= ?,club_juega= ?,serie_juega= ?,historial= ?,foto=? WHERE id = ?", r.FormValue("nombres"), r.FormValue("apellido_paterno"), r.FormValue("apellido_materno"), r.FormValue("email"), r.FormValue("edad"), r.FormValue("fecha_nacimiento"), r.FormValue("comuna"), r.FormValue("direccion"), r.FormValue("club_juega"), r.FormValue("serie_juega"), r.FormValue("historial"), nombreFoto, r.FormValue("id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}

	println(file)
	fmt.Fprintf(w, "Modificacion  exitosamente ")
	//http.Redirect(w, r, "/listUsers", http.StatusSeeOther)

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

/*--------------------------------------------------------------------------------------------clubs-------------*/
func templateCreateClubHandler(w http.ResponseWriter, r *http.Request) {
	//////////////////////////////////////////////////////////////////////////////
	rowcomuna, err := db.Query("SELECT nombre FROM Comuna")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rowcomuna.Close()

	/*itera sobre los resultados de una consulta a la base de datos y extrae información de cada fila.
	for rows.Next():

	Esta es la estructura básica de un bucle for en Go, en el que se llama al método Next() de rows para iterar sobre los resultados de la consulta.
	rows.Next() devuelve true mientras haya más filas en el conjunto de resultados. Cuando se alcanza el final, devuelve false, y el bucle termina.
	*/
	var Comunas []Comuna
	for rowcomuna.Next() {
		var Comuna Comuna //crea una instancia de Comuna para almacenar los datos de cada fila.
		if err := rowcomuna.Scan(&Comuna.Nombre); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		Comunas = append(Comunas, Comuna)
	}

	pageData := PageData{
		Comunas: Comunas,
	}
	// println(len(Clubes))

	templates.ExecuteTemplate(w, "create_club.html", pageData)
}

// // Crear CLUB
func CreateClubHandler(w http.ResponseWriter, r *http.Request) {

	_, err := db.Exec("INSERT INTO ClubDeportivo (nombre,comuna,direccion,representante,activo) VALUES (?, ?, ?, ?, ?)", r.FormValue("nombre"), r.FormValue("comuna"), r.FormValue("direccion"), r.FormValue("representante"), 1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//fmt.Fprintf(w, "Club Ingeresado exitosamente ")
	//templates.ExecuteTemplate(w, "usuarioCreado.html", nil)
	http.Redirect(w, r, "/listClubs", http.StatusSeeOther)
	//listUsersHandler(w, r)
}

func listClubHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id,nombre,comuna, representante,direccion FROM ClubDeportivo where activo=1")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var Clubs []Club
	for rows.Next() {
		var club Club
		if err := rows.Scan(&club.ID, &club.Nombre, &club.Comuna, &club.Direccion, &club.Representante); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		Clubs = append(Clubs, club)
	}
	// println(len(users))
	templates.ExecuteTemplate(w, "club_table.html", Clubs)
}

func DeleteClubHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	_, err := db.Exec("Update ClubDeportivo set activo=0 WHERE id = ?", id)
	// _, err := db.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/listClubs", http.StatusSeeOther)
}

func EditClubHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	row := db.QueryRow("SELECT id,nombre,representante,comuna,direccion FROM ClubDeportivo WHERE id = ?", id)
	var club Club
	err := row.Scan(&club.ID, &club.Nombre, &club.Representante, &club.Comuna, &club.Direccion)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//////////////////////////////////////////////////////////////////////////////
	rowcomuna, err := db.Query("SELECT nombre FROM Comuna")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rowcomuna.Close()

	var Comunas []Comuna
	for rowcomuna.Next() {
		var Comuna Comuna //crea una instancia de Comuna para almacenar los datos de cada fila.
		if err := rowcomuna.Scan(&Comuna.Nombre); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		Comunas = append(Comunas, Comuna)

	}
	//////////////////////////////////////////////////////////////////////////////

	//////////////////////////////////////////////////////////////////////////////
	pageData := PageData{
		Club:    club,
		Comunas: Comunas,
	}

	templates.ExecuteTemplate(w, "edit_club.html", pageData)

	// templates.ExecuteTemplate(w, "edit_user.html", user)

	// }

}

func EditClubHandlerPost(w http.ResponseWriter, r *http.Request) {

	// // Parsear los datos del formulario
	// if err := r.ParseForm(); err != nil {
	// 	http.Error(w, "Error al procesar los datos del formulario", http.StatusBadRequest)
	// 	return
	// }

	// // Capturar los valores enviados
	// nombre := r.FormValue("nombre")
	// representante := r.FormValue("representante")
	// comuna := r.FormValue("comuna")
	// direccion := r.FormValue("direccion")
	// id := r.FormValue("id")

	// // Debugging: Imprimir los valores
	// log.Printf("Datos recibidos: nombre=%s, representante=%s, comuna=%s, direccion=%s, id=%s", nombre, representante, comuna, direccion, id)

	_, err := db.Exec("UPDATE ClubDeportivo SET nombre = ?,representante= ?,comuna= ?,direccion= ? WHERE id = ?", r.FormValue("nombre"), r.FormValue("representante"), r.FormValue("comuna"), r.FormValue("direccion"), r.FormValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/listClubs", http.StatusSeeOther)

}
