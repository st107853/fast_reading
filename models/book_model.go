package models

import (
	"encoding/json"
	"html/template"
	"time"

	"gorm.io/gorm"
)

type Book struct {
	gorm.Model

	Name            string    `json:"name" form:"name" gorm:"not null"`
	Author          string    `json:"author" form:"author" gorm:"not null"`
	ReleaseDate     time.Time `json:"release_date" form:"release_date"`
	PublicationYear int       `json:"publication_year" form:"publication_year"`
	Released        bool      `json:"released" form:"released" gorm:"default:false;not null"`
	Description     string    `json:"description" form:"description" gorm:"type:text"`
	CoverPath       string    `json:"cover_path" form:"cover_path"`

	CreatorUserID uint `json:"creator_user_id"`

	Labels      []*Label `json:"labels" gorm:"many2many:book_labels;"`
	FavoritedBy []*User  `json:"favorited_by" gorm:"many2many:user_favorites;"`
}

type BookBase struct {
	gorm.Model

	Name          string       `json:"name" gorm:"column:name"`
	Author        string       `json:"author" gorm:"column:author"`
	Description   string       `json:"description" form:"description" gorm:"type:text"`
	BookLabels    []*Label     `json:"labels" gorm:"many2many:book_labels;"`
	CreatorUserID uint         `json:"creator_user_id"`
	CoverPath     string       `json:"cover_path" form:"cover_path"`
	Cover         template.URL `json:"book_cover"`
}

func (BookBase) TableName() string {
	return "books"
}

func (BookResponse) TableName() string {
	return "books"
}

type GetBook struct {
	BookID          uint            `json:"id" gorm:"column:id"`
	Name            string          `json:"name" gorm:"column:name"`
	Author          string          `json:"author" gorm:"column:author"`
	PublicationYear int             `json:"publication_year" gorm:"column:publication_year"`
	Description     string          `json:"description" form:"description" gorm:"type:text"`
	BookLabels      []*Label        `json:"labels" gorm:"many2many:book_labels;"`
	Chapters        []Chapter       `json:"chapter"`
	IsFavorited     bool            `json:"is_favorited"`
	IsCreator       bool            `json:"is_creator"`
	CreatorUserID   uint            `json:"creator_user_id"`
	AllLabels       []Label         `json:"all_labels"`
	Cover           template.URL    `json:"book_cover"`
	Progress        ReadingProgress `json:"progress" gorm:"column:reading_progress"`
}

type Chapter struct {
	gorm.Model

	BookID       uint   `json:"book_id" gorm:"column:book_id"`
	Title        string `json:"title" gorm:"column:title"`
	Text         string `json:"text" gorm:"column:text"`
	ChapterOrder int    `json:"chapter_order" gorm:"column:chapter_order"`
}

type BookResponse struct {
	ID        uint   `json:"id" gorm:"column:id"`
	Name      string `json:"name" gorm:"column:name"`
	Author    string `json:"author" gorm:"column:author"`
	CoverPath string `json:"cover_path" form:"cover_path"`
}

type SmallBookResponse struct {
	ID     uint         `json:"id"`
	Name   string       `json:"name"`
	Author string       `json:"author"`
	Cover  template.URL `json:"book_cover"`
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
	ID     uint         `json:"id"`
	Name   string       `json:"name"`
	Author string       `json:"author"`
	Cover  template.URL `json:"cover_path"`
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
			Cover:  template.URL(book.CoverPath),
		}
	}

	jsonData, err := json.Marshal(items)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}
