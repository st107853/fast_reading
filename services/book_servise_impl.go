package services

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Book struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Name        string             `json:"name" binding:"required" bson:"name"`
	Author      string             `json:"author" binding:"required" bson:"author"`
	Text        string             `json:"text" binding:"required" bson:"text"`
	ReleaseDate time.Time          `json:"date" bson:"date"`
}

type BookServiseImpl struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewBookService(collection *mongo.Collection, ctx context.Context) *BookServiseImpl {
	return &BookServiseImpl{collection: collection, ctx: ctx}
}

// InsertBook inserts a new book into the database.
func (bs *BookServiseImpl) InsertBook(book Book) (error, primitive.ObjectID) {
	// Generate a new unique ObjectID for the book
	book.ID = primitive.NewObjectID()

	_, err := bs.collection.InsertOne(bs.ctx, book)
	if err != nil {
		return fmt.Errorf("bsi: failed to insert book: %w", err), primitive.NilObjectID
	}

	return nil, book.ID
}

func (bs *BookServiseImpl) BookExist(bookName, bookAuthor string) (bool, error) {
	filter := bson.M{"name": bookName, "author": bookAuthor}

	count, err := bs.collection.CountDocuments(bs.ctx, filter)
	if err != nil {
		return false, fmt.Errorf("bsi: failed to count documents: %w", err)
	}

	return count > 0, nil
}

func (bs *BookServiseImpl) FindBookByID(bookID string) (Book, error) {
	var book Book

	// Convert the string ID to a primitive.ObjectID
	id, err := primitive.ObjectIDFromHex(bookID)
	if err != nil {
		return book, fmt.Errorf("bsi: invalid ObjectID: %w", err)
	}

	filter := bson.M{"_id": id}

	err = bs.collection.FindOne(bs.ctx, filter).Decode(&book)
	if err != nil {
		return book, fmt.Errorf("bsi: no documents found: %w", err)
	}

	return book, nil
}

// DeleteAll implements BookService.
func (bs *BookServiseImpl) DeleteAll() error {
	_, err := bs.collection.DeleteMany(bs.ctx, bson.M{}, nil)
	if err != nil {
		return fmt.Errorf("bsi: failed to delete all books: %w", err)
	}

	return err
}

// DeleteBook implements BookService.
func (bs *BookServiseImpl) DeleteBook(bookId string) error {
	id, err := primitive.ObjectIDFromHex(bookId)
	if err != nil {
		return fmt.Errorf("bsi: invalid ObjectID: %w", err)
	}

	filter := bson.M{"_id": id}

	result, err := bs.collection.DeleteOne(bs.ctx, filter)
	if err != nil {
		return fmt.Errorf("bsi: failed to delete book: %w", err)
	}

	fmt.Println("Deleted record: ", id.Hex(), " with result: ", result.DeletedCount)
	return nil
}

func (bs *BookServiseImpl) ListAllBooks() ([]Book, error) {
	var books []Book

	cursor, err := bs.collection.Find(bs.ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("bsi: failed to find all books: %w", err)
	}

	err = cursor.All(bs.ctx, &books)
	if err != nil {
		return nil, fmt.Errorf("bsi: failed to decode books: %w", err)
	}

	return books, nil
}

func (bs *BookServiseImpl) FindAll(bookName string) ([]Book, error) {
	var books []Book

	filter := bson.M{"name": bookName}

	cursor, err := bs.collection.Find(bs.ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("bsi: failed to find books by name: %w", err)
	}

	err = cursor.All(bs.ctx, &books)
	if err != nil {
		return nil, fmt.Errorf("bsi: failed to decode books: %w", err)
	}

	return books, nil
}

func (bs *BookServiseImpl) FindBook(bookName string) (Book, error) {
	var book Book

	filter := bson.D{{Key: "name", Value: bookName}}

	err := bs.collection.FindOne(bs.ctx, filter).Decode(&book)
	if err != nil {
		return book, fmt.Errorf("bsi: failed to find book by name: %w", err)
	}

	return book, nil
}

// UpdateBook implements BookService.
func (bs *BookServiseImpl) UpdateBook(bookId primitive.ObjectID, book Book) error {

	filter := bson.M{"id": bookId}
	update := bson.M{"$set": bson.M{"author": book.Author, "text": book.Text, "name": book.Name, "date of release": book.ReleaseDate}}

	result, err := bs.collection.UpdateOne(bs.ctx, filter, update)
	if err != nil {
		return fmt.Errorf("bsi: failed to update book: %w", err)
	}
	fmt.Println("Updated record: ", result)

	return nil
}
