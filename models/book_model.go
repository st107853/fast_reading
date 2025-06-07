package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Book struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Name        string             `json:"name" binding:"required" bson:"name"`
	Author      string             `json:"author" binding:"required" bson:"author"`
	Text        string             `json:"text" binding:"required" bson:"text"`
	ReleaseDate time.Time          `json:"date" bson:"date"`
}

type BookResponse struct {
	ID     primitive.ObjectID `json:"id" bson:"_id"`
	Name   string             `json:"name" bson:"name"`
	Author string             `json:"author" bson:"author"`
}
