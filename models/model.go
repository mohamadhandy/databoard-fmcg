package models

type AccessibleModel interface {
	Admin | []Admin | Brand | []Brand | Category | []Category | Product | []Product
}
