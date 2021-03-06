package server

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/tolopsy/foodpro/api/persistence"
	"github.com/tolopsy/foodpro/api/persistence/cache"
)

type Handler struct {
	db    persistence.DatabaseHandler
	cache persistence.CacheHandler
}

func NewHandler(db persistence.DatabaseHandler, cache persistence.CacheHandler) *Handler {
	return &Handler{
		db:    db,
		cache: cache,
	}
}

func (handler *Handler) FetchAllRecipes(ctx *gin.Context) {
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

func (handler *Handler) FetchOneRecipe(ctx *gin.Context) {
	id := ctx.Param("id")
	recipe, err := handler.db.GetRecipe(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, recipe)
}

func (handler *Handler) SearchRecipesByTag(ctx *gin.Context) {
	tag := ctx.Query("tag")
	recipes, err := handler.db.FindRecipesByTag(tag)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	ctx.JSON(http.StatusOK, recipes)
}

func (handler *Handler) CreateNewRecipe(ctx *gin.Context) {
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

func (handler *Handler) UpdateRecipe(ctx *gin.Context) {
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

func (handler *Handler) DeleteRecipe(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := handler.db.DeleteRecipe(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	handler.cache.ClearRecipes()
	ctx.JSON(http.StatusNoContent, gin.H{"message": "Recipe has been deleted"})
}
