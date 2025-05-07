package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/st107853/fast_reading/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Book struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"` // Change bson tag to "_id"
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
func (b *BookServiseImpl) InsertBook(book Book) error {
	// Generate a new unique ObjectID for the book
	book.ID = primitive.NewObjectID()

	collection := models.DB.Database(models.DBName).Collection(models.CollName)
	result, err := collection.InsertOne(context.TODO(), book)
	if err != nil {
		return fmt.Errorf("failed to insert book: %v", err)
	}

	fmt.Printf("Inserted a record with id: %v\n", result.InsertedID)
	return nil
}

func (b *BookServiseImpl) BookExist(bookName, bookAuthor string) bool {
	filter := bson.M{"name": bookName, "author": bookAuthor}

	collection := models.DB.Database(models.DBName).Collection(models.CollName)
	count, err := collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}

	return count > 0
}

func (b *BookServiseImpl) FindBookByID(bookID string) (Book, error) {
	var book Book

	// Convert the string ID to a primitive.ObjectID
	id, err := primitive.ObjectIDFromHex(bookID)
	if err != nil {
		return book, fmt.Errorf("invalid ObjectID: %v", err)
	}

	filter := bson.M{"_id": id}

	collection := models.DB.Database(models.DBName).Collection(models.CollName)
	err = collection.FindOne(context.TODO(), filter).Decode(&book)
	if err != nil {
		return book, fmt.Errorf("mongo: no documents in result: %v", err)
	}

	return book, nil
}

// DeleteAll implements BookService.
func (b *BookServiseImpl) DeleteAll() error {
	collection := models.DB.Database(models.DBName).Collection(models.CollName)
	delRes, err := collection.DeleteMany(context.TODO(), bson.M{}, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Deleted all records: ", delRes)
	return err
}

// DeleteBook implements BookService.
func (b *BookServiseImpl) DeleteBook(bookId string) error {
	id, err := primitive.ObjectIDFromHex(bookId)
	if err != nil {
		return err
	}

	filter := bson.M{"id": id}

	collection := models.DB.Database(models.DBName).Collection(models.CollName)
	result, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}

	fmt.Println("Deleted record: ", result)
	return nil
}

func (b *BookServiseImpl) ListAllBooks() []Book {
	var books []Book

	collection := models.DB.Database(models.DBName).Collection(models.CollName)
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	err = cursor.All(context.TODO(), &books)
	if err != nil {
		log.Fatal(err)
	}

	return books
}

func (b *BookServiseImpl) FindAll(bookName string) []Book {
	var books []Book

	filter := bson.M{"name": bookName}

	collection := models.DB.Database(models.DBName).Collection(models.CollName)
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}

	err = cursor.All(context.TODO(), &books)
	if err != nil {
		log.Fatal(err)
	}

	return books
}

func (b *BookServiseImpl) FindBook(bookName string) Book {
	var book Book

	filter := bson.D{{Key: "name", Value: bookName}}

	collection := models.DB.Database(models.DBName).Collection(models.CollName)
	err := collection.FindOne(context.TODO(), filter).Decode(&book)
	if err != nil {
		log.Fatal(err)
	}

	return book
}

// UpdateBook implements BookService.
func (b *BookServiseImpl) UpdateBook(bookId primitive.ObjectID, book Book) error {

	filter := bson.M{"id": bookId}
	update := bson.M{"$set": bson.M{"author": book.Author, "text": book.Text, "name": book.Name, "date of release": book.ReleaseDate}}

	collection := models.DB.Database(models.DBName).Collection(models.CollName)
	result, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	fmt.Println("Updated record: ", result)

	return nil
}
