package main

import (
	// Importa las plantillas generadas
	"context"
	"ejemplo2/models"
	"ejemplo2/templates"
	"fmt"
	"log"
	"net/http"
)

// Datos de ejemplo
var regions = []models.Region{
	{ID: 1, Name: "Región 1"},
	{ID: 2, Name: "Región 2"},
}

var provincias = []models.Provincia{
	{ID: 1, Name: "Provincia 1-1", RegionID: 1},
	{ID: 2, Name: "Provincia 1-2", RegionID: 1},
	{ID: 3, Name: "Provincia 2-1", RegionID: 2},
}

var comunas = []models.Comuna{
	{ID: 1, Name: "Comuna 1-1-1", ProvinciaID: 1},
	{ID: 2, Name: "Comuna 1-1-2", ProvinciaID: 1},
	{ID: 3, Name: "Comuna 2-1-1", ProvinciaID: 2},
}

func regionHandler(w http.ResponseWriter, r *http.Request) {
	// Crear un componente de región
	regionComponent := templates.RegionComponent(regions)
	err := regionComponent.Render(context.Background(), w)
	if err != nil {
		log.Fatal(err)
	}
}

func provinciaHandler(w http.ResponseWriter, r *http.Request) {
	// Recuperar la región seleccionada y filtrar las provincias
	regionID := r.URL.Query().Get("region_id")
	var filteredProvincias []models.Provincia
	for _, provincia := range provincias {
		if fmt.Sprintf("%d", provincia.RegionID) == regionID {
			filteredProvincias = append(filteredProvincias, provincia)
		}
	}

	// Crear un componente de provincia
	provinciaComponent := templates.ProvinciaComponent(filteredProvincias)
	err := provinciaComponent.Render(context.Background(), w)
	if err != nil {
		log.Fatal(err)
	}
}

func comunaHandler(w http.ResponseWriter, r *http.Request) {
	// Recuperar la provincia seleccionada y filtrar las comunas
	provinciaID := r.URL.Query().Get("provincia_id")
	var filteredComunas []models.Comuna
	for _, comuna := range comunas {
		if fmt.Sprintf("%d", comuna.ProvinciaID) == provinciaID {
			filteredComunas = append(filteredComunas, comuna)
		}
	}

	// Crear un componente de comuna
	comunaComponent := templates.ComunaComponent(filteredComunas)
	err := comunaComponent.Render(context.Background(), w)
	if err != nil {
		log.Fatal(err)
	}
}

func baseHandler(w http.ResponseWriter, r *http.Request) {
	// Crear un componente de comuna

	base := templates.Base("hola", templates.RegionComponent(regions))
	err := base.Render(r.Context(), w)
	if err != nil {
		log.Fatal(err)
	}

}

func main() {

	// nada := nada(1, 2)
	// println(nada)

	// Generar las plantillas con templ generate antes de ejecutar el servidor.

	// Configurar las rutas
	http.HandleFunc("/", baseHandler)
	http.HandleFunc("/region", regionHandler)
	http.HandleFunc("/provincia", provinciaHandler)
	http.HandleFunc("/comuna", comunaHandler)

	// Iniciar el servidor
	fmt.Println("Servidor iniciado en http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
