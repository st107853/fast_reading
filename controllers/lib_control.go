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
var continuePage = template.Must(template.New("continue_page.html").ParseFiles("./static/continue_page.html"))

type BookData struct {
	Title        string
	Books        []models.BookBase
	Labels       []*models.Label
	LastReleased []models.Book
}

type BookController struct {
	bookService services.BookService
	userService services.UserService
}

// UpdateLabelsRequest represents the expected JSON structure for updating book labels
type UpdateLabelsRequest struct {
	LabelIDs []bool `json:"label_ids"`
}

func NewBookController(bookService services.BookService, userService services.UserService) BookController {
	return BookController{bookService, userService}
}

func (bc *BookController) ListAllBooks(c *gin.Context) {
	keyword := c.Query("keyword")
	labelIDsString := c.Query("labels")
	readingOnly := c.Query("reading_only") == "true"

	userId, _ := c.Get("UserId")
	uID, _ := userId.(uint)

	if uID == 0 && readingOnly {
		c.Redirect(http.StatusFound, "/library/auth/login")
		return
	}

	var labelIDs []uint
	if labelIDsString != "" {
		idStrings := strings.Split(labelIDsString, ",")
		for _, s := range idStrings {
			if id, err := strconv.ParseUint(s, 10, 32); err == nil {
				labelIDs = append(labelIDs, uint(id))
			}
		}
	}

	books, err := bc.bookService.SearchBooks(keyword, labelIDs, readingOnly, uID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	jsonData, _ := models.MarshalBookList(books)
	c.Data(http.StatusOK, "application/json", jsonData)
}

func (bc *BookController) AllBooks(c *gin.Context) {
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

	lastReleased, err := bc.bookService.ListLastReleased(2)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	data := BookData{
		Title:        "All what we have",
		Books:        books,
		Labels:       labels,
		LastReleased: lastReleased,
	}

	// Execute the template and write the output to the response writer
	if err := mainPage.Execute(c.Writer, data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

func (bc *BookController) ContinueReading(c *gin.Context) {
	userId, _ := c.Get("UserId")
	uID, _ := userId.(uint)

	books, err := bc.bookService.SearchBooks("", []uint{}, true, uID)
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
		Labels: labels,
		Books:  books,
	}

	// Execute the template and write the output to the response writer
	if err := continuePage.Execute(c.Writer, data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

func (bc *BookController) CreateBook(c *gin.Context) {
	var input models.Book
	userId, _ := c.Get("UserId")
	uID, _ := userId.(uint)

	// c.ShouldBind() using binding/form for multipart/form-data
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data or missing fields: " + err.Error()})
		return
	}

	file, err := c.FormFile("cover_image")
	if err != nil && err != http.ErrMissingFile {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve cover file: " + err.Error()})
		return
	}

	bookID, serviceErr := bc.bookService.InsertBook(input, file, uID)

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

	var uri models.BookURI
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	chapter.BookID = uint(uri.BookID)
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

	chapter, err := bc.bookService.FindChapterByID(id)
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
	id := c.Param("chapter_id")

	var inputData models.Chapter
	if err := c.ShouldBindJSON(&inputData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedChapter, err := bc.bookService.UpdateChapter(id, inputData)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedChapter)
}

func (bc *BookController) ReleaseBook(c *gin.Context) {
	var uri models.BookURI
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := bc.bookService.ReleaseBook(uri.BookID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Book released successfully"})
}

func (bc *BookController) BookFavourite(c *gin.Context) {
	var book models.BookBase
	userId, _ := c.Get("UserId")
	uID, _ := userId.(uint)

	if err := c.ShouldBindUri(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	// Call the function to add the book to favorite books
	err := bc.userService.AddBookToFavoriteBooks(uID, book.BookID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add book to favorite books " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, book.Name)
}

func (bc *BookController) BookMark(c *gin.Context) {
	var uri models.ReadingProgress
	userId, _ := c.Get("UserId")
	uID, _ := userId.(uint)

	if uID == 0 {
		return
	}

	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	err := bc.userService.SaveBooksMark(uID, uri.BookID, uri.ChapterID, uri.LastIndex)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save book mark " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book mark saved successfully"})
}

// GetBook retrieves a book by its ID
func (bc *BookController) GetBook(c *gin.Context) {
	var uri models.BookURI
	userId, _ := c.Get("UserId")
	uID, _ := userId.(uint)

	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	book, err := bc.bookService.FindBookByID(uri.BookID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	progress := models.NewReadingProgress()

	if uID != 0 {
		// Check if book is favorited by current user (if authenticated)
		book.IsFavorited, err = bc.userService.IsBookFavorited(uID, book.BookID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		progress = bc.userService.GetBooksMark(uID, book.BookID)
		book.IsCreator = userId == book.CreatorUserID
	}

	book.Progress = *progress

	// Execute the bookPage template and write the output to the response writer
	if err := bookPage.Execute(c.Writer, book); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

// GetChapter retrieves a chapter by its ID
func (bc *BookController) GetChapter(c *gin.Context) {
	var uri models.ReadingProgress
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	book, err := bc.bookService.FindBooksChapterByIDs(uri.BookID, uri.ChapterID)
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
	var uri models.BookURI
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
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
	updatedBook, serviceErr := bc.bookService.UpdateBook(uri.BookID, file, input)
	if serviceErr != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": serviceErr.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedBook)
}

// DeleteBook deletes a book by its ID
func (bc *BookController) DeleteBook(c *gin.Context) {
	var uri models.BookURI

	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := bc.bookService.DeleteBook(uri.BookID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Book deleted"})
}

// DeleteBook deletes a book by its ID
func (bc *BookController) DeleteChapter(c *gin.Context) {
	id := c.Param("chapter_id")

	if err := bc.bookService.DeleteChapter(id); err != nil {
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
	var uri models.BookURI

	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	book, err := bc.bookService.FindBookByID(uri.BookID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	labels, err := bc.bookService.ListAllLabels()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	templateData := gin.H{
		"Book":      book,
		"AllLabels": labels,
	}

	// Execute the bookPage template and write the output to the response writer
	if err := addBook.Execute(c.Writer, templateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

func (bc *BookController) AddLabel(c *gin.Context) {
	var uri models.BookURI
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var req UpdateLabelsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON body"})
		return
	}

	var labelIDs []uint
	for i := 1; i < len(req.LabelIDs); i++ {
		if req.LabelIDs[i] {
			labelIDs = append(labelIDs, uint(i))
		}
	}

	err := bc.bookService.AddLabel(uri.BookID, labelIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Labels updated successfully"})
}
