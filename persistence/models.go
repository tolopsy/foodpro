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
