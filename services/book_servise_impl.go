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
	// Validate creator user id to avoid foreign key constraint violation
	if book.CreatorUserID == 0 {
		return 0, fmt.Errorf("bsi: missing CreatorUserID")
	}

	// Insert the book
	err := bs.collection.Create(&book).Error
	return book.ID, err
}

// BookExist checks if a book exist and return bool.
func (bs *BookServiseImpl) BookExist(bookName, bookAuthor string) (bool, error) {
	var count int64

	err := bs.collection.Model(&models.Book{}).Where("name = ? AND author = ?", bookName, bookAuthor).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("bsi: failed to count documents: %w", err)
	}

	return count > 0, nil
}

// FindBookByID finds and returns book by its ID.
func (bs *BookServiseImpl) FindBookByID(bookID string) (models.Book, error) {
	var book models.Book

	err := bs.collection.First(&book, bookID).Error
	if err != nil {
		return book, fmt.Errorf("bsi: failed to find book by ID: %w", err)
	}

	return book, nil
}

// FindBooksByCreatorID finds and returns books by the creator's ID.
func (bs *BookServiseImpl) FindBooksByCreatorID(creatorID uint) ([]models.Book, error) {
	var books []models.Book

	err := bs.collection.Where("creator_user_id = ?", creatorID).Find(&books).Error
	if err != nil {
		return nil, fmt.Errorf("bsi: failed to find books by creator ID: %w", err)
	}

	return books, nil
}

// FindFavoriteBooksByUserEmail finds and returns favorite books by user ID.
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

// FindChapterByID finds and returns chapter by its ID.
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

// FindChapterByID finds and returns chapter by its ID.
func (bs *BookServiseImpl) FindChapterByIDStr(id string) (models.Chapter, error) {
	var chapter models.Chapter

	err := bs.collection.First(&chapter, id).Error
	if err != nil {
		return chapter, fmt.Errorf("bsi: failed to find chapter by ID: %w", err)
	}

	return chapter, nil
}

// FindBooksChapterByIDs finds n'th book's chapter.
func (bs *BookServiseImpl) FindBooksChapterByIDs(bookId, chapterId string) (models.ChapterResponse, error) {
	var chapterResponse models.ChapterResponse

	err := bs.collection.Where("book_id = ? AND chapter_order = ?", bookId, chapterId).First(&chapterResponse.Chapter).Error
	if err != nil {
		return chapterResponse, fmt.Errorf("bsi: failed to find chapter by ID: %w", err)
	}

	err = bs.collection.First(&chapterResponse.Book, chapterResponse.Chapter.BookID).Error
	if err != nil {
		return chapterResponse, fmt.Errorf("bsi: failed to find chapter by ID: %w", err)
	}

	return chapterResponse, nil
}

// FindChaptersByBookID finds and returns chapters by book ID.
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
	if err := bs.collection.Unscoped().Delete(&models.Book{}, bookId).Error; err != nil {
		return fmt.Errorf("bsi: failed to hard delete book: %w", err)
	}

	return nil
}

// ListAllBooks finds and returns all books.
func (bs *BookServiseImpl) ListAllBooks() ([]models.Book, error) {
	var books []models.Book

	err := bs.collection.Limit(20).Find(&books).Error
	if err != nil {
		return nil, fmt.Errorf("bsi: failed to find all books: %w", err)
	}

	return books, nil
}

// FindAll finds and returns books by its name.
func (bs *BookServiseImpl) FindAll(bookName string) ([]models.Book, error) {
	var books []models.Book

	err := bs.collection.Where("name = ?", bookName).Find(&books).Error
	if err != nil {
		return nil, fmt.Errorf("bsi: failed to find books by name: %w", err)
	}

	return books, nil
}

// FindBook finds and returns book by its name.
func (bs *BookServiseImpl) FindBook(bookName string) (models.Book, error) {
	var book models.Book

	err := bs.collection.Where("name = ?", bookName).First(&book).Error
	if err != nil {
		return book, fmt.Errorf("bsi: failed to find book by name: %w", err)
	}

	return book, nil
}

// UpdateBook find and updates a book's fields.
func (bs *BookServiseImpl) UpdateBook(bookId uint, input models.Book) (models.Book, error) {
	var existingBook models.Book
	if err := bs.collection.First(&existingBook, bookId).Error; err != nil {
		return models.Book{}, fmt.Errorf("bsi: book with id %d not found: %w", bookId, err)
	}

	updateData := map[string]interface{}{
		"name":         input.Name,
		"author":       input.Author,
		"release_date": input.ReleaseDate,
		"description":  input.Description,
		"released":     input.Released,
	}

	if err := bs.collection.Model(&existingBook).Updates(updateData).Error; err != nil {
		return models.Book{}, fmt.Errorf("bsi: failed to update book: %w", err)
	}

	return existingBook, nil
}

// UpdateChapter find and updates a chapter's fields.
func (bs *BookServiseImpl) UpdateChapter(chapterId uint, chapter models.Chapter) (models.Chapter, error) {
	var existingChapter models.Chapter
	if err := bs.collection.First(&existingChapter, chapterId).Error; err != nil {
		return models.Chapter{}, fmt.Errorf("bsi: chapter with id %d not found: %w", chapterId, err)
	}

	updateData := map[string]interface{}{
		"title": chapter.Title,
		"text":  chapter.Text,
	}

	if chapter.ChapterOrder != 0 {
		updateData["chapter_order"] = chapter.ChapterOrder
	}

	if err := bs.collection.Model(&existingChapter).Updates(updateData).Error; err != nil {
		return models.Chapter{}, fmt.Errorf("bsi: failed to update chapter: %w", err)
	}

	return existingChapter, nil
}
