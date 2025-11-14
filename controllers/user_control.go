package controllers

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/st107853/fast_reading/models"
	"github.com/st107853/fast_reading/services"
)

var userPage = template.Must(template.New("user_page.html").Funcs(template.FuncMap{
	"extractNumericPart": extractNumericPart,
}).ParseFiles("./static/user_page.html"))

type UserController struct {
	userService services.UserService
	bookService services.BookService
}

type UserData struct {
	Name           string
	FavouriteBooks []models.Book
	CreatedBooks   []models.Book
}

func NewUserController(userService services.UserService, bookService services.BookService) UserController {
	return UserController{userService: userService, bookService: bookService}
}

func (uc *UserController) GetMe(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(*models.User)

	//var user = models.FilteredResponse(currentUser)

	// Fetch books created by the current user. If the book service fails,
	// log the error and render the page with favorites only.
	var created []models.Book
	if uc.bookService != nil {
		if cb, err := uc.bookService.FindBooksByCreatorID(currentUser.ID); err == nil {
			created = cb
		} else {
			// don't break the page if created-books query fails; render favorites.
			ctx.Error(err)
		}
	}

	var favorite []models.Book
	if uc.bookService != nil {
		if fb, err := uc.bookService.FindFavoriteBooksByUserEmail(currentUser.ID); err == nil {
			favorite = fb
		} else {
			// don't break the page if favorite-books query fails; render created.
			ctx.Error(err)
		}
	}

	data := UserData{
		Name:           currentUser.Name,
		FavouriteBooks: favorite,
		CreatedBooks:   created,
	}

	// Execute the template and write the output to the response writer
	if err := userPage.Execute(ctx.Writer, data); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"user": models.FilteredResponse(currentUser)}})
}
