package controllers

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/st107853/fast_reading/models"
	"github.com/st107853/fast_reading/services"
)

// Pages html
var mainPage = template.Must(template.New("main_page.html").ParseFiles("./static/main_page.html"))
var bookPage = template.Must(template.New("book_page.html").ParseFiles("./static/book_page.html"))
var bookChapter = template.Must(template.New("book_chapter.html").ParseFiles("./static/book_chapter.html"))
var addBook = template.Must(template.New("create_book_face.html").ParseFiles("./static/create_book_face.html"))
var addBookChapter = template.Must(template.New("create_book_chapter.html").ParseFiles("./static/create_book_chapter.html"))

type BookData struct {
	Title  string
	Books  []models.SmallBookResponse
	Labels []models.Label
}

type BookController struct {
	bookService services.BookService
	userService services.UserService
}

func NewBookController(bookService services.BookService, userService services.UserService) BookController {
	return BookController{bookService, userService}
}

func (bc *BookController) ListAllBooks(c *gin.Context) {
	keyword := c.Query("keyword")
	labelIDsString := c.Query("labels")

	// Transform labelIDsString into []uint
	var labelIDs []uint
	if labelIDsString != "" {
		idStrings := strings.Split(labelIDsString, ",")
		for _, s := range idStrings {
			if id, err := strconv.ParseUint(s, 10, 32); err == nil {
				labelIDs = append(labelIDs, uint(id))
			}
		}
	}

	books, err := bc.bookService.SearchBooks(keyword, labelIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Search failed: " + err.Error()})
		return
	}

	// Return marshaled JSON data
	jsonData, err := models.MarshalBookList(books)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal books: " + err.Error()})
		return
	}

	c.Data(http.StatusOK, "application/json", jsonData)
}

func (bc *BookController) AllBooks(c *gin.Context) {
	var books []models.SmallBookResponse

	books, err := bc.bookService.ListAllBooks()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	labels, err := bc.bookService.ListAllLabels()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	data := BookData{
		Title:  "All what we have",
		Books:  books,
		Labels: labels,
	}

	// Execute the template and write the output to the response writer
	if err := mainPage.Execute(c.Writer, data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

func (bc *BookController) CreateBook(c *gin.Context) {
	var input models.Book
	// c.ShouldBind() using binding/form for multipart/form-data
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data or missing fields: " + err.Error()})
		return
	}

	// Get the file from the form input
	// Note: We do not return an error if the file is missing
	file, err := c.FormFile("cover_image")
	if err != nil && err != http.ErrMissingFile {
		// Error, if it's something other than missing file
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve cover file: " + err.Error()})
		return
	}

	cookie, err := c.Cookie("user_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID cookie not found " + err.Error()})
		return
	}
	creatorUserID, err := strconv.ParseUint(cookie, 10, 32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	bookID, serviceErr := bc.bookService.InsertBook(input, file, uint(creatorUserID))

	if serviceErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create book: %s", serviceErr.Error())})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Book created successfully",
		"book_id": bookID,
	})
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

func (bc *BookController) EditBookChapter(c *gin.Context) {
	id := c.Param("chapter_id")
	var chapter models.Chapter

	chapter, err := bc.bookService.FindChapterByIDStr(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Execute the bookPage template and write the output to the response writer
	if err := addBookChapter.Execute(c.Writer, chapter); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

func (bc *BookController) UpdateBookChapter(c *gin.Context) {
	idParam := c.Param("chapter_id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chapter ID"})
		return
	}

	var inputData models.Chapter
	if err := c.ShouldBindJSON(&inputData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedChapter, err := bc.bookService.UpdateChapter(uint(id), inputData)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedChapter)
}

func (bc *BookController) ReleaseBook(c *gin.Context) {
	// idParam := c.Param("book_id")
	// bookId, err := strconv.ParseUint(idParam, 10, 64)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
	// 	return
	// }

	// book, err := bc.bookService.FindBookByID(idParam)
	// if err != nil {
	// 	c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	// 	return
	// }

	// book.Released = true
	// updatedBook, err := bc.bookService.UpdateBook(uint(bookId), nil, book)
	// if err != nil {
	// 	c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	// 	return
	// }

	// c.JSON(http.StatusOK, updatedBook)
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
	err = bc.userService.AddBookToFavoriteBooks(cookie, strconv.FormatUint(uint64(book.ID), 10))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add book to favorite books " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, book.Name)
}

// GetBooks retrieves all books from the database
func (bc *BookController) GetBooksByName(c *gin.Context) {
	var books []models.SmallBookResponse
	name := c.Param("name")
	books, err := bc.bookService.FindAllByName(name)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	labels, err := bc.bookService.ListAllLabels()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	data := BookData{
		Title:  "All what we have",
		Books:  books,
		Labels: labels,
	}

	// Execute the template and write the output to the response writer
	if err := mainPage.Execute(c.Writer, data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

// GetCreatedBooks retrieves all books created by the user
func (bc *BookController) GetCreatedBooks(c *gin.Context) {
	var books []models.SmallBookResponse
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

	labels, err := bc.bookService.ListAllLabels()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	data := BookData{
		Title:  "All what we have",
		Books:  books,
		Labels: labels,
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
	var book models.GetBook

	book, err := bc.bookService.FindBookByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
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
		book.IsFavorited, err = bc.userService.IsBookFavorited(uint(userid), book.BookID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		book.IsCreator = userid == uint64(book.CreatorUserID)
	}

	// Execute the bookPage template and write the output to the response writer
	if err := bookPage.Execute(c.Writer, book); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

// GetChapter retrieves a chapter by its ID
func (bc *BookController) GetChapter(c *gin.Context) {
	chapterId := c.Param("chapter_id")
	bookId := c.Param("book_id")
	var book models.ChapterResponse

	book, err := bc.bookService.FindBooksChapterByIDs(bookId, chapterId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Execute the bookPage template and write the output to the response writer
	if err := bookChapter.Execute(c.Writer, book); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

// UpdateBook handles the request, including file upload and service call.
func (bc *BookController) UpdateBook(c *gin.Context) {
	idParam := c.Param("book_id")
	bookId, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	var input models.Book

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data: " + err.Error()})
		return
	}

	// Get the file from the form input
	file, err := c.FormFile("cover_image")
	if err != nil && err != http.ErrMissingFile {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve cover file: " + err.Error()})
		return
	}

	// Call Service
	updatedBook, serviceErr := bc.bookService.UpdateBook(uint(bookId), file, input)
	if serviceErr != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": serviceErr.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedBook)
}

// DeleteBook deletes a book by its ID
func (bc *BookController) DeleteBook(c *gin.Context) {
	id := c.Param("book_id")

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

// DeleteBook deletes a book by its ID
func (bc *BookController) DeleteChapter(c *gin.Context) {
	id := c.Param("chapter_id")

	idParsed, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := bc.bookService.DeleteChapter(uint(idParsed)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Chapter deleted"})
}

func (bc *BookController) DeleteAllBooks(c *gin.Context) {
	if err := bc.bookService.DeleteAll(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Book deleted"})
}

func (bc *BookController) AddBook(c *gin.Context) {
	labels, err := bc.bookService.ListAllLabels()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	templateData := gin.H{
		"Name":        "Book title*",
		"Author":      "Book author*",
		"Description": "Book description*",
		"Chapters":    nil,
		"AllLabels":   labels,
	}

	if err := addBook.Execute(c.Writer, templateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

func (bc *BookController) AddBookChapter(c *gin.Context) {
	templateData := gin.H{
		"BookID": c.Param("book_id"),
	}
	if err := addBookChapter.Execute(c.Writer, templateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

func (bc *BookController) EditBook(c *gin.Context) {
	id := c.Param("book_id")
	var book models.GetBook
	book, err := bc.bookService.FindBookByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	book.AllLabels, err = bc.bookService.ListAllLabels()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Execute the bookPage template and write the output to the response writer
	if err := addBook.Execute(c.Writer, book); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

func (bc *BookController) AddLabel(c *gin.Context) {
	idParam := c.Param("book_id")
	labelIdParam := c.Param("label_id")

	// Convert idParam and labelIdParam to uint
	bookId, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	labelId, err := strconv.ParseUint(labelIdParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid label ID"})
		return
	}

	err = bc.bookService.AddLabel(uint(bookId), uint(labelId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Label added to book"})
}
