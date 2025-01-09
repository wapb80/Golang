package models

type Comuna struct {
	ID   string
	Name string
}

type Region struct {
	ID   string
	Name string
}

type Provincia struct {
	ID   string
	Name string
}

// Datos simulados
var Data = map[string]map[string][]Comuna{
	"Provincia1": {
		"Region1": {{ID: "1-1", Name: "Comuna1-1"}, {ID: "1-2", Name: "Comuna1-2"}},
		"Region2": {{ID: "2-1", Name: "Comuna2-1"}, {ID: "2-2", Name: "Comuna2-2"}},
	},
	"Provincia2": {
		"Region3": {{ID: "3-1", Name: "Comuna3-1"}, {ID: "3-2", Name: "Comuna3-2"}},
		"Region4": {{ID: "4-1", Name: "Comuna4-1"}, {ID: "4-2", Name: "Comuna4-2"}},
	},
}
