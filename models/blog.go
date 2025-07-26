package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Blog struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	BlogID    string             `json:"blog_id" bson:"blog_id"`
	Title     string             `json:"title" bson:"title"`
	Content   string             `json:"content" bson:"content"`
	Slug      string             `json:"slug" bson:"slug"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	Tags      []string           `json:"tags" bson:"tags"`
	AuthorID  string             `json:"author_id" bson:"author_id"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}
