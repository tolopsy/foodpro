package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/tolopsy/foodpro/handlers"
)

var err error
var ctx context.Context
var client *mongo.Client
var collection *mongo.Collection
var recipeHandler *handlers.RecipeHandler

func init() {
	ctx = context.Background()
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal("Error while pinging DB: " + err.Error())
	}
	log.Println("Connected to MongoDB")
	collection = client.Database(os.Getenv("MONGO_DATABASE")).Collection("recipes")
	recipeHandler = handlers.NewRecipeHandler(collection, ctx)
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
