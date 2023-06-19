package repositories

import "klikdaily-databoard/models"

type RepositoryResult[T any | models.AccessibleModel] struct {
	Data       T
	Error      error
	StatusCode int
	Message    string
}
