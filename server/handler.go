package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/tolopsy/foodpro/persistence"
)

type RecipeHandler struct {
	db persistence.DatabaseHandler
}

func NewRecipeHandler(dbHandler persistence.DatabaseHandler) *RecipeHandler {
	return &RecipeHandler{
		db: dbHandler,
	}
}

func (handler *RecipeHandler) FetchAllRecipes(ctx *gin.Context) {
	recipes, err := handler.db.FetchAllRecipes()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
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

	ctx.JSON(http.StatusOK, gin.H{"message": "Recipe has been updated"})
}

func (handler *RecipeHandler) DeleteRecipe(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := handler.db.DeleteRecipe(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusNoContent, gin.H{"message": "Recipe has been deleted"})
}
