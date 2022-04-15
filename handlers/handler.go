package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/tolopsy/foodpro/models"
)

type RecipeHandler struct {
	collection *mongo.Collection
	context    context.Context
}

func NewRecipeHandler(collection *mongo.Collection, ctx context.Context) *RecipeHandler {
	return &RecipeHandler{
		collection: collection,
		context:    ctx,
	}
}

func (handler *RecipeHandler) FetchAllRecipes(ctx *gin.Context) {
	cursor, err := handler.collection.Find(handler.context, bson.M{})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error while fetching all recipes from db -> " + err.Error()})
		return
	}

	var recipes []models.Recipe
	if err = cursor.All(handler.context, &recipes); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error while parsing recipes -> " + err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, recipes)
}

func (handler *RecipeHandler) FetchOneRecipe(ctx *gin.Context) {
	id := ctx.Param("id")
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID -> "+ err.Error()})
		return
	}
	documentArg := bson.M{"_id": objectId}
	result := handler.collection.FindOne(handler.context, documentArg)
	var recipe models.Recipe
	
	if err = result.Decode(&recipe); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error while parsing recipe from db -> " + err.Error()})
		return
	}
	
	ctx.JSON(http.StatusOK, recipe)
}

func (handler *RecipeHandler) SearchRecipesByTag(ctx *gin.Context) {
	tag := ctx.Query("tag")
	searchArg := bson.M{"tags": tag}
	cursor, err := handler.collection.Find(handler.context, searchArg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error while searching DB -> " + err.Error()})
		return
	}
	var recipes []models.Recipe
	
	if err = cursor.All(handler.context, &recipes); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error while parsing resultss from DB -> " + err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, recipes)
}

func (handler *RecipeHandler) CreateNewRecipe(ctx *gin.Context) {
	var recipe models.Recipe
	if err := ctx.ShouldBindJSON(&recipe); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error while parsing request data -> " + err.Error()})
		return
	}
	
	recipe.ID = primitive.NewObjectID()
	recipe.PublishedAt = time.Now()
	_, err := handler.collection.InsertOne(handler.context, recipe)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error while inserting new recipe -> " + err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, recipe)
}

func (handler *RecipeHandler) UpdateRecipe(ctx *gin.Context) {
	id := ctx.Param("id")
	var recipe models.Recipe
	if err := ctx.ShouldBindJSON(&recipe); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error while parsing request data -> " + err.Error()})
		return
	}

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid error -> " + err.Error()})
		return
	}

	documentUpdate := bson.D{
		{"name", recipe.Name},
		{"instructions", recipe.Instructions},
		{"ingredients", recipe.Ingredients},
		{"tags", recipe.Tags},
	}

	_, err = handler.collection.UpdateByID(handler.context, objectId, bson.D{{"$set", documentUpdate}})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error while updating recipe -> " + err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Recipe has been updated"})
}

func (handler *RecipeHandler) DeleteRecipe(ctx *gin.Context) {
	id := ctx.Param("id")
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID -> " + err.Error()})
		return
	}
	
	documentArg := bson.M{"_id": objectId}
	_, err = handler.collection.DeleteOne(handler.context, documentArg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error while trying to delete recipe -> " + err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Recipe has been deleted"})
}

