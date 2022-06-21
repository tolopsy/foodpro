package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/tolopsy/foodpro/api/provider"
	"github.com/tolopsy/foodpro/api/server"
	auth "github.com/tolopsy/foodpro/api/server/middleware/authentication"
	session_auth "github.com/tolopsy/foodpro/api/server/middleware/authentication/session"
)

var handler *server.Handler
var authMiddleware auth.AuthMiddleware

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error while loading dev environment variables: " + err.Error())
	}

	dbType, dbURI, dbName := os.Getenv("DB_TYPE"), os.Getenv("DB_URI"), os.Getenv("DB_NAME")
	cacheType, cacheHost, cachePassword := os.Getenv("CACHE_TYPE"), os.Getenv("CACHE_HOST"), os.Getenv("CACHE_PASSWORD")

	SESS_STORE_ADDRESS := os.Getenv("SESS_STORE_ADDRESS")
	SESS_STORE_PASSWORD := os.Getenv("SESS_STORE_PASSWORD")
	SESS_STORE_KEY := os.Getenv("SESS_STORE_KEY")

	db, err := provider.NewDBHandler(dbType, dbURI, dbName)
	if err != nil {
		log.Fatal("Error while obtaining db handler -> " + err.Error())
	}

	cache, err := provider.NewCacheHandler(cacheType, cacheHost, cachePassword)
	if err != nil {
		log.Fatal("Error while obtainiing cache server -> " + err.Error())
	}

	if err != nil {
		log.Fatal("Error while initializing session store -> " + err.Error())
	}

	handler = server.NewHandler(db, cache)
	authMiddleware, err = session_auth.NewSessionAuth(
		SESS_STORE_KEY,
		SESS_STORE_ADDRESS,
		SESS_STORE_PASSWORD,
		db.VerifyUser,
	)
	if err != nil {
		log.Fatal("Error while initializing authentication middleware -> " + err.Error())
	}
}

func main() {
	engine := gin.Default()

	auth.LoadSpecialFeatures(authMiddleware, engine)

	engine.GET("/recipes", handler.FetchAllRecipes)
	engine.GET("recipes/:id", handler.FetchOneRecipe)
	engine.GET("/recipes/search", handler.SearchRecipesByTag)
	engine.POST("/sign-in", authMiddleware.SignIn)
	engine.GET("/sign-out", authMiddleware.SignOut)

	authorized := engine.Group("/")
	authorized.Use(authMiddleware.Authenticate())
	authorized.POST("/recipes", handler.CreateNewRecipe)
	authorized.PATCH("/recipes/:id", handler.UpdateRecipe)
	authorized.DELETE("/recipes/:id", handler.DeleteRecipe)

	engine.Run(":8080")
}



/*
func loadUsersIntoDb() {
	users := map[string]string{
		"admin":   "password",
		"guest":   "guest",
		"tolopsy": "foodpro",
	}

	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("DB_URI")))
	if err != nil {
		log.Fatal("kILELEYI", err.Error())
	}
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal("kITUNLELEYI", err.Error())
	}

	userCollection := client.Database(os.Getenv("DB_NAME")).Collection("users")

	// basic hashing & salting
	h := sha256.New()
	for username, password := range users {
		userCollection.InsertOne(ctx, bson.M{
			"username": username,
			"password": string(h.Sum([]byte(password))),
		})
	}
}


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
