package domain

type Ing struct {
	Amount int    `json:"amount"`
	Type   string `json:"type"`
}

type Recipe struct {
	ID          string `json:"id"`
	AuthorID    string `json:"user_id"`
	Name        string `json:"name"`
	Ingredients []Ing  `json:"ingredients"`
	Temperature int    `json:"temperature"`
}
