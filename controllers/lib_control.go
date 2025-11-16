package controllers

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/st107853/fast_reading/config"
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
var bookChapter = template.Must(template.New("book_chapter.html").ParseFiles("./static/book_chapter.html"))
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

	// Add creator user ID from cookie
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

	book_id, err := bc.bookService.InsertBook(book)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	config, err := config.LoadConfig(".")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "error": "Could not load config " + err.Error()})
		return
	}

	c.SetCookie("book_id", strconv.Itoa(int(book_id)), config.AccessTokenMaxAge*60, "/", "localhost", false, false)

	// Return created book ID and name so client can navigate to chapter creation
	c.JSON(http.StatusCreated, gin.H{"book_id": book_id, "name": book.Name})
}

func (bc *BookController) CreateChapter(c *gin.Context) {
	var chapter models.Chapter
	if err := c.ShouldBindJSON(&chapter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Try to get book id from the URL param first, fallback to cookie
	var bookIdUint uint64
	var err error
	if param := c.Param("book_id"); param != "" {
		bookIdUint, err = strconv.ParseUint(param, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Book ID param " + err.Error()})
			return
		}
	} else {
		cookie, err := c.Cookie("book_id")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Book ID cookie not found " + err.Error()})
			return
		}
		bookIdUint, err = strconv.ParseUint(cookie, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Book ID " + err.Error()})
			return
		}
	}
	chapter.BookID = uint(bookIdUint)
	// Persist chapter via service
	if bc.bookService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "book service not available"})
		return
	}

	id, err := bc.bookService.InsertChapter(chapter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save chapter: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"chapter_id": id})
}

func (bc *BookController) ReleaseBook(c *gin.Context) {
	bookId := c.Param("book_id")
	var book models.Book
	book, err := bc.bookService.FindBookByID(bookId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	book.Released = true
	err = bc.bookService.UpdateBook(book.Model.ID, book)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book released successfully"})
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

// GetCreatedBooks retrieves all books created by the user
func (bc *BookController) GetCreatedBooks(c *gin.Context) {
	var books []models.Book
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

	books, err = bc.bookService.FindBooksByCreatorID(uint(userid))

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
	id := c.Param("book_id")
	isFavorited := false
	isCreator := false
	var book models.Book

	book, err := bc.bookService.FindBookByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	chapters, err := bc.bookService.FindChaptersByBookID(book.Model.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	cookie, err := c.Cookie("user_id")
	if err == nil {
		userid, err := strconv.ParseUint(cookie, 10, 32)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Check if book is favorited by current user (if authenticated)
		isFavorited, err = bc.userService.IsBookFavorited(uint(userid), book.Model.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		isCreator = userid == uint64(book.CreatorUserID)
	}

	// Pass both book and isFavorited to template
	templateData := gin.H{
		"Name":        book.Name,
		"Author":      book.Author,
		"ID":          book.Model.ID,
		"Description": book.Description,
		"Chapters":    chapters,
		"IsFavorited": isFavorited,
		"IsCreator":   isCreator,
	}

	// Execute the bookPage template and write the output to the response writer
	if err := bookPage.Execute(c.Writer, templateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

// GetChapter retrieves a chapter by its ID
func (bc *BookController) GetChapter(c *gin.Context) {
	id := c.Param("chapter_id")
	var chapter models.ChapterResponse

	chapterId, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	chapter, err = bc.bookService.FindChapterByID(uint(chapterId))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Execute the bookPage template and write the output to the response writer
	if err := bookChapter.Execute(c.Writer, chapter); err != nil {
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
	id := c.Param("id")

	idParsed, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := bc.bookService.DeleteBook(uint(idParsed)); err != nil {
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

	templateData := gin.H{
		"Name":        "Book title*",
		"Author":      "Book author*",
		"Description": "Book description*",
		"Chapters":    nil,
	}

	if err := addBook.Execute(c.Writer, templateData); err != nil {
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

func (bc *BookController) EditBook(c *gin.Context) {
	id := c.Param("book_id")
	var book models.Book

	book, err := bc.bookService.FindBookByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Execute the bookPage template and write the output to the response writer
	if err := addBook.Execute(c.Writer, book); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}
