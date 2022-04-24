package provider

import (
	"github.com/tolopsy/foodpro/persistence"
	"github.com/tolopsy/foodpro/persistence/db/mongolayer"
)

type DBTYPE string

const (
	MONGO_DB DBTYPE = "mongodb"
)

func NewDBHandler(dbType, dbURI, dbName string) (persistence.DatabaseHandler, error) {
	switch DBTYPE(dbType) {
	case MONGO_DB:
		return mongolayer.NewMongoDBHandler(dbURI, dbName)
	}

	return nil, nil
}
