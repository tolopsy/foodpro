package mongolayer

import (
	"context"
	"time"

	"github.com/tolopsy/foodpro/persistence"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type DBHandler struct {
	collection *mongo.Collection
	context    context.Context
}

func NewMongoDBHandler(dbURI, dbName string) (*DBHandler, error) {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbURI))
	if err != nil {
		return nil, err
	}
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		return nil, err
	}
	collection := client.Database(dbName).Collection("recipes")
	return &DBHandler{
		collection: collection,
		context:    ctx,
	}, nil
}

func (db *DBHandler) FetchAllRecipes() ([]persistence.Recipe, error) {
	cursor, err := db.collection.Find(db.context, bson.M{})
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
	result := db.collection.FindOne(db.context, documentArg)

	if err = result.Decode(&recipe); err != nil {
		return recipe, err
	}
	return recipe, nil
}

func (db *DBHandler) FindRecipesByTag(tag string) ([]persistence.Recipe, error) {
	searchArg := bson.M{"tags": tag}
	cursor, err := db.collection.Find(db.context, searchArg)
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
	_, err := db.collection.InsertOne(db.context, recipe)
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
	_, err = db.collection.UpdateByID(db.context, objectId, update)
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
	_, err = db.collection.DeleteOne(db.context, documentFilter)
	if err != nil {
		return err
	}
	return nil
}
