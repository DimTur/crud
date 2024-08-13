package recipedb

import (
	"crud/internal/domain"
	"crud/internal/repository/cache"
)

type DB interface {
	Get(id string) (*domain.Recipe, error)
	GetAll(page, limit int, sortBy string) (*cache.PaginatedResponse, error)
	Set(id string, recipe *domain.Recipe) error
	Delete(id string) error
}
