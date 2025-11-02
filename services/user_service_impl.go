package services

import (
	"context"

	"github.com/st107853/fast_reading/models"
	"gorm.io/gorm"
)

type UserServiceImpl struct {
	collection *gorm.DB
	ctx        context.Context
}

func NewUserServiceImpl(collection *gorm.DB, ctx context.Context) UserService {
	return &UserServiceImpl{collection, ctx}
}

func (us *UserServiceImpl) FindUserById(id string) (*models.User, error) {
	var user *models.User
	if err := us.collection.WithContext(us.ctx).Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (us *UserServiceImpl) FindUserByEmail(email string) (*models.User, error) {
	var user *models.User
	if err := us.collection.WithContext(us.ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (us *UserServiceImpl) AddBookToFavoriteBooks(email string, bookID uint) error {
	var user models.User
	var book models.Book

	if err := us.collection.WithContext(us.ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return err
	}
	if err := us.collection.WithContext(us.ctx).First(&book, bookID).Error; err != nil {
		return err
	}

	// Check if already exists
	if us.collection.Model(&user).Where("id = ?", bookID).Association("FavoriteBooks").Count() > 0 {
		// Remove if exists
		return us.collection.Model(&user).Association("FavoriteBooks").Delete(&book)
	}

	// Otherwise, add
	return us.collection.Model(&user).Association("FavoriteBooks").Append(&book)
}
