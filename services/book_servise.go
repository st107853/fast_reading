package services

import (
	"mime/multipart"

	"github.com/st107853/fast_reading/models"
)

type BookService interface {
	InsertBook(input models.Book, file *multipart.FileHeader, creatorUserID uint) (uint, error)
	FindBookByID(bookId uint) (models.GetBook, error)
	FindBooksByCreatorID(creatorId uint) ([]models.BookBase, error)
	FindFavoriteBooksByUserID(userId uint) ([]models.BookBase, error)
	FindStartedBooks(userID uint) ([]models.BookBase, error)
	InsertChapter(chapter models.Chapter) (uint, error)
	FindChapterByID(id string) (models.Chapter, error)
	FindBooksChapterByIDs(bookId, chapterId uint) (models.ChapterResponse, error)
	DeleteAll() error
	DeleteBook(bookId uint) error
	DeleteChapter(chapterId string) error
	ListAllBooks() ([]models.BookBase, error)
	ListAllLabels() ([]*models.Label, error)
	ListLastReleased(n int) ([]models.Book, error)
	ReleaseBook(bookId uint) error
	UpdateBook(bookId uint, file *multipart.FileHeader, book models.Book) (models.Book, error)
	UpdateChapter(chapterId string, chapter models.Chapter) (models.Chapter, error)
	AddLabel(bookId uint, labelIds []uint) error
	SearchBooks(keyword string, labelIds []uint, readingOnly bool, userID uint) ([]models.BookBase, error)
}
