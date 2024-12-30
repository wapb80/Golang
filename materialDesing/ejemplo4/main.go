package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"reflect"
	"strconv"

	_ "github.com/lib/pq"
)

// Handler para generar el gráfico
func chartHandler(w http.ResponseWriter, r *http.Request) {
	archivo := r.URL.Query().Get("archivo")
	// Crear un gráfico de barras
	// bar := charts.NewBar()
	// bar.SetGlobalOptions(charts.WithTitleOpts(opts.Title{Title: "Gráfico de Barras"}))
	// bar.SetXAxis([]string{"Lunes", "Martes", "Miércoles", "Jueves", "Viernes", "Sábado", "Domingo"}).
	// 	AddSeries("Categoría A", generateBarItems())

	// // Renderizar el gráfico como HTML
	// w.Header().Set("Content-Type", "text/html")
	// bar.Render(w)

	fmt.Println(reflect.TypeOf(archivo))
	fmt.Print(archivo)

	// Si el valor de archivo es una cadena JSON, intentar decodificarlo
	var archivoJSON map[string]interface{}
	err := json.Unmarshal([]byte(archivo), &archivoJSON)
	if err != nil {
		http.Error(w, "Error al decodificar archivo", http.StatusBadRequest)
		return
	}

	// Imprimir el objeto decodificado
	// fmt.Println("Objeto decodificado:", archivoJSON)

	// Responder con el objeto decodificado
	// fmt.Fprintf(w, "Objeto decodificado: %+v", archivoJSON)

	renderTemplate(w, "graficosPrueba.html", nil)
}

var (
	tmpl = template.Must(template.ParseGlob("templates/*.html"))
	db   *sql.DB
)

// Configurar la conexión a la base de datos
func initDB() {
	var err error
	connStr := "postgres://postgres:Sead_2023%23@192.168.8.2:5432/encvulne"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error al conectar a la base de datos: %v", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatalf("No se puede conectar a la base de datos: %v", err)
	}
	log.Println("Conexión a la base de datos exitosa")
}

func renderTemplate(w http.ResponseWriter, tmplName string, data interface{}) {
	err := tmpl.ExecuteTemplate(w, tmplName, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func menuReportesHandler(w http.ResponseWriter, r *http.Request) {
	// Consultar categorías y regiones desde la base de datos
	selectorPanel := queryList("SELECT id, opcion FROM opciones_menu where menu= 1 and categoria='selectorPanel' ")
	anio := queryList("select 0 as id, 'Todos' as opcion union all SELECT id, opcion FROM opciones_menu where categoria='anio' ")
	encuesta := queryList("select 0 as id, 'Todas' as opcion union all SELECT id, opcion FROM opciones_menu where categoria='encuesta' ")
	regiones := queryList("select 0 as codregion, 'Todas' as nombreregion union all SELECT codregion, nombreregion FROM dimregion")

	renderTemplate(w, "menu_reportes.html", map[string]interface{}{
		"SelectorPanel": selectorPanel,
		"Anio":          anio,
		"Encuesta":      encuesta,
		"Regiones":      regiones,
	})

}

func provinciasHandler(w http.ResponseWriter, r *http.Request) {
	regionIDStr := r.URL.Query().Get("region")
	//println(regionIDStr)
	// Intentar convertir el valor a un entero
	if regionIDStr != "0" {
		regionID, err := strconv.Atoi(regionIDStr)
		if err != nil {
			log.Printf("Error al convertir region_id: %v", err)
			http.Error(w, "Parámetro region_id inválido", http.StatusBadRequest)
			return
		}
		query := "select 0 as codprovincia, 'Todas' as nombreprovincia union all SELECT codprovincia,nombreprovincia FROM dimprovincia WHERE codregion = $1"
		provincias := queryList(query, regionID)

		renderTemplate(w, "select_option_region.html", map[string]interface{}{
			"Options": provincias,
		})

	}
}

func filtrosHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "filtros.html", nil)
	// err := tmpl.Execute(w, "filtros.html")
	// if err != nil {
	// 	log.Printf("Error al renderizar el fragmento: %v", err) // Imprimir el error
	// 	http.Error(w, fmt.Sprintf("Error al renderizar el fragmento: %v", err), http.StatusInternalServerError)
	// }

}

func graficosHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	var filters map[string]string
	// Parsear el JSON recibido
	err := json.NewDecoder(r.Body).Decode(&filters)
	if err != nil {
		http.Error(w, "Error al procesar los datos", http.StatusBadRequest)
		return
	}
	// aca verificar si existe el html o se debe crear

	// Aquí puedes procesar los filtros recibidos
	fmt.Println("Filtros recibidos:", filters)
	// renderTemplate(w, "grafico.html", nil)
}

func comunasHandler(w http.ResponseWriter, r *http.Request) {
	provinciaIDStr := r.URL.Query().Get("provincia")
	// Intentar convertir el valor a un entero
	provinciaID, err := strconv.Atoi(provinciaIDStr)
	if err != nil {
		log.Printf("Error al convertir region_id: %v", err)
		http.Error(w, "Parámetro region_id inválido", http.StatusBadRequest)
		return
	}
	query := "select 0 as codcomuna, 'Todas' as nombrecomuna union all SELECT codcomuna, nombrecomuna FROM dimcomuna WHERE codprovincia = $1"
	comunas := queryList(query, provinciaID)

	renderTemplate(w, "select_options.html", map[string]interface{}{
		"Options": comunas,
	})
}

// Helper para consultar datos y formatearlos
func queryList(query string, args ...interface{}) []map[string]string {
	rows, err := db.Query(query, args...)
	if err != nil {
		log.Printf("Error en la consulta: %v", err)
		return nil
	}
	defer rows.Close()

	var result []map[string]string
	for rows.Next() {
		var id, nombre string
		if err := rows.Scan(&id, &nombre); err != nil {
			log.Printf("Error al escanear: %v", err)
			continue
		}
		result = append(result, map[string]string{"ID": id, "Nombre": nombre})
	}
	return result
}

func main() {
	// Servir el directorio estático
	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	initDB()
	defer db.Close()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "base.html", map[string]interface{}{
			"Title": "Bienvenida",
		})
	})

	http.HandleFunc("/menu/reportes", menuReportesHandler)
	http.HandleFunc("/botonFiltros", filtrosHandler)
	http.HandleFunc("/provincias", provinciasHandler)
	http.HandleFunc("/comunas", comunasHandler)
	http.HandleFunc("/graficos", graficosHandler)
	http.HandleFunc("/chart", chartHandler)

	log.Println("Servidor iniciado en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
