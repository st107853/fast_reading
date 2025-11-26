package services

import (
	"github.com/st107853/fast_reading/models"
)

type BookService interface {
	InsertBook(book models.Book) (uint, error)
	BookExist(bookName, bookAuthor string) (bool, error)
	FindBookByID(bookID string) (models.Book, error)
	FindBooksByCreatorID(creatorID uint) ([]models.Book, error)
	FindFavoriteBooksByUserEmail(userID uint) ([]models.Book, error)
	InsertChapter(chapter models.Chapter) (uint, error)
	FindChapterByID(id uint) (models.ChapterResponse, error)
	FindChapterByIDStr(id string) (models.Chapter, error)
	FindBooksChapterByIDs(bookId, chapterId string) (models.ChapterResponse, error)
	FindChaptersByBookID(bookID uint) ([]models.Chapter, error)
	DeleteAll() error
	DeleteBook(bookId uint) error
	DeleteChapter(chapterId uint) error
	ListAllBooks() ([]models.Book, error)
	ListAllLabels() ([]models.Label, error)
	FindAll(bookName string) ([]models.Book, error)
	FindBook(bookName string) (models.Book, error)
	UpdateBook(bookId uint, book models.Book) (models.Book, error)
	UpdateChapter(chapterId uint, chapter models.Chapter) (models.Chapter, error)
	AddLabel(bookId, LabelId uint) error
	ListAllBooksLabels(bookId string) ([]models.Label, error)
	SearchBooks(keyword string, labelIDs []uint) ([]models.Book, error)
}
