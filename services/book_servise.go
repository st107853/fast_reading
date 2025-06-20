package services

import "go.mongodb.org/mongo-driver/bson/primitive"

type BookService interface {
	UpdateBook(bookId primitive.ObjectID, book Book) error
	DeleteBook(bookId string) error
	DeleteAll() error
	InsertBook(book Book) (error, primitive.ObjectID)
	BookExist(bookName, bookAuthor string) (bool, error)
	FindBookByID(bookID string) (Book, error)
	ListAllBooks() ([]Book, error)
	FindAll(bookName string) ([]Book, error)
	FindBook(bookName string) (Book, error)
}
