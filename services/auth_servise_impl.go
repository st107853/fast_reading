package services

import (
	"context"
	"errors"
	"fmt"
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
	// 1️⃣ Normalize and prepare data
	user.Email = strings.ToLower(user.Email)
	user.PasswordConfirm = ""
	user.Verified = true
	user.Role = "user"
	user.CreatedAt = time.Now()
	user.UpdatedAt = user.CreatedAt

	// 2️⃣ Hash password
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	user.Password = hashedPassword

	// 3️⃣ Check for existing user (email must be unique)
	var existing models.User
	if err := uc.collection.WithContext(uc.ctx).
		Where("email = ?", user.Email).
		First(&existing).Error; err == nil {
		return nil, fmt.Errorf("user with that email already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}

	// 4️⃣ Create user
	newUser := models.User{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
		Role:     user.Role,
		Verified: user.Verified,
	}

	if err := uc.collection.WithContext(uc.ctx).Create(&newUser).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// 5️⃣ Prepare response (no password)
	dbResponse := &models.DBResponse{
		ID:        newUser.ID,
		Name:      newUser.Name,
		Email:     newUser.Email,
		Role:      newUser.Role,
		Verified:  newUser.Verified,
		CreatedAt: newUser.CreatedAt,
	}

	return dbResponse, nil
}

// SignInUser authenticates a user based on their credentials.
// (Currently not implemented.)
func (uc *AuthServiceImpl) SignInUser(*models.SignInInput) (*models.DBResponse, error) {
	return nil, nil
}
