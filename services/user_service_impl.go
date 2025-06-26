package services

import (
	"context"
	"fmt"
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
		return nil, fmt.Errorf("usi: %w", err)
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
		return nil, fmt.Errorf("usi: %w", err)
	}

	return user, nil
}

// AddBookToCreatedBooks appends a book name to the user's CreatedBooks array.
func (us *UserServiceImpl) AddBookToCreatedBooks(email string, bookId primitive.ObjectID, bookName, bookAuthor string) error {
	// Ensure book format matches your schema
	book := models.BookResponse{
		ID:     bookId,
		Name:   bookName,
		Author: bookAuthor,
	}

	// First, ensure createdbooks exists and is an array
	_, err := us.collection.UpdateOne(
		us.ctx,
		bson.M{
			"email": email,
			"$or": []bson.M{
				{"createdbooks": bson.M{"$exists": false}},
				{"createdbooks": nil},
			},
		},
		bson.M{
			"$set": bson.M{"createdbooks": bson.A{}},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to initialize createdbooks: %w", err)
	}

	// Now push the book
	_, err = us.collection.UpdateOne(
		us.ctx,
		bson.M{"email": email},
		bson.M{"$push": bson.M{"createdbooks": book}},
	)
	if err != nil {
		return fmt.Errorf("usi: updating user error: %w", err)
	}

	return nil
}

func (us *UserServiceImpl) AddBookToFavoriteBooks(email string, bookId primitive.ObjectID, bookName, bookAuthor string) error {

	// Check if the book is already in the favourite list
	count, err := us.collection.CountDocuments(
		us.ctx,
		bson.M{
			"email":         email,
			"favourite._id": bookId,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to check favourite books: %w", err)
	}
	if count > 0 {
		// Query by user's email
		query := bson.M{"email": email}

		// $pull the book with matching _id from createdbooks
		update := bson.M{
			"$pull": bson.M{
				"favourite": bson.M{
					"_id": bookId,
				},
			},
		}

		_, err = us.collection.UpdateOne(us.ctx, query, update)
		if err != nil {
			return fmt.Errorf("failed to pull book from favourite: %w", err)
		}

		return nil
	}

	// Ensure book format matches your schema
	book := models.BookResponse{
		ID:     bookId,
		Name:   bookName,
		Author: bookAuthor,
	}

	// First, ensure createdbooks exists and is an array
	_, err = us.collection.UpdateOne(
		us.ctx,
		bson.M{
			"email": email,
			"$or": []bson.M{
				{"favourite": bson.M{"$exists": false}},
				{"favourite": nil},
			},
		},
		bson.M{
			"$set": bson.M{"favourite": bson.A{}},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to initialize favourite: %w", err)
	}

	// Now push the book
	_, err = us.collection.UpdateOne(
		us.ctx,
		bson.M{"email": email},
		bson.M{"$push": bson.M{"favourite": book}},
	)
	if err != nil {
		return fmt.Errorf("usi: updating user error: %w", err)
	}

	return nil
}

func (us *UserServiceImpl) DeleteBookFromCreatedBooks(email, bookId string) error {
	// Convert bookId string to ObjectID
	bookIdObj, err := primitive.ObjectIDFromHex(bookId)
	if err != nil {
		return fmt.Errorf("invalid bookId: %w", err)
	}

	// Query by user's email
	query := bson.M{"email": email}

	// $pull the book with matching _id from createdbooks
	update := bson.M{
		"$pull": bson.M{
			"createdbooks": bson.M{
				"_id": bookIdObj,
			},
		},
	}

	_, err = us.collection.UpdateOne(us.ctx, query, update)
	if err != nil {
		return fmt.Errorf("failed to pull book from createdbooks: %w", err)
	}

	return nil
}
