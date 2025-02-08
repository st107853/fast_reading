package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/st107853/fast_reading/models"
)

func AllBooks(c *gin.Context) {
	var books []models.Book
	books = models.ListAllBooks()
	c.JSON(http.StatusOK, books)
}

func CreateBook(c *gin.Context) {
	var book models.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := models.InsertBook(book)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, book)
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
	name := c.Param("name")
	var book models.Book
	book = models.FindBook(name)
	c.JSON(http.StatusOK, book)
}

// UpdateBook updates an existing book
func UpdateBook(c *gin.Context) {
	id := c.Param("id")
	realId := [12]byte{}
	for i, v := range id {
		realId[i] = byte(v)
	}
	var book = models.Book{ID: realId, Name: c.Param("name"), Author: c.Param("author"), Text: c.Param("text")}
	if err := models.UpdateBook(id, book); err != nil {
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
