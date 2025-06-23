package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/st107853/fast_reading/models"
	"github.com/st107853/fast_reading/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// AuthServiceImpl is the implementation of the AuthService interface.
type AuthServiceImpl struct {
	collection *mongo.Collection
	ctx        context.Context
}

// NewAuthService creates a new instance of AuthServiceImpl.
func NewAuthService(collection *mongo.Collection, ctx context.Context) AuthService {
	return &AuthServiceImpl{collection, ctx}
}

// SignUpUser registers a new user in the database.
// It hashes the user's password, sets default values, and ensures the email is unique.
func (uc *AuthServiceImpl) SignUpUser(user *models.SignUpInput) (*models.DBResponse, error) {
	// Set default values for the user.
	user.CreatedAt = time.Now()
	user.UpdatedAt = user.CreatedAt
	user.Email = strings.ToLower(user.Email) // Normalize email to lowercase.
	user.PasswordConfirm = ""
	user.Favourite = []models.BookResponse{}
	user.Written = []models.BookResponse{}
	user.Verified = true
	user.Role = "user"

	// Hash the user's password.
	hashedPassword, _ := utils.HashPassword(user.Password)
	user.Password = hashedPassword

	// Add the new user to the database.
	res, err := uc.collection.InsertOne(uc.ctx, &user)
	if err != nil {
		// Handle duplicate email error.
		if er, ok := err.(mongo.WriteException); ok && er.WriteErrors[0].Code == 11000 {
			return nil, fmt.Errorf("asi: user with that email already exist")
		}
		return nil, err
	}

	// Create a unique index for the email field.
	opt := options.Index()
	opt.SetUnique(true)
	index := mongo.IndexModel{Keys: bson.M{"email": 1}, Options: opt}
	if _, err := uc.collection.Indexes().CreateOne(uc.ctx, index); err != nil {
		return nil, fmt.Errorf("asi: could not create index for email. error: %w", err)
	}

	// Retrieve the newly created user from the database.
	var newUser *models.DBResponse
	query := bson.M{"_id": res.InsertedID}
	err = uc.collection.FindOne(uc.ctx, query).Decode(&newUser)
	if err != nil {
		return nil, fmt.Errorf("asi: %w", err)
	}

	return newUser, nil
}

// SignInUser authenticates a user based on their credentials.
// (Currently not implemented.)
func (uc *AuthServiceImpl) SignInUser(*models.SignInInput) (*models.DBResponse, error) {
	return nil, nil
}
