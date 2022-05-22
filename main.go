package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/tolopsy/foodpro/provider"
	"github.com/tolopsy/foodpro/server"
	auth "github.com/tolopsy/foodpro/server/middleware/authentication"
	jwt_auth "github.com/tolopsy/foodpro/server/middleware/authentication/jwt"
)

var handler *server.Handler
var authMiddleware auth.AuthMiddleware

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error while loading dev environment variables: " + err.Error())
	}

	dbType, dbURI, dbName := os.Getenv("DB_TYPE"), os.Getenv("DB_URI"), os.Getenv("DB_NAME")
	cacheType, cacheHost, cachePassword := os.Getenv("CACHE_TYPE"), os.Getenv("CACHE_HOST"), os.Getenv("CACHE_PASSWORD")

	db, err := provider.NewDBHandler(dbType, dbURI, dbName)
	if err != nil {
		log.Fatal("Error while obtaining db handler -> " + err.Error())
	}

	cache, err := provider.NewCacheHandler(cacheType, cacheHost, cachePassword)
	if err != nil {
		log.Fatal("Error while obtainiing cache server -> " + err.Error())
	}
	handler = server.NewHandler(db, cache)
	authMiddleware = jwt_auth.NewJWTAuth(os.Getenv("JWT_SECRET"), db.VerifyUser)
}

func main() {
	engine := gin.Default()
	engine.GET("/recipes", handler.FetchAllRecipes)
	engine.GET("recipes/:id", handler.FetchOneRecipe)
	engine.GET("/recipes/search", handler.SearchRecipesByTag)
	engine.POST("/sign-in", authMiddleware.SignIn)

	auth.LoadSpecialHandlers(authMiddleware, engine)

	authorized := engine.Group("/")
	authorized.Use(authMiddleware.Authenticate())
	authorized.POST("/recipes", handler.CreateNewRecipe)
	authorized.PATCH("/recipes/:id", handler.UpdateRecipe)
	authorized.DELETE("/recipes/:id", handler.DeleteRecipe)

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
