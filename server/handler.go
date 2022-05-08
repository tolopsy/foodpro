package server

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/tolopsy/foodpro/persistence"
	"github.com/tolopsy/foodpro/persistence/cache"
)

type RecipeHandler struct {
	db    persistence.DatabaseHandler
	cache persistence.CacheHandler
}

func NewRecipeHandler(db persistence.DatabaseHandler, cache persistence.CacheHandler) *RecipeHandler {
	return &RecipeHandler{
		db:    db,
		cache: cache,
	}
}

func (handler *RecipeHandler) FetchAllRecipes(ctx *gin.Context) {
	fetchFromDB := false
	recipes, err := handler.cache.GetRecipes()
	if err == cache.ErrorKeyDoesNotExist {
		fetchFromDB = true
	} else if err != nil {
		log.Println("Error while fetching recipes from cache -> " + err.Error())
		fetchFromDB = true
	}

	if fetchFromDB {
		recipes, err = handler.db.FetchAllRecipes()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		handler.cache.SetRecipes(recipes)
	}

	ctx.JSON(http.StatusOK, recipes)
}

func (handler *RecipeHandler) FetchOneRecipe(ctx *gin.Context) {
	id := ctx.Param("id")
	recipe, err := handler.db.GetRecipe(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, recipe)
}

func (handler *RecipeHandler) SearchRecipesByTag(ctx *gin.Context) {
	tag := ctx.Query("tag")
	recipes, err := handler.db.FindRecipesByTag(tag)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	ctx.JSON(http.StatusOK, recipes)
}

func (handler *RecipeHandler) CreateNewRecipe(ctx *gin.Context) {
	var recipe persistence.Recipe
	if err := ctx.ShouldBindJSON(&recipe); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := handler.db.AddRecipe(&recipe); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	handler.cache.ClearRecipes()
	ctx.JSON(http.StatusOK, recipe)
}

func (handler *RecipeHandler) UpdateRecipe(ctx *gin.Context) {
	id := ctx.Param("id")
	var recipe persistence.Recipe
	if err := ctx.ShouldBindJSON(&recipe); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error while parsing request data -> " + err.Error()})
		return
	}

	if err := handler.db.UpdateRecipe(id, recipe); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	handler.cache.ClearRecipes()
	ctx.JSON(http.StatusOK, gin.H{"message": "Recipe has been updated"})
}

func (handler *RecipeHandler) DeleteRecipe(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := handler.db.DeleteRecipe(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	handler.cache.ClearRecipes()
	ctx.JSON(http.StatusNoContent, gin.H{"message": "Recipe has been deleted"})
}
