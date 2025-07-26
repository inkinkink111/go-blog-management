package repositories

import (
	"context"
	"inkinkink111/go-blog-management/db"
	"inkinkink111/go-blog-management/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BlogRepository struct {
	collection *mongo.Collection
}

func NewBlogRepository() *BlogRepository {
	return &BlogRepository{
		collection: db.DB.Collection("blogs"), // your collection name
	}
}

func (br *BlogRepository) GetAllBlogs(page, limit int, tags []string) ([]models.Blog, int64, error) {
	var blogs []models.Blog
	// Get total blogs count
	totalCount, err := br.collection.CountDocuments(context.TODO(), bson.M{})
	if err != nil {
		return nil, 0, err
	}
	// Pagination
	skip := (page - 1) * limit
	options := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit))
	// Filter
	filter := bson.M{}
	if len(tags) > 0 {
		filter["tags"] = bson.M{"$in": tags}
	}
	cursor, err := br.collection.Find(context.TODO(), filter, options)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		var blog models.Blog
		err := cursor.Decode(&blog)
		if err != nil {
			return nil, 0, err
		}
		blogs = append(blogs, blog)
	}
	return blogs, totalCount, nil
}

func (br *BlogRepository) GetBlogByID(blogID string) (*models.Blog, error) {
	var blog models.Blog
	err := br.collection.FindOne(context.TODO(), bson.M{"blog_id": blogID}).Decode(&blog)
	if err != nil {
		return nil, err
	}
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &blog, nil
}

func (br *BlogRepository) InsertBlog(blog *models.Blog) error {
	_, err := br.collection.InsertOne(context.TODO(), blog)
	if err != nil {
		return err
	}
	return nil
}

func (br *BlogRepository) UpdateBlog(blog *models.Blog) error {
	_, err := br.collection.UpdateOne(context.TODO(), bson.M{"blog_id": blog.BlogID}, bson.M{"$set": blog})
	if err != nil {
		return err
	}
	return nil
}

func (br *BlogRepository) DeleteBlog(blogID string) error {
	_, err := br.collection.DeleteOne(context.TODO(), bson.M{"blog_id": blogID})
	if err != nil {
		return err
	}
	return nil
}
