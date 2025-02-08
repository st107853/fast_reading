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
	ID          primitive.ObjectID `json:"id" bson:"_id.omitempty"`
	Name        string             `json:"name" bson:"name"`
	Author      string             `json:"author bson:"author"`
	Text        string             `json:"text" bson:"text"`
	ReleaseDate time.Time          `json:"date of release" bson:"date"`
}

func InsertBook(book Book) error {
	collection := db.Database(dbName).Collection(collName)
	inserted, err := collection.InsertOne(context.TODO(), book)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted a record with id: ", inserted.InsertedID)
	return err
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

func UpdateBook(bookId string, book Book) error {
	id, err := primitive.ObjectIDFromHex(bookId)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": id}
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

	filter := bson.M{"_id": id}

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
	delRes, err := collection.DeleteMany(context.TODO(), bson.D{{}}, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Deleted all records: ", delRes)
	return err
}
