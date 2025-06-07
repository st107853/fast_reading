package services

import (
	"context"
	"strings"

	"github.com/st107853/fast_reading/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserServiceImpl struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewUserServiceImpl(collection *mongo.Collection, ctx context.Context) UserService {
	return &UserServiceImpl{collection, ctx}
}

func (us *UserServiceImpl) FindUserById(id string) (*models.DBResponse, error) {
	oid, _ := primitive.ObjectIDFromHex(id)

	var user *models.DBResponse

	query := bson.M{"_id": oid}
	err := us.collection.FindOne(us.ctx, query).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &models.DBResponse{}, err
		}
		return nil, err
	}

	return user, nil
}

func (us *UserServiceImpl) FindUserByEmail(email string) (*models.DBResponse, error) {
	var user *models.DBResponse

	query := bson.M{"email": strings.ToLower(email)}
	err := us.collection.FindOne(us.ctx, query).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &models.DBResponse{}, err
		}
		return nil, err
	}

	return user, nil
}

// AddBookToCreatedBooks appends a book name to the user's CreatedBooks array.
func (us *UserServiceImpl) AddBookToCreatedBooks(email string, bookId primitive.ObjectID, bookName, bookAuthor string) error {
	query := bson.M{"email": strings.ToLower(email)}
	book := models.BookResponse{
		ID:     bookId,
		Name:   bookName,
		Author: bookAuthor,
	}
	update := bson.M{"$push": bson.M{"createdbooks": book}}

	_, err := us.collection.UpdateOne(us.ctx, query, update)
	if err != nil {
		return err
	}
	return nil
}
