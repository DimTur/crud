package handler

import "crud/internal/domain"

type RecipeReq struct {
	ID          string       `json:"id"`
	AuthorID    string       `json:"user_id"`
	Name        string       `json:"name"`
	Ingredients []domain.Ing `json:"ingredients"`
	Temperature int          `json:"temperature"`
}
