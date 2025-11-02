package services

import (
	"context"
	"fmt"

	"github.com/st107853/fast_reading/models"
	"gorm.io/gorm"
)

type BookServiseImpl struct {
	collection *gorm.DB
	ctx        context.Context
}

func NewBookService(collection *gorm.DB, ctx context.Context) *BookServiseImpl {
	return &BookServiseImpl{collection: collection, ctx: ctx}
}

// InsertBook inserts a new book into the database.
func (bs *BookServiseImpl) InsertBook(book models.Book) (uint, error) {
	// Generate a new unique ObjectID for the book
	err := bs.collection.Create(&book).Error
	return book.ID, err
}

func (bs *BookServiseImpl) BookExist(bookName, bookAuthor string) (bool, error) {
	var count int64

	err := bs.collection.Model(&models.Book{}).Where("name = ? AND author = ?", bookName, bookAuthor).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("bsi: failed to count documents: %w", err)
	}

	return count > 0, nil
}

func (bs *BookServiseImpl) FindBookByID(bookID string) (models.Book, error) {
	var book models.Book

	// Convert the string ID to a primitive.ObjectID
	err := bs.collection.First(&book, bookID).Error
	if err != nil {
		return book, fmt.Errorf("bsi: failed to find book by ID: %w", err)
	}

	return book, nil
}

// DeleteAll implements BookService.
func (bs *BookServiseImpl) DeleteAll() error {
	return bs.collection.Exec("DELETE FROM books").Error
}

// DeleteBook implements BookService.
func (bs *BookServiseImpl) DeleteBook(bookId uint) error {
	var book models.Book
	if err := bs.collection.First(&book, bookId).Error; err != nil {
		return fmt.Errorf("bsi: failed to find book by ID: %w", err)
	}

	return bs.collection.Delete(&book).Error
}

func (bs *BookServiseImpl) ListAllBooks() ([]models.Book, error) {
	var books []models.Book

	err := bs.collection.Limit(20).Find(&books).Error
	if err != nil {
		return nil, fmt.Errorf("bsi: failed to find all books: %w", err)
	}

	return books, nil
}

func (bs *BookServiseImpl) FindAll(bookName string) ([]models.Book, error) {
	var books []models.Book

	err := bs.collection.Where("name = ?", bookName).Find(&books).Error
	if err != nil {
		return nil, fmt.Errorf("bsi: failed to find books by name: %w", err)
	}

	return books, nil
}

func (bs *BookServiseImpl) FindBook(bookName string) (models.Book, error) {
	var book models.Book

	err := bs.collection.Where("name = ?", bookName).First(&book).Error
	if err != nil {
		return book, fmt.Errorf("bsi: failed to find book by name: %w", err)
	}

	return book, nil
}

// UpdateBook implements BookService.
func (bs *BookServiseImpl) UpdateBook(bookId uint, book models.Book) error {

	if err := bs.collection.Model(&models.Book{}).Where("id = ?", bookId).Updates(book).Error; err != nil {
		return fmt.Errorf("bsi: failed to update book: %w", err)
	}

	return nil
}
