package mongolayer

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type DBHandler struct {
	recipeCollection *mongo.Collection
	userCollection   *mongo.Collection
	context          context.Context
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

	recipeCollection := client.Database(dbName).Collection("recipes")
	userCollection := client.Database(dbName).Collection("users")

	return &DBHandler{
		recipeCollection: recipeCollection,
		userCollection:   userCollection,
		context:          ctx,
	}, nil
}
