package models

import (
	"encoding/json"
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

	BookID       uint   `json:"book_id" gorm:"column:book_id"`
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

// STRUCTURE FOR JSON OUTPUT
type BookListItem struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Author string `json:"author"`
}

// MarshalBookList converts a slice of Book models into a JSON byte slice,
// containing only the ID, Name, and Author fields.
func MarshalBookList(books []Book) ([]byte, error) {
	items := make([]BookListItem, len(books))

	for i, book := range books {
		items[i] = BookListItem{
			ID:     book.ID,
			Name:   book.Name,
			Author: book.Author,
		}
	}

	jsonData, err := json.Marshal(items)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}
