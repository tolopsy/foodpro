package persistence

type DatabaseHandler interface {
	FetchAllRecipes() ([]Recipe, error)
	GetRecipe(string) (Recipe, error)
	FindRecipesByTag(string) ([]Recipe, error)
	AddRecipe(*Recipe) error
	UpdateRecipe(string, Recipe) error
	DeleteRecipe(string) error
}

type CacheHandler interface {
	SetRecipes([]Recipe) error
	GetRecipes() ([]Recipe, error)
	ClearRecipes() error
}
