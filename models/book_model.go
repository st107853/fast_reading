package models

import (
	"time"

	"gorm.io/gorm"
)

type Book struct {
	gorm.Model

	Name        string    `json:"name" gorm:"not null"`
	Author      string    `json:"author" gorm:"not null"`
	ReleaseDate time.Time `json:"release_date"`
	Released    bool      `json:"released" gorm:"default:false;not null"`

	CreatorUserID uint `json:"creator_user_id"`

	Chapters    []Chapter `json:"chapters" gorm:"foreignKey:BookID;constraint:OnDelete:CASCADE;"`
	Labels      []*Label  `json:"labels" gorm:"many2many:book_labels;"`
	FavoritedBy []*User   `json:"favorited_by" gorm:"many2many:user_favorites;"`
}

type Chapter struct {
	gorm.Model

	BookID uint   `json:"book_id" gorm:"book_id"` // Foreign Key
	Title  string `json:"title" gorm:"title"`
	Text   string `json:"text" gorm:"text"`
	Order  int    `json:"order" gorm:"chapter_order"`
}

type BookResponse struct {
	ID     uint   `json:"id" gorm:"_id"`
	Name   string `json:"name" gorm:"name"`
	Author string `json:"author" gorm:"author"`
}

type Label struct {
	gorm.Model

	Name string `json:"name" gorm:"unique;not null"`
}
