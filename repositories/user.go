package repositories

import (
	"context"
	"errors"
	"inkinkink111/go-blog-management/db"
	"inkinkink111/go-blog-management/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		collection: db.DB.Collection("users"), // your collection name
	}
}

func (ur *UserRepository) InsertUser(user *models.User) error {
	_, err := ur.collection.InsertOne(context.TODO(), user)
	if err != nil {
		// Check for duplicate email error
		if mongo.IsDuplicateKeyError(err) {
			return errors.New("email already exists")
		}
		return err
	}
	return nil
}

func (ur *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	filter := bson.M{"email": email}
	result := ur.collection.FindOne(context.TODO(), filter)
	var user models.User
	if err := result.Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

// func (db *userRepo) InsertUser(data models.User) error {
// 	// Check if user already exists
// 	filter := bson.M{"email": data.Email}
// 	count, err := db.Collection.CountDocuments(db.Context, filter)
// 	if err != nil {
// 		return errors.New("failed to count documents")
// 	}
// 	if count > 0 {
// 		return errors.New("user already exists")
// 	}
// 	// Insert user
// 	_, err0 := db.Collection.InsertOne(db.Context, data)
// 	if err0 != nil {
// 		return errors.New("failed to insert user")
// 	}

// 	return nil
// }
