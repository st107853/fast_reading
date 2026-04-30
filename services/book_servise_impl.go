package services

import (
	"context"
	"fmt"
	_ "image/jpeg"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

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

// InsertBook inserts a new book into the database and saves the cover file if provided.
func (bs *BookServiseImpl) InsertBook(book models.Book, file *multipart.FileHeader, creatorUserID uint) (uint, error) {

	fmt.Println("Got book: ", book)

	// Set the creator user ID before inserting
	book.CreatorUserID = creatorUserID

	fmt.Println("Set creator id")

	if err := bs.collection.Create(&book).Error; err != nil {
		return 0, fmt.Errorf("failed to insert book record: %w", err)
	}

	bookId := book.BookID

	fmt.Println("Book inserted with ID: ", bookId)

	// Save the cover file if provided
	if file != nil {

		fmt.Println("Started cover")

		ext := filepath.Ext(file.Filename)

		storageDir := "./covers"                          // Define a storage directory for covers
		coverFileName := fmt.Sprintf("%d%s", bookId, ext) // Cover name = book ID + extension
		dstPath := filepath.Join(storageDir, coverFileName)

		fmt.Println("Cover file name: ", coverFileName)
		fmt.Println("Destination path: ", dstPath)

		if err := saveFileToDisk(file, dstPath); err != nil {
			return bookId, err
		}

		fmt.Println("Cover file saved to disk")

		// Updating the book record with the cover path
		book.CoverPath = models.FormatCoverURL(coverFileName)

		fmt.Println("Formatted cover URL: ", book.CoverPath)

		err := bs.collection.Model(&book).Update("cover_path", coverFileName).Error
		if err != nil {
			os.Remove(dstPath)
			return bookId, fmt.Errorf("failed to update book with cover_path: %w", err)
		}
	}

	return bookId, nil
}

// FindBookByID finds and returns book by its ID.
func (bs *BookServiseImpl) FindBookByID(bookID uint) (models.GetBook, error) {
	var result models.GetBook

	err := bs.collection.Model(&models.Book{}).
		Preload("Chapters", func(db *gorm.DB) *gorm.DB {
			return db.Order("chapter_order ASC")
		}).
		Preload("BookLabels").
		Where("id = ?", bookID).
		First(&result).Error

	if err != nil {
		return result, fmt.Errorf("bsi: failed to find book: %w", err)
	}

	result.CoverPath = models.FormatCoverURL(string(result.CoverPath))

	return result, nil
}

// FindBooksByCreatorID finds and returns books by the creator's ID.
func (bs *BookServiseImpl) FindBooksByCreatorID(creatorID uint) ([]models.BookBase, error) {
	var books []models.BookBase

	err := bs.collection.Where("creator_user_id = ?", creatorID).Find(&books).Error
	if err != nil {
		return nil, fmt.Errorf("bsi: failed to find all books: %w", err)
	}

	return books, nil
}

// FindFavoriteBooksByUserEmail finds and returns favorite books by user ID.
func (bs *BookServiseImpl) FindFavoriteBooksByUserID(userID uint) ([]models.BookBase, error) {
	var books []models.BookBase

	err := bs.collection.Joins("JOIN user_favorites ON user_favorites.book_id = books.id").
		Where("user_favorites.user_id = ?", userID).
		Find(&books).Error
	if err != nil {
		return nil, fmt.Errorf("bsi: failed to find favorite books by user ID: %w", err)
	}

	return books, nil
}

func (bs *BookServiseImpl) FindStartedBooks(userID uint) ([]models.BookBase, error) {
	var books []models.BookBase
	err := bs.collection.Joins("JOIN reading_progress ON reading_progress.book_id = books.id").
		Where("reading_progress.user_id = ?", userID).
		Find(&books).Error
	if err != nil {
		return nil, fmt.Errorf("bsi: failed to find started books by user ID: %w", err)
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

	return chapter.ChapterID, nil
}

// FindChapterByID finds and returns chapter by its ID.
func (bs *BookServiseImpl) FindChapterByID(id string) (models.Chapter, error) {
	var chapter models.Chapter

	err := bs.collection.First(&chapter, id).Error
	if err != nil {
		return chapter, fmt.Errorf("bsi: failed to find chapter by ID: %w", err)
	}

	return chapter, nil
}

// FindBooksChapterByIDs finds n'th book's chapter.
func (bs *BookServiseImpl) FindBooksChapterByIDs(bookId, chapterId uint) (models.ChapterResponse, error) {
	var chapterResponse models.ChapterResponse

	err := bs.collection.Where("book_id = ? AND chapter_order = ?", bookId, chapterId).First(&chapterResponse.Chapter).Error
	if err != nil {
		return chapterResponse, fmt.Errorf("bsi: failed to find chapter by ID: %w", err)
	}

	err = bs.collection.First(&chapterResponse.BookBase, chapterResponse.Chapter.BookID).Error
	if err != nil {
		return chapterResponse, fmt.Errorf("bsi: failed to find chapter by ID: %w", err)
	}

	return chapterResponse, nil
}

// DeleteAll deletes all books.
func (bs *BookServiseImpl) DeleteAll() error {
	return bs.collection.Exec("DELETE FROM books").Error
}

// DeleteBook delete one book by its ID.
func (bs *BookServiseImpl) DeleteBook(bookId uint) error {
	if err := bs.collection.Unscoped().Delete(&models.Book{}, bookId).Error; err != nil {
		os.RemoveAll(filepath.Join("covers", fmt.Sprintf("%d.jpeg", bookId)))
		return fmt.Errorf("bsi: failed to hard delete book: %w", err)
	}

	return nil
}

// DeleteChapter deletes one chapter by its ID.
func (bs *BookServiseImpl) DeleteChapter(chapterId string) error {
	if err := bs.collection.Unscoped().Delete(&models.Chapter{}, chapterId).Error; err != nil {
		return fmt.Errorf("bsi: failed to hard delete chapter: %w", err)
	}

	return nil
}

// ListAllBooks finds and returns all books.
func (bs *BookServiseImpl) ListAllBooks() ([]models.BookBase, error) {
	var books []models.BookBase

	err := bs.collection.Limit(20).Where("released = ?", true).Find(&books).Error
	if err != nil {
		return nil, fmt.Errorf("bsi: failed to find all books: %w", err)
	}

	return books, nil
}

// ListAllLabels finds and returns all labels.
func (bs *BookServiseImpl) ListAllLabels() ([]*models.Label, error) {
	var labels []*models.Label

	err := bs.collection.Find(&labels).Error
	if err != nil {
		return nil, fmt.Errorf("bsi: failed to find all labels: %w", err)
	}

	return labels, nil
}

func (bs *BookServiseImpl) ListLastReleased(n int) ([]models.Book, error) {
	var books []models.Book
	err := bs.collection.Where("released = ?", true).Order("release_date DESC").Limit(n).Find(&books).Error
	if err != nil {
		return nil, fmt.Errorf("bsi: failed to find last released books: %w", err)
	}

	for i := range books {

		// Find book's labels
		err = bs.collection.Joins("JOIN book_labels ON book_labels.label_id = labels.id").
			Where("book_labels.book_id = ?", books[i].BookID).
			Find(&books[i].BookLabels).Error

		if err != nil {
			return books, fmt.Errorf("bsi: failed to find all labels: %w", err)
		}
	}

	return books, nil

}

// ReleaseBook sets the release status of a book to true.
func (bs *BookServiseImpl) ReleaseBook(bookId uint) error {
	var book models.Book

	if err := bs.collection.First(&book, bookId).Error; err != nil {
		return fmt.Errorf("bsi: book not found: %w", err)
	}

	updates := map[string]interface{}{
		"released": !book.Released,
	}

	if book.ReleaseDate.Year() < 2 {
		updates["release_date"] = time.Now()
	}

	if err := bs.collection.Model(&book).Updates(updates).Error; err != nil {
		return fmt.Errorf("bsi: failed to update book status: %w", err)
	}

	return nil
}

// UpdateBook find and updates a book's fields.
func (bs *BookServiseImpl) UpdateBook(bookId uint, file *multipart.FileHeader, input models.Book) (models.Book, error) {
	var existingBook models.Book

	if err := bs.collection.First(&existingBook, bookId).Error; err != nil {
		return models.Book{}, fmt.Errorf("bsi: book with id %d not found: %w", bookId, err)
	}

	updateData := map[string]interface{}{
		"name":             input.Name,
		"author":           input.Author,
		"publication_year": input.PublicationYear,
		"description":      input.Description,
	}

	// Conditionally add CoverPath ONLY if it was set by the controller (i.e., a file was uploaded)
	if input.CoverPath != "" {
		updateData["cover_path"] = input.CoverPath
	}

	if file != nil {
		ext := filepath.Ext(file.Filename)

		// Cover name = book ID + extension
		coverFileName := fmt.Sprintf("%d%s", bookId, ext)
		targetPath := filepath.Join("covers", coverFileName)

		// Creation of 'covers' directory if it doesn't exist
		if _, statErr := os.Stat("covers"); os.IsNotExist(statErr) {
			if mkdirErr := os.MkdirAll("covers", os.ModePerm); mkdirErr != nil {
				return existingBook, fmt.Errorf("failed to create storage directory for cover: %w", mkdirErr)
			}
		}

		// Copying the file
		src, openErr := file.Open()
		if openErr != nil {
			return existingBook, fmt.Errorf("failed to open uploaded file: %w", openErr)
		}
		defer src.Close()

		dst, createErr := os.Create(targetPath)
		if createErr != nil {
			return existingBook, fmt.Errorf("failed to create destination file for cover: %w", createErr)
		}
		defer dst.Close()

		if _, copyErr := io.Copy(dst, src); copyErr != nil {
			return existingBook, fmt.Errorf("failed to save file content: %w", copyErr)
		}

		// Updating the book record with the cover path
		existingBook.CoverPath = models.FormatCoverURL(targetPath)
		// Saving changes to the database
		if err := bs.collection.Model(&existingBook).Update("CoverPath", targetPath).Error; err != nil {
			// TODO: In case of an update error delete the saved cover file
			return existingBook, fmt.Errorf("failed to update book with CoverPath: %w", err)
		}
	}

	if err := bs.collection.Model(&existingBook).Updates(updateData).Error; err != nil {
		return models.Book{}, fmt.Errorf("bsi: failed to update book: %w", err)
	}

	return existingBook, nil
}

// UpdateChapter find and updates a chapter's fields.
func (bs *BookServiseImpl) UpdateChapter(chapterId string, chapter models.Chapter) (models.Chapter, error) {
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

func (bs *BookServiseImpl) AddLabel(bookId uint, labelIds []uint) error {
	var book models.Book

	if err := bs.collection.First(&book, bookId).Error; err != nil {
		return fmt.Errorf("bsi: book with id %d not found: %w", bookId, err)
	}

	fmt.Printf("Adding labels %v to book %d\n", labelIds, bookId)

	var labels []models.Label
	if len(labelIds) > 0 {
		if err := bs.collection.Find(&labels, labelIds).Error; err != nil {
			return fmt.Errorf("bsi: failed to fetch labels: %w", err)
		}
	}

	err := bs.collection.Model(&book).Association("BookLabels").Replace(labels)
	if err != nil {
		return fmt.Errorf("bsi: failed to replace labels for book %d: %w", bookId, err)
	}
	return nil
}

func searchScope(keyword string, labelIDs []uint) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// Search by keyword in book name (case-insensitive)
		if keyword != "" {
			db = db.Where("name ILIKE ?", "%"+keyword+"%")
		}

		// Filter by labels if labelIDs are provided
		if len(labelIDs) > 0 {
			// Создаем подзапрос к таблице связей
			subQuery := db.Session(&gorm.Session{NewDB: true}).
				Table("book_labels").
				Select("book_id").
				Where("label_id IN (?)", labelIDs).
				Group("book_id").
				Having("COUNT(DISTINCT label_id) = ?", len(labelIDs))

			db = db.Where("id IN (?)", subQuery)
		}

		return db
	}
}

func (bs *BookServiseImpl) SearchBooks(keyword string, labelIDs []uint, readingOnly bool, userID uint) ([]models.BookBase, error) {
	var books []models.BookBase

	query := bs.collection.Model(&models.BookBase{})

	if readingOnly {
		query = query.Joins("JOIN reading_progress ON reading_progress.book_id = books.id").
			Where("reading_progress.user_id = ?", userID).
			Distinct("books.*")
	}

	err := query.Scopes(searchScope(keyword, labelIDs)).
		Where("released = ?", true).
		Find(&books).Error

	if err != nil {
		return nil, fmt.Errorf("bsi: failed to search books: %w", err)
	}

	return books, nil
}

func saveFileToDisk(file *multipart.FileHeader, dstPath string) error {
	src, err := file.Open()
	if err != nil {
		return fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return fmt.Errorf("failed to save file content: %w", err)
	}
	return nil
}
