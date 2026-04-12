package services

import (
	"github.com/st107853/fast_reading/models"
)

type UserService interface {
	FindUserById(string) (*models.User, error)
	FindUserByEmail(string) (*models.User, error)
	AddBookToFavoriteBooks(email string, bookId uint) error
	GetBooksMark(userId, bookId uint) *models.ReadingProgress
	SaveBooksMark(userId, bookId, chapterId, lastIndex uint) error
	IsBookFavorited(userId, bookId uint) (bool, error)
}
