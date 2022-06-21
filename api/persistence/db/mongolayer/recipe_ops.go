package mongolayer

import (
	"time"

	"github.com/tolopsy/foodpro/api/persistence"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/bson"
)

func (db *DBHandler) FetchAllRecipes() ([]persistence.Recipe, error) {
	cursor, err := db.recipeCollection.Find(db.context, bson.M{})
	if err != nil {
		return nil, err
	}

	var recipes []persistence.Recipe
	if err = cursor.All(db.context, &recipes); err != nil {
		return nil, err
	}

	return recipes, nil
}

func (db *DBHandler) GetRecipe(id string) (persistence.Recipe, error) {
	var recipe persistence.Recipe
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return recipe, err
	}

	documentArg := bson.M{"_id": objectId}
	result := db.recipeCollection.FindOne(db.context, documentArg)

	if err = result.Decode(&recipe); err != nil {
		return recipe, err
	}
	return recipe, nil
}

func (db *DBHandler) FindRecipesByTag(tag string) ([]persistence.Recipe, error) {
	searchArg := bson.M{"tags": tag}
	cursor, err := db.recipeCollection.Find(db.context, searchArg)
	if err != nil {
		return nil, err
	}

	var recipes []persistence.Recipe
	if err = cursor.All(db.context, &recipes); err != nil {
		return nil, err
	}
	return recipes, nil
}

func (db *DBHandler) AddRecipe(recipe *persistence.Recipe) error {
	recipe.ID = primitive.NewObjectID()
	recipe.PublishedAt = time.Now()
	_, err := db.recipeCollection.InsertOne(db.context, recipe)
	if err != nil {
		return err
	}
	return nil
}

func (db *DBHandler) UpdateRecipe(id string, recipe persistence.Recipe) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{"$set": &recipe}
	_, err = db.recipeCollection.UpdateByID(db.context, objectId, update)
	if err != nil {
		return err
	}
	return nil
}

func (db *DBHandler) DeleteRecipe(id string) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	documentFilter := bson.M{"_id": objectId}
	_, err = db.recipeCollection.DeleteOne(db.context, documentFilter)
	if err != nil {
		return err
	}
	return nil
}
