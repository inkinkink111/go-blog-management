package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Blog struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	BlogID    string             `json:"blog_id" bson:"blog_id" example:"1234567890"`
	Title     string             `json:"title" bson:"title" example:"My Blog Title"`
	Content   string             `json:"content" bson:"content" example:"Blog content"`
	Slug      string             `json:"slug" bson:"slug" example:"my-blog-title"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at" example:"2021-01-01T00:00:00Z"`
	Tags      []string           `json:"tags" bson:"tags" example:"golang,redis"`
	AuthorID  string             `json:"author_id" bson:"author_id" example:"1234567890"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at" example:"2021-01-01T00:00:00Z"`
}
