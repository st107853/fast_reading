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

func (us *UserServiceImpl) FindUserById(id uint) (*models.User, error) {
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

func (us *UserServiceImpl) AddBookToFavoriteBooks(id, bookId uint) error {
	var user models.User

	if err := us.collection.WithContext(us.ctx).Where("id = ?", id).First(&user).Error; err != nil {
		return err
	}

	book := models.BookBase{
		BookID: bookId,
	}

	association := us.collection.Model(&user).Association("FavoriteBooks")
	exists := false

	var count int64
	us.collection.Table("user_favorites").Where("user_id = ? AND book_id = ?", user.ID, bookId).Count(&count)
	if count > 0 {
		exists = true
	}

	if exists {
		return association.Delete(&book)
	}

	return association.Append(&book)
}

func (us *UserServiceImpl) SaveBooksMark(userId uint, bookId uint, chapterID uint, lastIndex uint) error {
	progress := models.ReadingProgress{
		UserID:    userId,
		BookID:    bookId,
		ChapterID: chapterID,
		LastIndex: lastIndex,
	}

	return us.collection.WithContext(us.ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "book_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"chapter_id", "last_index"}),
	}).Create(&progress).Error
}

func (us *UserServiceImpl) GetBooksMark(userId uint, bookId uint) *models.ReadingProgress {
	var progress = models.NewReadingProgress()
	us.collection.WithContext(us.ctx).Where("user_id = ? AND book_id = ?", userId, bookId).Limit(1).Find(progress)

	return progress
}

func (us *UserServiceImpl) IsBookFavorited(userID uint, bookId uint) (bool, error) {
	var user models.User
	if err := us.collection.WithContext(us.ctx).First(&user, userID).Error; err != nil {
		return false, err
	}

	count := us.collection.Model(&user).Where("id = ?", bookId).Association("FavoriteBooks").Count()
	return count > 0, nil
}
