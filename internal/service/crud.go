package service

import (
	"crud/internal/domain"
	"crud/internal/repository/cache"
	"crud/internal/repository/recipedb"

	"github.com/google/uuid"
)

var recipes recipedb.DB

func Init(DB recipedb.DB) {
	recipes = DB
}

func Get(id string) (*domain.Recipe, error) {
	return recipes.Get(id)
}

func GetAll(page, limit int, sortBy string) (*cache.PaginatedResponse, error) {
	return recipes.GetAll(page, limit, sortBy)
}

func Delete(id string) error {
	return recipes.Delete(id)
}

func AddOrUpd(r *domain.Recipe) error {

	if r.ID == "" {
		r.ID = uuid.New().String()
	}

	return recipes.Set(r.ID, r)
}
