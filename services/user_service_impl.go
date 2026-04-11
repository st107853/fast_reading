package services

import (
	"context"

	"github.com/st107853/fast_reading/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func (us *UserServiceImpl) AddBookToFavoriteBooks(email string, bookID string) error {
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

func (us *UserServiceImpl) SaveBooksMark(userId uint, bookID uint, chapterID uint, lastIndex int) error {
	progress := models.ReadingProgress{
		UserID:    userId,
		BookID:    bookID,
		ChapterID: chapterID,
		LastIndex: lastIndex,
	}

	return us.collection.WithContext(us.ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "book_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"chapter_id", "last_index"}),
	}).Create(&progress).Error
}

func (us *UserServiceImpl) GetBooksMark(userId uint, bookID uint) *models.ReadingProgress {
	var progress models.ReadingProgress
	progress.ChapterID = 1
	us.collection.WithContext(us.ctx).Where("user_id = ? AND book_id = ?", userId, bookID).First(&progress)

	return &progress
}

func (us *UserServiceImpl) IsBookFavorited(userID uint, bookID uint) (bool, error) {
	var user models.User
	if err := us.collection.WithContext(us.ctx).First(&user, userID).Error; err != nil {
		return false, err
	}

	count := us.collection.Model(&user).Where("id = ?", bookID).Association("FavoriteBooks").Count()
	return count > 0, nil
}
