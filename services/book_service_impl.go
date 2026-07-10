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

type BookServiceImpl struct {
	collection *gorm.DB
	ctx        context.Context
}

func NewBookService(collection *gorm.DB, ctx context.Context) *BookServiceImpl {
	return &BookServiceImpl{collection: collection, ctx: ctx}
}

// InsertBook inserts a new book into the database and saves the cover file if provided.
func (bs *BookServiceImpl) InsertBook(book models.Book, file *multipart.FileHeader, creatorUserID uint) (uint, error) {
	book.CreatorUserID = creatorUserID

	// Wrap the entire DB work in a transaction
	var bookID uint
	err := bs.collection.Transaction(func(tx *gorm.DB) error {

		if err := tx.Create(&book).Error; err != nil {
			return ErrDomainWithMsg("Failed to create book.", err)
		}
		bookID = book.BookID

		if file == nil {
			return nil
		}

		// Prepare the cover path and update inside the same transaction
		ext := filepath.Ext(file.Filename)
		coverFileName := fmt.Sprintf("%d%s", bookID, ext)

		if err := tx.Model(&book).Update("cover_path", coverFileName).Error; err != nil {
			return ErrDomainWithMsg("Failed to update cover path.", err)
		}

		return nil
	})

	if err != nil {
		return 0, ErrInternal(err)
	}

	// Write file ONLY after transaction committed successfully
	if file != nil {
		ext := filepath.Ext(file.Filename)
		coverFileName := fmt.Sprintf("%d%s", bookID, ext)
		dstPath := filepath.Join("./covers", coverFileName)

		if err := saveFileToDisk(file, dstPath); err != nil {
			bs.collection.Model(&models.Book{}).
				Where("id = ?", bookID).
				Update("cover_path", nil)

			return bookID, ErrDomainWithMsg("Book created but cover upload failed", err)
		}
	}

	return bookID, nil
}

// FindBookByID finds and returns book by its ID.
func (bs *BookServiceImpl) FindBookByID(bookID uint) (models.GetBook, error) {
	var result models.GetBook

	err := bs.collection.Model(&models.Book{}).
		Preload("Chapters", func(db *gorm.DB) *gorm.DB {
			return db.Order("chapter_order ASC")
		}).
		Preload("BookLabels").
		Where("id = ?", bookID).
		First(&result).Error

	if err != nil {
		return result, ErrDomainWithMsg("bsi: failed to find book", err)
	}

	return result, nil
}

// FindBooksByCreatorID finds and returns books by the creator's ID.
func (bs *BookServiceImpl) FindBooksByCreatorID(creatorID uint) ([]models.BookBase, []models.Label, error) {
	var books []models.BookBase
	var ids []uint

	bs.collection.Table("books").Where("creator_user_id = ?", creatorID).Pluck("id", &ids)

	err := bs.collection.Where("id IN ?", ids).Find(&books).Error
	if err != nil {
		return nil, nil, ErrDomainWithMsg("bsi: failed to find favorite books by user ID", err)
	}

	labels, err := getLabels(bs, ids)

	return books, labels, err
}

func (bs *BookServiceImpl) IsBookCreator(bookID, userID uint) bool {
	var count int64
	err := bs.collection.Model(&models.Book{}).
		Where("id = ? AND creator_user_id = ?", bookID, userID).
		Count(&count).Error

	if err != nil {
		return false
	}

	return count > 0
}

// FindFavoriteBooksByUserEmail finds and returns favorite books by user ID.
func (bs *BookServiceImpl) FindFavoriteBooksByUserID(userID uint) ([]models.BookBase, []models.Label, error) {
	var books []models.BookBase
	var ids []uint
	bs.collection.Table("user_favorites").Where("user_id = ?", userID).Pluck("book_id", &ids)

	err := bs.collection.Where("id IN ?", ids).Find(&books).Error
	if err != nil {
		return nil, nil, ErrDomainWithMsg("bsi: failed to find favorite books by user ID", err)
	}

	labels, err := getLabels(bs, ids)

	return books, labels, err
}

func (bs *BookServiceImpl) FindStartedBooks(userID uint) ([]models.BookBase, error) {
	var books []models.BookBase
	err := bs.collection.Joins("JOIN reading_progress ON reading_progress.book_id = books.id").
		Where("reading_progress.user_id = ?", userID).
		Find(&books).Error
	if err != nil {
		return nil, ErrDomainWithMsg("bsi: failed to find started books by user ID", err)
	}

	return books, nil
}

// InsertChapter inserts a new chapter into the database and assigns an order if not set.
func (bs *BookServiceImpl) InsertChapter(chapter models.Chapter) (uint, error) {
	// If no order provided, calculate next order for the book
	if chapter.ChapterOrder == 0 {
		var count int64
		if err := bs.collection.Model(&models.Chapter{}).Where("book_id = ?", chapter.BookID).Count(&count).Error; err == nil {
			chapter.ChapterOrder = int(count) + 1
		}
	}

	if err := bs.collection.Create(&chapter).Error; err != nil {
		return 0, ErrDomainWithMsg("bsi: failed to insert chapter", err)
	}

	return chapter.ChapterID, nil
}

// FindChapterByID finds and returns chapter by its ID.
func (bs *BookServiceImpl) FindChapterByID(id string) (models.Chapter, error) {
	var chapter models.Chapter

	err := bs.collection.First(&chapter, id).Error
	if err != nil {
		return chapter, ErrDomainWithMsg("bsi: failed to find chapter by ID", err)
	}

	return chapter, nil
}

// FindBooksChapterByIDs finds n'th book's chapter.
func (bs *BookServiceImpl) FindBooksChapterByIDs(bookId, chapterId uint) (models.ChapterResponse, error) {
	var chapterResponse models.ChapterResponse

	err := bs.collection.Where("book_id = ? AND chapter_order = ?", bookId, chapterId).First(&chapterResponse.Chapter).Error
	if err != nil {
		return chapterResponse, ErrDomainWithMsg("bsi: failed to find chapter by ID", err)
	}

	err = bs.collection.First(&chapterResponse.BookBase, chapterResponse.Chapter.BookID).Error
	if err != nil {
		return chapterResponse, ErrDomainWithMsg("bsi: failed to find chapter by ID", err)
	}

	return chapterResponse, nil
}

// DeleteAll deletes all books.
func (bs *BookServiceImpl) DeleteAll() error {
	return bs.collection.Exec("DELETE FROM books").Error
}

// DeleteBook delete one book by its ID.
func (bs *BookServiceImpl) DeleteBook(bookId uint) error {
	if err := bs.collection.Unscoped().Delete(&models.Book{}, bookId).Error; err != nil {
		return ErrDomainWithMsg("bsi: failed to hard delete book", err)
	}

	covers, err := filepath.Glob(filepath.Join("covers", fmt.Sprintf("%d.*", bookId)))
	if err == nil {
		for _, coverPath := range covers {
			_ = os.Remove(coverPath)
		}
	}

	return nil
}

// DeleteChapter deletes one chapter by its ID.
func (bs *BookServiceImpl) DeleteChapter(chapterId string) error {
	if err := bs.collection.Unscoped().Delete(&models.Chapter{}, chapterId).Error; err != nil {
		return ErrDomainWithMsg("bsi: failed to hard delete chapter", err)
	}

	return nil
}

// ListAllBooks finds and returns all books.
func (bs *BookServiceImpl) ListAllBooks() ([]models.BookBase, error) {
	var books []models.BookBase

	err := bs.collection.Limit(20).Where("released = ?", true).Find(&books).Error
	if err != nil {
		return nil, ErrDomainWithMsg("bsi: failed to find all books", err)
	}

	return books, nil
}

// ListAllLabels finds and returns all labels.
func (bs *BookServiceImpl) ListAllLabels() ([]*models.Label, error) {
	var labels []*models.Label

	err := bs.collection.Find(&labels).Error
	if err != nil {
		return nil, ErrDomainWithMsg("bsi: failed to find all labels", err)
	}

	return labels, nil
}

func (bs *BookServiceImpl) ListLastReleased(n int) ([]models.Book, error) {
	var books []models.Book
	err := bs.collection.Where("released = ?", true).Order("release_date DESC").Limit(n).Find(&books).Error
	if err != nil {
		return nil, ErrDomainWithMsg("bsi: failed to find last released books", err)
	}

	for i := range books {

		// Find book's labels
		err = bs.collection.Joins("JOIN book_labels ON book_labels.label_id = labels.id").
			Where("book_labels.book_id = ?", books[i].BookID).
			Find(&books[i].BookLabels).Error

		if err != nil {
			return books, ErrDomainWithMsg("bsi: failed to find all labels", err)
		}
	}

	return books, nil

}

// ReleaseBook sets the release status of a book to true.
func (bs *BookServiceImpl) ReleaseBook(bookId uint) error {
	var book models.Book

	if err := bs.collection.First(&book, bookId).Error; err != nil {
		return ErrDomainWithMsg("bsi: book not found", err)
	}

	updates := map[string]interface{}{
		"released": !book.Released,
	}

	if book.ReleaseDate.Year() < 2 {
		updates["release_date"] = time.Now()
	}

	if err := bs.collection.Model(&book).Updates(updates).Error; err != nil {
		return ErrDomainWithMsg("bsi: failed to update book status", err)
	}

	return nil
}

// UpdateBook find and updates a book's fields.
func (bs *BookServiceImpl) UpdateBook(bookId uint, file *multipart.FileHeader, input models.Book) (models.Book, error) {
	var existingBook models.Book

	if err := bs.collection.First(&existingBook, bookId).Error; err != nil {
		return models.Book{}, ErrDomainWithMsg("bsi: book with id %d not found", err)
	}

	// Determine the new cover filename before the transaction,
	// so we can store it in the DB and write the file after commit.
	var newCoverFileName string
	var oldCoverPath string = string(existingBook.CoverPath)

	if file != nil {
		ext := filepath.Ext(file.Filename)
		if ext == "" {
			ext = ".jpg"
		}
		newCoverFileName = fmt.Sprintf("%d%s", bookId, ext)
	}

	// --- Step 1: Update everything inside a transaction ---
	err := bs.collection.Transaction(func(tx *gorm.DB) error {
		updateData := map[string]interface{}{
			"name":             input.Name,
			"author":           input.Author,
			"publication_year": input.PublicationYear,
			"description":      input.Description,
		}

		if newCoverFileName != "" {
			updateData["cover_path"] = newCoverFileName
		}

		if err := tx.Model(&existingBook).Updates(updateData).Error; err != nil {
			return ErrDomainWithMsg("failed to update book fields", err)
		}

		return nil
	})

	if err != nil {
		// DB failed — no file was touched, nothing to clean up
		return models.Book{}, err
	}

	// --- Step 2: Write file ONLY after transaction committed ---
	if file != nil {
		dstPath := filepath.Join("covers", newCoverFileName)

		if err := os.MkdirAll("covers", os.ModePerm); err != nil {
			return existingBook, ErrDomainWithMsg("failed to create covers directory", err)
		}

		if err := saveFileToDisk(file, dstPath); err != nil {
			// File write failed AFTER DB commit — DB has the new filename
			// but file doesn't exist yet. Roll back just the cover_path in DB.
			rollbackErr := bs.collection.
				Model(&existingBook).
				Update("cover_path", oldCoverPath).
				Error
			if rollbackErr != nil {
				// Log both errors — this needs manual intervention
				return existingBook, ErrDomainWithMsg("CRITICAL: file save failed AND cover_path rollback failed — book has stale cover_path in DB", err)
			}
			return existingBook, ErrDomainWithMsg("cover upload failed, cover_path restored to previous value", err)
		}

		// --- Step 3: Clean up the OLD cover file if extension changed ---
		// e.g. old was .png, new is .jpg — delete the orphaned old file
		if oldCoverPath != "" && oldCoverPath != newCoverFileName {
			oldFilePath := filepath.Join("covers", oldCoverPath)
			if removeErr := os.Remove(oldFilePath); removeErr != nil && !os.IsNotExist(removeErr) {
				// Non-fatal: log but don't fail the request
				ErrDomainWithMsg("warning: could not remove old cover file ", removeErr)
			}
		}

		existingBook.CoverPath = models.FormatCoverURL(newCoverFileName)
	}

	return existingBook, nil
}

// UpdateChapter find and updates a chapter's fields.
func (bs *BookServiceImpl) UpdateChapter(chapterId uint, chapter models.Chapter) (models.Chapter, error) {
	var existingChapter models.Chapter
	if err := bs.collection.First(&existingChapter, chapterId).Error; err != nil {
		return models.Chapter{}, ErrDomainWithMsg("bsi: chapter with id %d not found", err)
	}

	updateData := map[string]interface{}{
		"title": chapter.Title,
		"text":  chapter.Text,
	}

	if chapter.ChapterOrder != 0 {
		updateData["chapter_order"] = chapter.ChapterOrder
	}

	if err := bs.collection.Model(&existingChapter).Updates(updateData).Error; err != nil {
		return models.Chapter{}, ErrDomainWithMsg("bsi: failed to update chapter", err)
	}

	return existingChapter, nil
}

func (bs *BookServiceImpl) AddLabel(bookId uint, labelIds []uint) error {
	var book models.Book

	if err := bs.collection.First(&book, bookId).Error; err != nil {
		return ErrDomainWithMsg("bsi: book with id %d not found", err)
	}

	var labels []models.Label
	if len(labelIds) > 0 {
		if err := bs.collection.Find(&labels, labelIds).Error; err != nil {
			return ErrDomainWithMsg("bsi: failed to fetch labels", err)
		}
	}

	err := bs.collection.Model(&book).Association("BookLabels").Replace(labels)
	if err != nil {
		return ErrDomainWithMsg("bsi: failed to replace labels for book %d", err)
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

// Filter code:
// 0 - regular
// 1 - continue reading
// 2 - created
// 3 - favourite

func (bs *BookServiceImpl) SearchBooks(keyword string, labelIDs []uint, filterCode string, userID uint) ([]models.BookBase, error) {
	var books []models.BookBase

	query := bs.collection.Model(&models.BookBase{})

	switch filterCode {
	case "1":
		query = query.Joins("JOIN reading_progress ON reading_progress.book_id = books.id").
			Where("reading_progress.user_id = ?", userID).
			Distinct("books.*")
	case "2":
		query = query.
			Where("creator_user_id = ?", userID).
			Distinct("books.*")
	case "3":
		query = query.Joins("JOIN user_favorites ON user_favorites.book_id = books.id").
			Where("user_favorites.user_id = ?", userID).
			Distinct("books.*")
	default:
		query = query.Where("released = ?", true)
	}

	err := query.Scopes(searchScope(keyword, labelIDs)).Find(&books).Error

	if err != nil {
		return nil, ErrDomainWithMsg("bsi: failed to search books", err)
	}

	return books, nil
}

func saveFileToDisk(file *multipart.FileHeader, dstPath string) error {
	src, err := file.Open()
	if err != nil {
		return ErrDomainWithMsg("failed to open uploaded file", err)
	}
	defer src.Close()

	dst, err := os.Create(dstPath)
	if err != nil {
		return ErrDomainWithMsg("failed to create destination file", err)
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return ErrDomainWithMsg("failed to save file content", err)
	}
	return nil
}

func getLabels(bs *BookServiceImpl, ids []uint) ([]models.Label, error) {
	var labels []models.Label

	err := bs.collection.
		Table("labels").
		Joins("JOIN book_labels ON book_labels.label_id = labels.id").
		Where("book_labels.book_id IN ?", ids).
		Distinct("labels.*").
		Find(&labels).Error

	return labels, err
}
