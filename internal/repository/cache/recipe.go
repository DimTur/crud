package cache

import (
	"context"
	"crud/internal/domain"
	"errors"
	"sort"
	"sync"
)

type RecipeCache struct {
	pool map[string]*domain.Recipe
	mtx  sync.RWMutex
}

const RecipeDumpFileName = "recipes.json"

func RecipeCacheInit(ctx context.Context, wg *sync.WaitGroup) (*RecipeCache, error) {
	var c RecipeCache
	c.pool = make(map[string]*domain.Recipe)

	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		makeDump(RecipeDumpFileName, c.pool)
	}()

	if err := loadFromDump(RecipeDumpFileName, &c.pool); err != nil {
		return nil, err
	}

	return &c, nil
}

func (c *RecipeCache) Get(id string) (*domain.Recipe, error) {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	if val, ok := c.pool[id]; ok {
		return val, nil
	}
	return nil, errors.New("recipe not found")
}

func (c *RecipeCache) Set(id string, recipe *domain.Recipe) error {

	c.mtx.Lock()
	c.pool[id] = recipe
	c.mtx.Unlock()

	return nil
}
func (c *RecipeCache) Delete(id string) error {

	c.mtx.Lock()
	delete(c.pool, id)
	c.mtx.Unlock()

	return nil
}

type RecipeEntry struct {
	ID     string
	Recipe *domain.Recipe
}

// Becouse we using cache I use page and limit: all data in memory.
// If we use SQL we need to use limit and offset.
// Pagination is best done from the database level.
func (c *RecipeCache) GetAll(page, limit int, sortBy string) ([]*domain.Recipe, error) {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	if len(c.pool) == 0 {
		return nil, errors.New("no recipes")
	}

	recipes := make([]*domain.Recipe, 0, len(c.pool))
	for _, recipe := range c.pool {
		recipes = append(recipes, recipe)
	}

	if sortBy == "name" {
		sort.Slice(recipes, func(i, j int) bool {
			return recipes[i].Name < recipes[j].Name
		})
	}

	startIdx := (page - 1) * limit
	endIdx := startIdx + limit

	if startIdx >= len(recipes) {
		return nil, errors.New("page out of range")
	}
	if endIdx > len(recipes) {
		endIdx = len(recipes)
	}

	return recipes[startIdx:endIdx], nil
}
