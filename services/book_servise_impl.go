package services

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"html/template"
	"image"
	_ "image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

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

// InsertBook inserts a new book into the database and saves the cover file if provided.
func (bs *BookServiseImpl) InsertBook(book models.Book, file *multipart.FileHeader, creatorUserID uint) (uint, error) {

	// Set the creator user ID before inserting
	book.CreatorUserID = creatorUserID

	if err := bs.collection.Create(&book).Error; err != nil {
		return 0, fmt.Errorf("failed to insert book record: %w", err)
	}

	bookId := book.ID

	// Save the cover file if provided
	if file != nil {
		ext := filepath.Ext(file.Filename)

		// Cover name = book ID + extension
		coverFileName := fmt.Sprintf("%d%s", bookId, ext)
		targetPath := filepath.Join("covers", coverFileName)

		// Creation of 'covers' directory if it doesn't exist
		if _, statErr := os.Stat("covers"); os.IsNotExist(statErr) {
			if mkdirErr := os.MkdirAll("covers", os.ModePerm); mkdirErr != nil {
				// В случае ошибки создания папки, возвращаем ошибку, но запись книги уже есть
				return bookId, fmt.Errorf("failed to create storage directory for cover: %w", mkdirErr)
			}
		}

		// Copying the file
		src, openErr := file.Open()
		if openErr != nil {
			return bookId, fmt.Errorf("failed to open uploaded file: %w", openErr)
		}
		defer src.Close()

		dst, createErr := os.Create(targetPath)
		if createErr != nil {
			return bookId, fmt.Errorf("failed to create destination file for cover: %w", createErr)
		}
		defer dst.Close()

		if _, copyErr := io.Copy(dst, src); copyErr != nil {
			return bookId, fmt.Errorf("failed to save file content: %w", copyErr)
		}

		// Updating the book record with the cover path
		book.CoverPath = targetPath

		// Saving changes to the DB
		if err := bs.collection.Model(&book).Update("CoverPath", targetPath).Error; err != nil {
			// TODO: In case of an update error delete the saved cover file
			return bookId, fmt.Errorf("failed to update book with CoverPath: %w", err)
		}
	}

	return bookId, nil
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
func (bs *BookServiseImpl) FindBookByID(bookID string) (models.GetBook, error) {
	var base models.BookBase
	var result models.GetBook

	// Find the book base information
	err := bs.collection.First(&base, bookID).Error
	if err != nil {
		return result, fmt.Errorf("bsi: failed to find book by ID: %w", err)
	}

	// Find book's chapters
	err = bs.collection.Where("book_id = ?", bookID).Order("chapter_order ASC").Find(&result.Chapters).Error
	if err != nil {
		return result, fmt.Errorf("bsi: failed to find chapters by book ID: %w", err)
	}

	// Find book's labels
	err = bs.collection.Joins("JOIN book_labels ON book_labels.label_id = labels.id").
		Where("book_labels.book_id = ?", bookID).
		Find(&result.BookLabels).Error

	if err != nil {
		return result, fmt.Errorf("bsi: failed to find all labels: %w", err)
	}

	// Get book's cover
	if base.CoverPath != "" {
		img, err := getCover(base.CoverPath)
		if err == nil {
			// Конвертируем объект изображения в Base64 строку
			base64Str, err := imageToSafeBase64(img)
			if err == nil {
				result.Cover = base64Str // Теперь это строка "data:image/png;base64,..."
			}
		}
	}

	// Map base fields to result
	result.BookID = base.ID
	result.Name = base.Name
	result.Author = base.Author
	result.Description = base.Description
	result.CreatorUserID = base.CreatorUserID
	return result, nil
}

func getCover(filePath string) (image.Image, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	image, _, err := image.Decode(f)
	return image, err
}

// imageToSafeBase64 converts an image.Image to a base64-encoded string wrapped in template.URL.
func imageToSafeBase64(img image.Image) (template.URL, error) {
	if img == nil {
		return "", nil
	}

	buffer := new(bytes.Buffer)
	if err := png.Encode(buffer, img); err != nil {
		return "", err
	}

	encoded := base64.StdEncoding.EncodeToString(buffer.Bytes())
	safeURL := template.URL(fmt.Sprintf("data:image/png;base64,%s", encoded))

	return safeURL, nil
}

// FindBooksByCreatorID finds and returns books by the creator's ID.
func (bs *BookServiseImpl) FindBooksByCreatorID(creatorID uint) ([]models.SmallBookResponse, error) {
	var books []models.BookResponse
	var result []models.SmallBookResponse

	err := bs.collection.Where("creator_user_id = ?", creatorID).Find(&books).Error
	if err != nil {
		return nil, fmt.Errorf("bsi: failed to find all books: %w", err)
	}
	result = make([]models.SmallBookResponse, len(books))

	for i, book := range books {
		result[i] = models.SmallBookResponse{
			ID:     book.ID,
			Name:   book.Name,
			Author: book.Author,
		}
		if book.CoverPath != "" {
			img, err := getCover(book.CoverPath)
			if err == nil {
				base64Str, err := imageToSafeBase64(img)
				if err == nil {
					result[i].Cover = base64Str
				}
			}
		}
	}

	return result, nil
}

// FindFavoriteBooksByUserEmail finds and returns favorite books by user ID.
func (bs *BookServiseImpl) FindFavoriteBooksByUserEmail(userID uint) ([]models.SmallBookResponse, error) {
	var books []models.BookResponse
	var result []models.SmallBookResponse

	err := bs.collection.Joins("JOIN user_favorites ON user_favorites.book_id = books.id").
		Where("user_favorites.user_id = ?", userID).
		Find(&books).Error
	if err != nil {
		return nil, fmt.Errorf("bsi: failed to find favorite books by user ID: %w", err)
	}

	result = make([]models.SmallBookResponse, len(books))

	for i, book := range books {
		result[i] = models.SmallBookResponse{
			ID:     book.ID,
			Name:   book.Name,
			Author: book.Author,
		}
		if book.CoverPath != "" {
			img, err := getCover(book.CoverPath)
			if err == nil {
				base64Str, err := imageToSafeBase64(img)
				if err == nil {
					result[i].Cover = base64Str
				}
			}
		}
	}

	return result, nil
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
func (bs *BookServiseImpl) ListAllBooks() ([]models.SmallBookResponse, error) {
	var books []models.BookResponse
	var result []models.SmallBookResponse

	err := bs.collection.Limit(20).Find(&books).Error
	if err != nil {
		return nil, fmt.Errorf("bsi: failed to find all books: %w", err)
	}
	result = make([]models.SmallBookResponse, len(books))

	for i, book := range books {
		result[i] = models.SmallBookResponse{
			ID:     book.ID,
			Name:   book.Name,
			Author: book.Author,
		}
		if book.CoverPath != "" {
			img, err := getCover(book.CoverPath)
			if err == nil {
				base64Str, err := imageToSafeBase64(img)
				if err == nil {
					result[i].Cover = base64Str
				}
			}
		}
	}

	return result, nil
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
		Where("book_labels.book_id = ?", bookId).
		Find(&labels).Error

	if err != nil {
		return nil, fmt.Errorf("bsi: failed to find all labels: %w", err)
	}

	return labels, nil
}

// FindAll finds and returns books by its name.
func (bs *BookServiseImpl) FindAllByName(bookName string) ([]models.SmallBookResponse, error) {
	var books []models.BookResponse
	var result []models.SmallBookResponse
	err := bs.collection.Where("name = ?", bookName).Find(&books).Error
	if err != nil {
		return nil, fmt.Errorf("bsi: failed to find all books: %w", err)
	}
	result = make([]models.SmallBookResponse, len(books))

	for i, book := range books {
		result[i] = models.SmallBookResponse{
			ID:     book.ID,
			Name:   book.Name,
			Author: book.Author,
		}
		if book.CoverPath != "" {
			img, err := getCover(book.CoverPath)
			if err == nil {
				base64Str, err := imageToSafeBase64(img)
				if err == nil {
					result[i].Cover = base64Str
				}
			}
		}
	}

	return result, nil
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
func (bs *BookServiseImpl) UpdateBook(bookId uint, file *multipart.FileHeader, input models.Book) (models.Book, error) {
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
		existingBook.CoverPath = targetPath
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
