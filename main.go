package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var err error
var ctx context.Context
var client *mongo.Client
var loadInitialData bool = false
var collection *mongo.Collection

type Recipe struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	Name         string             `json:"name" bson:"name"`
	Tags         []string           `json:"tags" bson:"tags"`
	Ingredients  []string           `json:"ingredients" bson:"ingredients"`
	Instructions []string           `json:"instructions" bson:"instructions"`
	PublishedAt  time.Time          `json:"publishedAt" bson:"publishedAt"`
}

var recipes []Recipe

func CreateNewRecipeHandler(c *gin.Context) {
	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error while creating new recipe: " + err.Error(),
		})
		return
	}
	recipe.ID = primitive.NewObjectID()
	recipe.PublishedAt = time.Now()
	_, err = collection.InsertOne(ctx, recipe)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while inserting new recipe: " + err.Error()})
		return
	}
	c.JSON(http.StatusCreated, recipe)
}

func FetchAllRecipesHandler(c *gin.Context) {
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(ctx)

	recipes := make([]Recipe, 0)
	for cursor.Next(ctx) {
		var recipe Recipe
		err = cursor.Decode(&recipe)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while decoding recipes from db" + err.Error()})
		}
		recipes = append(recipes, recipe)
	}

	c.JSON(http.StatusOK, recipes)
}

func UpdateRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error while parsing recipe to update: " + err.Error(),
		})
		return
	}

	index := -1
	for i := 0; i < len(recipes); i++ {
		if recipes[i].ID.String() == id {
			index = i
			break
		}
	}

	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Recipe not found",
		})
		return
	}

	recipes[index] = recipe
	c.JSON(http.StatusOK, recipes[index])
}

func DeleteRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	index := -1

	for i := 0; i < len(recipes); i++ {
		if recipes[i].ID.String() == id {
			index = i
			break
		}
	}

	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Recipe not found",
		})
		return
	}

	recipes = append(recipes[:index], recipes[index+1:]...)
	c.JSON(http.StatusOK, gin.H{
		"message": "Recipe has been deleted",
	})
}

func SearchRecipesHandler(c *gin.Context) {
	tag := c.Query("tag")
	recipeListSearch := make([]Recipe, 0)

	for i := 0; i < len(recipes); i++ {
		found := false
		for _, recipeTag := range recipes[i].Tags {
			if strings.EqualFold(recipeTag, tag) {
				found = true
			}
		}

		if found {
			recipeListSearch = append(recipeListSearch, recipes[i])
		}
	}

	c.JSON(http.StatusOK, recipeListSearch)
}

func loadRecipesIntoDb() {
	recipes = make([]Recipe, 0)
	file, _ := ioutil.ReadFile("recipes.json")
	_ = json.Unmarshal(file, &recipes)

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

func init() {
	ctx = context.Background()
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB")
	collection = client.Database(os.Getenv("MONGO_DATABASE")).Collection("recipes")

	// loads preliminary data into db when we need it.
	// handy for test situations or when using a new empty db in local.
	// TODO: use environment variable.
	if loadInitialData {
		loadRecipesIntoDb()
	}
}

func main() {
	engine := gin.Default()
	engine.POST("/recipes", CreateNewRecipeHandler)
	engine.GET("/recipes", FetchAllRecipesHandler)
	engine.GET("/recipes/search", SearchRecipesHandler)
	engine.PUT("/recipes/:id", UpdateRecipeHandler)
	engine.DELETE("/recipes/:id", DeleteRecipeHandler)
	engine.Run()
}
