package services

import (
	"context"
	"fmt"

	"github.com/st107853/fast_reading/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

	fmt.Println("I'm at bsi 48 bookID:", bookID) // Debugging line

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

// DeleteAll deletes all books.
func (bs *BookServiseImpl) DeleteAll() error {
	return bs.collection.Exec("DELETE FROM books").Error
}

// DeleteBook delete one book by its ID.
func (bs *BookServiseImpl) DeleteBook(bookId uint) error {
	if err := bs.collection.Unscoped().Delete(&models.Book{}, bookId).Error; err != nil {
		return fmt.Errorf("bsi: failed to hard delete book: %w", err)
	}

	return nil
}

// DeleteChapter deletes one chapter by its ID.
func (bs *BookServiseImpl) DeleteChapter(chapterId uint) error {
	if err := bs.collection.Unscoped().Delete(&models.Chapter{}, chapterId).Error; err != nil {
		return fmt.Errorf("bsi: failed to hard delete chapter: %w", err)
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

// ListAllLabels finds and returns all labels.
func (bs *BookServiseImpl) ListAllLabels() ([]models.Label, error) {
	var labels []models.Label

	err := bs.collection.Find(&labels).Error
	if err != nil {
		return nil, fmt.Errorf("bsi: failed to find all labels: %w", err)
	}

	return labels, nil
}

// ListAllLabels finds and returns all labels.
func (bs *BookServiseImpl) ListAllBooksLabels(bookId string) ([]models.Label, error) {
	var labels []models.Label

	err := bs.collection.Joins("JOIN book_labels ON book_labels.label_id = labels.id").
		// ИСПРАВЛЕНО: Явно указываем таблицу book_labels
		Where("book_labels.book_id = ?", bookId).
		Find(&labels).Error

	if err != nil {
		return nil, fmt.Errorf("bsi: failed to find all labels: %w", err)
	}

	return labels, nil
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

// AddLabel adds a label to a book.
func (bs *BookServiseImpl) AddLabel(bookId, labelId uint) error {
	var book models.Book
	if err := bs.collection.Preload("Labels").First(&book, bookId).Error; err != nil {
		return fmt.Errorf("bsi: book with id %d not found: %w", bookId, err)
	}

	var label models.Label
	if err := bs.collection.First(&label, labelId).Error; err != nil {
		return fmt.Errorf("bsi: label with id %d not found: %w", labelId, err)
	}

	// Check if the label is already associated with the book
	for _, lbl := range book.Labels {
		if lbl.ID == label.ID {
			bs.collection.Model(&book).Association("Labels").Delete(&label)
			return nil // Label already associated, no action needed
		}
	}

	// Append the label to the book's Labels slice
	if err := bs.collection.Model(&book).Association("Labels").Append(&label); err != nil {
		return fmt.Errorf("bsi: failed to add label to book: %w", err)
	}

	return nil
}

func SearchScope(keyword string, labelIDs []uint) func(db *gorm.DB) *gorm.DB {

	return func(db *gorm.DB) *gorm.DB {

		if keyword != "" {
			db = db.Where("name ILIKE ?", "%"+keyword+"%")
		}

		if len(labelIDs) > 0 {

			cleanDB := db.
				Session(&gorm.Session{NewDB: true, SkipDefaultTransaction: true}).
				Omit(clause.Associations)

			subQuery := cleanDB.
				Table("book_labels").
				Select("book_id").
				Where("label_id IN (?)", labelIDs).
				Group("book_id")

			db = db.Where("id IN (?)", subQuery)
		}

		return db
	}
}

func (bs *BookServiseImpl) SearchBooks(keyword string, labelIDs []uint) ([]models.Book, error) {
	var books []models.Book

	err := bs.collection.Scopes(SearchScope(keyword, labelIDs)).Find(&books).Error

	if err != nil {
		return nil, fmt.Errorf("bsi: failed to search books: %w", err)
	}

	return books, nil
}
