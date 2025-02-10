package models

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Book struct {
	ID          primitive.ObjectID `json:"id" bson:"id"`
	Name        string             `json:"name" binding:"required" bson:"name"`
	Author      string             `json:"author" binding:"required" bson:"author"`
	Text        string             `json:"text" binding:"required" bson:"text"`
	ReleaseDate time.Time          `json:"date of release" bson:"date"`
}

func InsertBook(book Book) error {
	// Generate a new unique ObjectID for the book
	book.ID = primitive.NewObjectID()

	collection := db.Database(dbName).Collection(collName)
	result, err := collection.InsertOne(context.TODO(), book)
	if err != nil {
		return fmt.Errorf("failed to insert book: %v", err)
	}

	fmt.Printf("Inserted a record with id: %v\n", result.InsertedID)
	return nil
}

func InsertMany(books []Book) error {
	// Convert to a slice of interface{}.
	newBooks := make([]interface{}, len(books))
	for i, book := range books {
		newBooks[i] = book
	}

	collection := db.Database(dbName).Collection(collName)
	result, err := collection.InsertMany(context.TODO(), newBooks)

	if err != nil {
		log.Fatal(err)
	}

	log.Println(result)
	return err
}

func UpdateBook(bookId primitive.ObjectID, book Book) error {

	filter := bson.M{"id": bookId}
	update := bson.M{"$set": bson.M{"author": book.Author, "text": book.Text, "name": book.Name, "date of release": book.ReleaseDate}}

	collection := db.Database(dbName).Collection(collName)
	result, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	fmt.Println("Updated record: ", result)

	return nil
}

func DeleteBook(bookId string) error {
	id, err := primitive.ObjectIDFromHex(bookId)
	if err != nil {
		return err
	}

	filter := bson.M{"id": id}

	collection := db.Database(dbName).Collection(collName)
	result, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}

	fmt.Println("Deleted record: ", result)
	return nil
}

func FindBook(bookName string) Book {
	var book Book

	filter := bson.D{{Key: "name", Value: bookName}}

	collection := db.Database(dbName).Collection(collName)
	err := collection.FindOne(context.TODO(), filter).Decode(&book)
	if err != nil {
		log.Fatal(err)
	}

	return book
}

func FindBookByID(bookID string) (Book, error) {
	var book Book

	// Convert the string ID to a primitive.ObjectID
	id, err := primitive.ObjectIDFromHex(bookID)
	if err != nil {
		return book, fmt.Errorf("invalid ObjectID: %v", err)
	}

	filter := bson.M{"id": id}

	collection := db.Database(dbName).Collection(collName)
	err = collection.FindOne(context.TODO(), filter).Decode(&book)
	if err != nil {
		return book, fmt.Errorf("mongo: no documents in result: %v", err)
	}

	return book, nil
}

func FindAll(bookName string) []Book {
	var books []Book

	filter := bson.D{{Key: "name", Value: bookName}}

	collection := db.Database(dbName).Collection(collName)
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

func ListAllBooks() []Book {
	var books []Book

	collection := db.Database(dbName).Collection(collName)
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

func DeleteAll() error {
	collection := db.Database(dbName).Collection(collName)
	delRes, err := collection.DeleteMany(context.TODO(), bson.M{}, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Deleted all records: ", delRes)
	return err
}
