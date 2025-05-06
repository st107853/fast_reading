package services

import "go.mongodb.org/mongo-driver/bson/primitive"

type BookService interface {
	UpdateBook(bookId primitive.ObjectID, book Book) error
	DeleteBook(bookId string) error
	DeleteAll() error
	InsertBook(book Book) error
	BookExist(bookName, bookAuthor string) bool
	FindBookByID(bookID string) (Book, error)
	ListAllBooks() []Book
	FindAll(bookName string) []Book
	FindBook(bookName string) Book
}
