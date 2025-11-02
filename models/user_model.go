package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `json:"name" gorm:"not null"`
	Email    string `json:"email" gorm:"unique;not null"`
	Password string `json:"-" gorm:"not null"`
	Role     string `json:"role" gorm:"default:'user';not null"`
	Verified bool   `json:"verified" gorm:"default:false;not null"`

	FavoriteBooks []*Book `json:"favorite_books" gorm:"many2many:user_favorites;"`
}

// SignUpInput specify the fields required to register a new user.
type SignUpInput struct {
	Name            string    `json:"name" gorm:"name"`
	Email           string    `json:"email" gorm:"email" binding:"required"`
	Password        string    `json:"password" gorm:"password" binding:"required,min=8"`
	PasswordConfirm string    `json:"passwordConfirm" gorm:"passwordConfirm,omitempty" binding:"required"`
	Role            string    `json:"role" gorm:"role"`
	Verified        bool      `json:"verified" gorm:"verified"`
	CreatedAt       time.Time `json:"created_at" gorm:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"updated_at"`
}

// SignInInput specify the fields required to sign in a user.
type SignInInput struct {
	Email    string `json:"email" gorm:"email" binding:"required"`
	Password string `json:"password" gorm:"password" binding:"required"`
}

type DBResponse struct {
	ID              uint           `json:"id" gorm:"_id"`
	Name            string         `json:"name" gorm:"name"`
	Email           string         `json:"email" gorm:"email"`
	Password        string         `json:"password" gorm:"password"`
	PasswordConfirm string         `json:"passwordConfirm,omitempty" gorm:"passwordConfirm,omitempty"`
	Favourite       []BookResponse `json:"favourite" gorm:"-"`
	Role            string         `json:"role" gorm:"role"`
	Verified        bool           `json:"verified" gorm:"verified"`
	CreatedAt       time.Time      `json:"created_at" gorm:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at" gorm:"updated_at"`
}

// UserResponse specify the fields that should be included in the JSON response.
type UserResponse struct {
	ID        uint      `json:"id,omitempty" gorm:"_id,omitempty"`
	Name      string    `json:"name,omitempty" gorm:"name,omitempty"`
	Email     string    `json:"email,omitempty" gorm:"email,omitempty"`
	Role      string    `json:"role,omitempty" gorm:"role,omitempty"`
	CreatedAt time.Time `json:"created_at" gorm:"created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"updated_at"`
}

func FilteredResponse(user *DBResponse) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
