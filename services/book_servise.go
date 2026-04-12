package services

import (
	"mime/multipart"

	"github.com/st107853/fast_reading/models"
)

type BookService interface {
	InsertBook(input models.Book, file *multipart.FileHeader, creatorUserID uint) (uint, error)
	FindBookByID(bookId uint) (models.GetBook, error)
	FindBooksByCreatorID(creatorId uint) ([]models.SmallBookResponse, error)
	FindFavoriteBooksByUserEmail(userId uint) ([]models.SmallBookResponse, error)
	InsertChapter(chapter models.Chapter) (uint, error)
	FindChapterByID(id string) (models.Chapter, error)
	FindBooksChapterByIDs(bookId, chapterId uint) (models.ChapterResponse, error)
	DeleteAll() error
	DeleteBook(bookId uint) error
	DeleteChapter(chapterId string) error
	ListAllBooks() ([]models.SmallBookResponse, error)
	ListAllLabels() ([]models.Label, error)
	ListLastReleased(n int) ([]models.BookBase, error)
	ReleaseBook(bookId uint) error
	UpdateBook(bookId uint, file *multipart.FileHeader, book models.Book) (models.Book, error)
	UpdateChapter(chapterId string, chapter models.Chapter) (models.Chapter, error)
	AddLabel(bookId uint, labelIds []uint) error
	SearchBooks(keyword string, labelIds []uint) ([]models.Book, error)
}
