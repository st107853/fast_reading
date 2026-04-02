package services

import (
	"mime/multipart"

	"github.com/st107853/fast_reading/models"
)

type BookService interface {
	InsertBook(input models.Book, file *multipart.FileHeader, creatorUserID uint) (uint, error)
	BookExist(bookName, bookAuthor string) (bool, error)
	FindBookByID(bookID string) (models.GetBook, error)
	FindBooksByCreatorID(creatorID uint) ([]models.SmallBookResponse, error)
	FindFavoriteBooksByUserEmail(userID uint) ([]models.SmallBookResponse, error)
	InsertChapter(chapter models.Chapter) (uint, error)
	FindChapterByID(id uint) (models.ChapterResponse, error)
	FindChapterByIDStr(id string) (models.Chapter, error)
	FindBooksChapterByIDs(bookId, chapterId string) (models.ChapterResponse, error)
	FindChaptersByBookID(bookID uint) ([]models.Chapter, error)
	DeleteAll() error
	DeleteBook(bookId uint) error
	DeleteChapter(chapterId uint) error
	ListAllBooks() ([]models.SmallBookResponse, error)
	ListAllLabels() ([]models.Label, error)
	FindAllByName(bookName string) ([]models.SmallBookResponse, error)
	ReleaseBook(bookId uint) error
	UpdateBook(bookId uint, file *multipart.FileHeader, book models.Book) (models.Book, error)
	UpdateChapter(chapterId uint, chapter models.Chapter) (models.Chapter, error)
	AddLabel(bookId uint, labelIds []uint) error
	ListAllBooksLabels(bookId string) ([]models.Label, error)
	SearchBooks(keyword string, labelIDs []uint) ([]models.Book, error)
}
