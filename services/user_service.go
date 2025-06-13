package services

import (
	"github.com/st107853/fast_reading/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserService interface {
	FindUserById(string) (*models.DBResponse, error)
	FindUserByEmail(string) (*models.DBResponse, error)
	AddBookToCreatedBooks(email string, bookId primitive.ObjectID, bookName, bookAuthor string) error
	DeleteBookFromCreatedBooks(email, bookId string) error
}
