package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/tolopsy/foodpro/api/provider"
	"github.com/tolopsy/foodpro/api/server"
	auth "github.com/tolopsy/foodpro/api/server/middleware/authentication"
	cors_middleware "github.com/tolopsy/foodpro/api/server/middleware/cors"
	session_auth "github.com/tolopsy/foodpro/api/server/middleware/authentication/session"
)

var handler *server.Handler
var authMiddleware auth.AuthMiddleware
var corsRule cors.Config

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

	// cors rule
	corsRule = cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "OPTIONS"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
}

func main() {
	engine := gin.Default()
	engine.Use(cors_middleware.NewCorsMiddleware(corsRule))
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

