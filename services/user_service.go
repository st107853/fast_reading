package services

import (
	"github.com/st107853/fast_reading/models"
)

type UserService interface {
	FindUserById(string) (*models.User, error)
	FindUserByEmail(string) (*models.User, error)
	AddBookToFavoriteBooks(email string, bookID string) error
	IsBookFavorited(userID uint, bookID uint) (bool, error)
}
