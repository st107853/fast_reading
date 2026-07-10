package services

import (
	"context"
	"strings"
	"time"

	"github.com/st107853/fast_reading/models"
	"github.com/st107853/fast_reading/utils"
	"gorm.io/gorm"
)

// AuthServiceImpl is the implementation of the AuthService interface.
type AuthServiceImpl struct {
	collection *gorm.DB
	ctx        context.Context
}

// NewAuthService creates a new instance of AuthServiceImpl.
func NewAuthService(collection *gorm.DB, ctx context.Context) AuthService {
	return &AuthServiceImpl{collection, ctx}
}

// SignUpUser registers a new user in the database.
// It hashes the user's password, sets default values, and ensures the email is unique.
func (uc *AuthServiceImpl) SignUpUser(user *models.SignUpInput) (*models.DBResponse, error) {
	// Normalize and prepare data
	user.Email = strings.ToLower(user.Email)
	user.PasswordConfirm = ""
	user.Verified = true
	user.Role = "user"
	user.CreatedAt = time.Now()
	user.UpdatedAt = user.CreatedAt

	// Hash password
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return nil, ErrDomainWithMsg("failed to hash password", err)
	}
	user.Password = hashedPassword

	// Check for existing user (email must be unique)
	var count int64

	err = uc.collection.WithContext(uc.ctx).
		Model(&models.User{}).
		Where("email = ?", user.Email).
		Count(&count).Error

	if err != nil {
		return nil, ErrDomainWithMsg("failed to check existing user", err)
	}

	// Если count > 0, значит пользователь существует
	if count > 0 {
		return nil, ErrUser("user with that email already exists", nil)
	}

	// Create user
	newUser := models.User{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
		Role:     user.Role,
		Verified: user.Verified,
	}

	if err := uc.collection.WithContext(uc.ctx).Create(&newUser).Error; err != nil {
		return nil, ErrDomainWithMsg("failed to create user", err)
	}

	// Prepare response (no password)
	dbResponse := &models.DBResponse{
		ID:       newUser.ID,
		Name:     newUser.Name,
		Email:    newUser.Email,
		Role:     newUser.Role,
		Verified: newUser.Verified,
	}

	return dbResponse, nil
}
