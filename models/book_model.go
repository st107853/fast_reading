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
	Description string    `json:"description" gorm:"type:text"`

	CreatorUserID uint `json:"creator_user_id"`

	Labels      []*Label `json:"labels" gorm:"many2many:book_labels;"`
	FavoritedBy []*User  `json:"favorited_by" gorm:"many2many:user_favorites;"`
}

type Chapter struct {
	gorm.Model

	BookID       uint   `json:"book_id" gorm:"column:book_id"` // Foreign Key
	Title        string `json:"title" gorm:"column:title"`
	Text         string `json:"text" gorm:"column:text"`
	ChapterOrder int    `json:"chapter_order" gorm:"column:chapter_order"`
}

type BookResponse struct {
	ID     string `json:"id" gorm:"column:book_id"`
	Name   string `json:"name" gorm:"column:name"`
	Author string `json:"author" gorm:"column:author"`
}

type ChapterResponse struct {
	Book    Book    `json:"book"`
	Chapter Chapter `json:"chapter"`
}

type Label struct {
	gorm.Model

	Name string `json:"name" gorm:"unique;not null"`
}
