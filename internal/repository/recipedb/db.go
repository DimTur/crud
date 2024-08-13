package recipedb

import (
	"crud/internal/domain"
)

type DB interface {
	Get(id string) (*domain.Recipe, error)
	GetAll(page, limit int, sortBy string) ([]*domain.Recipe, error)
	Set(id string, recipe *domain.Recipe) error
	Delete(id string) error
}
