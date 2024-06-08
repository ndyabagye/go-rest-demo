package recipes

// Recipe represents a recipe
type Recipe struct {
	Name        string       `json:"name"`
	Ingredients []Ingredient `json:"ingredients"`
}

// Ingredient represents the individual ingredients
type Ingredient struct {
	Name string `json:"name"`
}
