package controllers

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/st107853/fast_reading/services"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Initialize the template with the custom function
var mainPage = template.Must(template.New("main_page.html").Funcs(template.FuncMap{
	"extractNumericPart": extractNumericPart,
}).ParseFiles("./static/main_page.html"))

// Book reading page
var bookPage = template.Must(template.New("book_page.html").ParseFiles("./static/book_page.html"))

var addBook = template.Must(template.New("create_book.html").ParseFiles("./static/create_book.html"))

type Data struct {
	Title string
	Books []services.Book
}

type BookController struct {
	bookService services.BookService
}

func NewBookController(bookService services.BookService) BookController {
	return BookController{bookService}
}

// Custom function to extract the numeric part of the ObjectID
func extractNumericPart(id primitive.ObjectID) string {
	// Assuming the numeric part is the last part of the ObjectID
	new := id.Hex()
	return new
}

func (bc *BookController) AllBooks(c *gin.Context) {
	var books []services.Book
	books = bc.bookService.ListAllBooks()

	data := Data{
		Title: "All what we have",
		Books: books,
	}

	// Execute the template and write the output to the response writer
	if err := mainPage.Execute(c.Writer, data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

func (bc *BookController) CreateBook(c *gin.Context) {
	var book services.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the book already exists
	existingBook := bc.bookService.BookExist(book.Name, book.Author)
	if existingBook != false {
		c.JSON(http.StatusConflict, gin.H{"error": "Book already exists"})
		return
	}

	err := bc.bookService.InsertBook(book)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error2": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, book.Name)
}

// GetBooks retrieves all books from the database
func (bc *BookController) GetBooksByName(c *gin.Context) {
	var books []services.Book
	name := c.Param("name")
	books = bc.bookService.FindAll(name)

	data := Data{
		Title: "All what we have",
		Books: books,
	}

	// Execute the template and write the output to the response writer
	if err := mainPage.Execute(c.Writer, data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

// GetBook retrieves a book by its ID
func (bc *BookController) GetBook(c *gin.Context) {
	id := c.Param("id")
	var book services.Book
	book, err := bc.bookService.FindBookByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Execute the bookPage template and write the output to the response writer
	if err := bookPage.Execute(c.Writer, book); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

// UpdateBook updates an existing book
func (bc *BookController) UpdateBook(c *gin.Context) {
	var book services.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := bc.bookService.UpdateBook(book.ID, book); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}
	c.JSON(http.StatusOK, book)
}

// DeleteBook deletes a book by its ID
func (bc *BookController) DeleteBook(c *gin.Context) {
	id := c.Param("id")
	if err := bc.bookService.DeleteBook(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Book deleted"})
}

func (bc *BookController) DeleteAllBooks(c *gin.Context) {
	if err := bc.bookService.DeleteAll(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Book deleted"})
}

func (bc *BookController) AddBook(c *gin.Context) {
	if err := addBook.Execute(c.Writer, nil); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}
