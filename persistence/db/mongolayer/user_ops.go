package mongolayer

import (
	"crypto/sha256"

	"github.com/tolopsy/foodpro/persistence"
	"go.mongodb.org/mongo-driver/bson"
)

// Basic verification as the hashing and salting herein is
// only to demo implementation
func (db *DBHandler) VerifyUser(user persistence.User) bool {
	h := sha256.New()

	userCredentials := bson.M{"username": user.Username, "password": string(h.Sum([]byte(user.Password)))}
	result := db.userCollection.FindOne(db.context, userCredentials)
	
	return result.Err() == nil
}
