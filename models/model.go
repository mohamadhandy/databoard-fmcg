package models

type AccessibleModel interface {
	Admin | []Admin | Brand | []Brand
}
