package persistence

import (
	"time"
)

type Recipe struct {
	ID           interface{} `json:"id,omitempty" bson:"_id,omitempty"`
	Name         string      `json:"name,omitempty" bson:"name,omitempty"`
	Tags         []string    `json:"tags,omitempty" bson:"tags,omitempty"`
	Ingredients  []string    `json:"ingredients,omitempty" bson:"ingredients,omitempty"`
	Instructions []string    `json:"instructions,omitempty" bson:"instructions,omitempty"`
	PublishedAt  time.Time   `json:"publishedAt,omitempty" bson:"publishedAt,omitempty"`
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (user *User) VerifyUser() bool {
	// acceptedUsername & acceptedPassword are just stubs and not 
	// typical to a proper verification process
	acceptedUsername := "admin"
	acceptedPassword := "password"

	if user.Username != acceptedUsername || user.Password != acceptedPassword {
		return false
	}
	return true
}