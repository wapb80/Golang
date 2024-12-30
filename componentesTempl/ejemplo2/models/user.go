package models

type Region struct {
	ID   int
	Name string
}

type Provincia struct {
	ID       int
	Name     string
	RegionID int
}

type Comuna struct {
	ID          int
	Name        string
	ProvinciaID int
}

type User struct {
	Name  string
	Email string
	Age   int
}
