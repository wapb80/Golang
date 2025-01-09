package main

import (
	"context"
	"ejemplo3/models"
	"ejemplo3/templates"
	"net/http"
)

// Handler para la página principal
func indexHandler(w http.ResponseWriter, r *http.Request) {
	provincias := []models.Provincia{
		{ID: "Provincia1", Name: "Provincia 1"},
		{ID: "Provincia2", Name: "Provincia 2"},
	}

	var comunas = []models.Comuna{
		{ID: "1", Name: "Comuna 1-1-1"},
		{ID: "2", Name: "Comuna 1-1-2"},
		{ID: "3", Name: "Comuna 2-1-1"},
	}

	templates.IndexPage(provincias, comunas).Render(context.Background(), w)
	// regionComponent := templates.RegionComponent(regions)
	// err := regionComponent.Render(context.Background(), w)
	// if err != nil {
	// 	log.Fatal(err)
	// }
}

// Handler para cargar regiones
func regionsHandler(w http.ResponseWriter, r *http.Request) {
	provincia := r.URL.Query().Get("provincia")
	regions := []models.Region{}

	if reg, ok := models.Data[provincia]; ok {
		for region := range reg {
			regions = append(regions, models.Region{ID: region, Name: region})
		}
	}

	templates.SelectComponent("region", "region", "Seleccionar Región", regions).Render(r.Context(), w)
}

// Handler para cargar comunas
func communesHandler(w http.ResponseWriter, r *http.Request) {
	provincia := r.URL.Query().Get("provincia")
	region := r.URL.Query().Get("region")
	comunas := []models.Comuna{}

	if reg, ok := models.Data[provincia]; ok {
		if com, ok := reg[region]; ok {
			comunas = append(comunas, com...)
		}
	}

	templates.ComunaComponent(comunas).Render(r.Context(), w)
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/regions", regionsHandler)
	http.HandleFunc("/communes", communesHandler)

	// Servir estáticos
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.ListenAndServe(":8080", nil)
}
