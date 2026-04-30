package models

import (
	"encoding/json"
	"html/template"
	"time"
)

const StaticCoversPath = "/covers/"

type BookBase struct {
	BookID    uint         `uri:"book_id" json:"id" gorm:"column:id;primaryKey"`
	Name      string       `json:"name" form:"name" gorm:"not null"`
	Author    string       `json:"author" form:"author" gorm:"not null"`
	CoverPath template.URL `json:"cover_path" gorm:"column:cover_path"`
}

type Book struct {
	BookBase

	ReleaseDate     time.Time `json:"release_date" form:"release_date"`
	PublicationYear int       `json:"publication_year" form:"publication_year"`
	Released        bool      `json:"released" form:"released" gorm:"default:false;not null"`
	Description     string    `json:"description" form:"description" gorm:"type:text"`

	CreatorUserID uint `json:"creator_user_id"`

	BookLabels  []*Label `json:"book_labels" gorm:"many2many:book_labels;joinForeignKey:book_id; joinReferences:label_id"`
	FavoritedBy []*User  `json:"favorited_by" gorm:"many2many:user_favorites;joinForeignKey:book_id;joinReferences:user_id"`
}

type BookURI struct {
	BookID uint `uri:"book_id" binding:"required"`
}

func (BookBase) TableName() string {
	return "books"
}

type GetBook struct {
	BookBase

	Description     string          `json:"description"`
	PublicationYear int             `json:"publication_year"`
	Chapters        []*Chapter      `json:"chapters" gorm:"foreignKey:BookID;references:BookID"`
	CreatorUserID   uint            `json:"creator_user_id"`
	IsFavorited     bool            `json:"is_favorited" gorm:"-"`
	IsCreator       bool            `json:"is_creator" gorm:"-"`
	BookLabels      []*Label        `json:"book_labels" gorm:"many2many:book_labels;joinForeignKey:book_id;joinReferences:label_id"`
	Progress        ReadingProgress `json:"progress" gorm:"-"`
}

type Chapter struct {
	ChapterID uint `uri:"chapter_id" json:"id" gorm:"column:id"`

	BookID       uint   `json:"book_id" gorm:"column:book_id"`
	Title        string `json:"title" gorm:"column:title"`
	Text         string `json:"text" gorm:"column:text"`
	ChapterOrder int    `json:"chapter_order" gorm:"column:chapter_order"`
}

type ChapterResponse struct {
	BookBase
	Chapter Chapter `json:"chapter"`
}

type Label struct {
	LabelID uint   `json:"id" gorm:"column:id"`
	Name    string `json:"name" gorm:"unique;not null"`
}

func FormatCoverURL(path string) template.URL {
	if path == "" || path == "null" {
		return template.URL("/static/default_cover.png")
	}

	return template.URL(StaticCoversPath + path)
}

func MarshalBookList(books []BookBase) ([]byte, error) {
	for i := range books {
		books[i].CoverPath = FormatCoverURL(string(books[i].CoverPath))
	}
	return json.Marshal(books)
}
