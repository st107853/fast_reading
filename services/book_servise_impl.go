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

	err := bs.collection.First(&book, bookID).Error
	if err != nil {
		return book, fmt.Errorf("bsi: failed to find book by ID: %w", err)
	}

	return book, nil
}

func (bs *BookServiseImpl) FindBooksByCreatorID(creatorID uint) ([]models.Book, error) {
	var books []models.Book

	err := bs.collection.Where("creator_user_id = ?", creatorID).Find(&books).Error
	if err != nil {
		return nil, fmt.Errorf("bsi: failed to find books by creator ID: %w", err)
	}

	return books, nil
}

func (bs *BookServiseImpl) FindFavoriteBooksByUserEmail(userID uint) ([]models.Book, error) {
	var books []models.Book

	err := bs.collection.Joins("JOIN user_favorites ON user_favorites.book_id = books.id").
		Where("user_favorites.user_id = ?", userID).
		Find(&books).Error
	if err != nil {
		return nil, fmt.Errorf("bsi: failed to find favorite books by user ID: %w", err)
	}

	return books, nil
}

// InsertChapter inserts a new chapter into the database and assigns an order if not set.
func (bs *BookServiseImpl) InsertChapter(chapter models.Chapter) (uint, error) {
	// If no order provided, calculate next order for the book
	if chapter.ChapterOrder == 0 {
		var count int64
		if err := bs.collection.Model(&models.Chapter{}).Where("book_id = ?", chapter.BookID).Count(&count).Error; err == nil {
			chapter.ChapterOrder = int(count) + 1
		}
	}

	if err := bs.collection.Create(&chapter).Error; err != nil {
		return 0, fmt.Errorf("bsi: failed to insert chapter: %w", err)
	}

	return chapter.ID, nil
}

func (bs *BookServiseImpl) FindChapterByID(id uint) (models.ChapterResponse, error) {
	var chapterResponse models.ChapterResponse

	err := bs.collection.First(&chapterResponse.Chapter, id).Error
	if err != nil {
		return chapterResponse, fmt.Errorf("bsi: failed to find chapter by ID: %w", err)
	}

	err = bs.collection.First(&chapterResponse.Book, chapterResponse.Chapter.BookID).Error
	if err != nil {
		return chapterResponse, fmt.Errorf("bsi: failed to find chapter by ID: %w", err)
	}

	return chapterResponse, nil
}

func (bs *BookServiseImpl) FindChaptersByBookID(bookID uint) ([]models.Chapter, error) {
	var chapters []models.Chapter

	err := bs.collection.Where("book_id = ?", bookID).Order("chapter_order ASC").Find(&chapters).Error
	if err != nil {
		return nil, fmt.Errorf("bsi: failed to find chapters by book ID: %w", err)
	}

	return chapters, nil
}

// DeleteAll implements BookService.
func (bs *BookServiseImpl) DeleteAll() error {
	return bs.collection.Exec("DELETE FROM books").Error
}

// DeleteBook implements BookService.
func (bs *BookServiseImpl) DeleteBook(bookId uint) error {
	var book models.Book
	if err := bs.collection.Where("id = ?", bookId).First(&book).Error; err != nil {
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
