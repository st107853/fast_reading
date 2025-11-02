package controllers

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/st107853/fast_reading/models"
	"github.com/st107853/fast_reading/services"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Initialize the template with the custom function
var mainPage = template.Must(template.New("main_page.html").Funcs(template.FuncMap{
	"extractNumericPart": extractNumericPart,
}).ParseFiles("./static/main_page.html"))

// Book reading page
var bookPage = template.Must(template.New("book_page.html").ParseFiles("./static/book_page.html"))

var addBook = template.Must(template.New("create_book_face.html").ParseFiles("./static/create_book_face.html"))

var addBookChapter = template.Must(template.New("create_book_chapter.html").ParseFiles("./static/create_book_chapter.html"))

type BookData struct {
	Title string
	Books []models.Book
}

type BookController struct {
	bookService services.BookService
	userService services.UserService
}

func NewBookController(bookService services.BookService, userService services.UserService) BookController {
	return BookController{bookService, userService}
}

// Custom function to extract the numeric part of the ObjectID
func extractNumericPart(id primitive.ObjectID) string {
	// Assuming the numeric part is the last part of the ObjectID
	new := id.Hex()
	return new
}

func (bc *BookController) AllBooks(c *gin.Context) {
	var books []models.Book
	books, err := bc.bookService.ListAllBooks()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	data := BookData{
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
	var book models.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the book already exists
	existingBook, err := bc.bookService.BookExist(book.Name, book.Author)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if existingBook {
		c.JSON(http.StatusConflict, gin.H{"error": "Book already exists"})
		return
	}

	cookie, err := c.Cookie("user_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID cookie not found " + err.Error()})
		return
	}

	userid, err := strconv.ParseUint(cookie, 10, 32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	book.CreatorUserID = uint(userid)

	_, err = bc.bookService.InsertBook(book)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Call the function to add the book to created books
	c.JSON(http.StatusCreated, book.Name)
}

func (bc *BookController) BookFavourite(c *gin.Context) {
	var book models.BookResponse
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cookie, err := c.Cookie("email")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID cookie not found " + err.Error()})
		return
	}

	// Call the function to add the book to favorite books
	err = bc.userService.AddBookToFavoriteBooks(cookie, book.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add book to favorite books " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, book.Name)
}

// GetBooks retrieves all books from the database
func (bc *BookController) GetBooksByName(c *gin.Context) {
	var books []models.Book
	name := c.Param("name")
	books, err := bc.bookService.FindAll(name)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	data := BookData{
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
	var book models.Book
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
	var book models.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := book.Model.ID
	if err := bc.bookService.UpdateBook(id, book); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, book)
}

// DeleteBook deletes a book by its ID
func (bc *BookController) DeleteBook(c *gin.Context) {
	cookie, err := c.Cookie("id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID cookie not found " + err.Error()})
		return
	}
	id, err := strconv.ParseUint(cookie, 10, 32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := bc.bookService.DeleteBook(uint(id)); err != nil {
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

func (bc *BookController) AddBookChapter(c *gin.Context) {
	if err := addBookChapter.Execute(c.Writer, nil); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}
