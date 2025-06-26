package services

import (
	"github.com/st107853/fast_reading/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookService interface {
	UpdateBook(bookId primitive.ObjectID, book models.Book) error
	DeleteBook(bookId string) error
	DeleteAll() error
	InsertBook(book models.Book) (error, primitive.ObjectID)
	BookExist(bookName, bookAuthor string) (bool, error)
	FindBookByID(bookID string) (models.Book, error)
	ListAllBooks() ([]models.Book, error)
	FindAll(bookName string) ([]models.Book, error)
	FindBook(bookName string) (models.Book, error)
}
