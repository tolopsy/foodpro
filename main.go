package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while fetching all recipes" + err.Error()})
		return
	}
	defer cursor.Close(ctx)

	recipes := make([]Recipe, 0)
	err = cursor.All(ctx, &recipes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while decoding recipes from db: " + err.Error()})
		return
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

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid ID"})
		return
	}
	filterById := bson.M{"_id": objectId}
	updateToMake := bson.D{{"$set", bson.D{
		{"name", recipe.Name},
		{"instructions", recipe.Instructions},
		{"ingredients", recipe.Ingredients},
		{"tags", recipe.Tags},
	}}}

	_, err = collection.UpdateOne(ctx, filterById, updateToMake)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while updating a recipe: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Recipe has been updated"})
}

func DeleteRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	_, err = collection.DeleteOne(ctx, bson.M{"_id": objectId})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occured during deletion " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Recipe has been deleted"})
}

func FetchOneRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	searchArg := bson.M{"_id": objectId}
	result := collection.FindOne(ctx, searchArg)

	var recipe Recipe
	err = result.Decode(&recipe)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while parsing fetched recipe" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, recipe)
}

func SearchRecipesByTagHandler(c *gin.Context) {
	tag := c.Query("tag")
	cursor, err := collection.Find(ctx, bson.M{"tags": tag})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while searching for recipe to delete " + err.Error()})
		return
	}
	defer cursor.Close(ctx)

	var listOfRecipes []Recipe
	err = cursor.All(ctx, &listOfRecipes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while parsing query results: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, listOfRecipes)
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

func init() {
	ctx = context.Background()
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal("Error while pinging DB: " + err.Error())
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
	engine.GET("recipes/:id", FetchOneRecipeHandler)
	engine.GET("/recipes/search", SearchRecipesByTagHandler)
	engine.PUT("/recipes/:id", UpdateRecipeHandler)
	engine.DELETE("/recipes/:id", DeleteRecipeHandler)
	engine.Run()
}
