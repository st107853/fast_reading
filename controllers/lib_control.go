package controllers

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/st107853/fast_reading/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Custom function to extract the numeric part of the ObjectID
func extractNumericPart(id primitive.ObjectID) string {
	// Assuming the numeric part is the last part of the ObjectID
	new := id.Hex()
	fmt.Println("id from func: ", new)
	return new
}

// Initialize the template with the custom function
var tmpl = template.Must(template.New("template.html").Funcs(template.FuncMap{
	"extractNumericPart": extractNumericPart,
}).ParseFiles("template.html"))

// Book reading page
var index = template.Must(template.New("index.html").ParseFiles("index.html"))

type Data struct {
	Title string
	Books []models.Book
}

func AllBooks(c *gin.Context) {
	var books []models.Book
	books = models.ListAllBooks()

	data := Data{
		Title: "All what we have",
		Books: books,
	}

	// Execute the template and write the output to the response writer
	if err := tmpl.Execute(c.Writer, data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

func CreateBook(c *gin.Context) {

	var book models.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := models.InsertBook(book)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error2": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, book.Name)
}

// GetBooks retrieves all books from the database
func GetBooks(c *gin.Context) {
	var books []models.Book
	name := c.Param("name")
	books = models.FindAll(name)
	c.JSON(http.StatusOK, books)
}

// GetBook retrieves a book by its ID
func GetBook(c *gin.Context) {
	id := c.Param("id")
	fmt.Println("id: ", id)
	var book models.Book
	book, err := models.FindBookByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Execute the index template and write the output to the response writer
	if err := index.Execute(c.Writer, book); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

// UpdateBook updates an existing book
func UpdateBook(c *gin.Context) {
	var book models.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := models.UpdateBook(book.ID, book); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}
	c.JSON(http.StatusOK, book)
}

// DeleteBook deletes a book by its ID
func DeleteBook(c *gin.Context) {
	id := c.Param("id")
	if err := models.DeleteBook(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Book deleted"})
}

func DeleteAllBooks(c *gin.Context) {
	if err := models.DeleteAll(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Book deleted"})
}
