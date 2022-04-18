package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/tolopsy/foodpro/provider"
	"github.com/tolopsy/foodpro/server"
)

var recipeHandler *server.RecipeHandler

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error while loading dev environment variables: " + err.Error())
	}

	dbType, dbURI, dbName := os.Getenv("DB_TYPE"), os.Getenv("DB_URI"), os.Getenv("DB_NAME")
	db, err := provider.NewDBHandler(dbType, dbURI, dbName)
	if err != nil {
		log.Fatal("Error while obtaining db handler: " + err.Error())
	}
	recipeHandler = server.NewRecipeHandler(db)
}

func main() {
	engine := gin.Default()
	engine.GET("/recipes", recipeHandler.FetchAllRecipes)
	engine.GET("recipes/:id", recipeHandler.FetchOneRecipe)
	engine.GET("/recipes/search", recipeHandler.SearchRecipesByTag)
	engine.POST("/recipes", recipeHandler.CreateNewRecipe)
	engine.PATCH("/recipes/:id", recipeHandler.UpdateRecipe)
	engine.DELETE("/recipes/:id", recipeHandler.DeleteRecipe)
	engine.Run()
}

/*
func loadRecipesIntoDb() {
	recipes = make([]Recipe, 0)
	file, _ := ioutil.ReadFile("recipes.json")
	err = json.Unmarshal(file, &recipes)
	if err != nil {
		log.Fatal("Error while loading recipes data: " + err.Error())
	}

	var listOfRecipes []interface{}
	for _, recipe := range recipes {
		listOfRecipes = append(listOfRecipes, recipe)
	}

	insertManyResult, err := collection.InsertMany(ctx, listOfRecipes)
	if err != nil {
		log.Fatal("Error while inserting initial recipes data " + err.Error())
	}
	log.Println("Inserted recipes: ", len(insertManyResult.InsertedIDs))
}
*/
