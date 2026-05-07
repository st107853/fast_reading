package controllers

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/st107853/fast_reading/models"
	"github.com/st107853/fast_reading/services"
)

var userPage = template.Must(template.New("user_page.html").ParseFiles("./static/user_page.html", "./static/template.html"))

type UserController struct {
	userService services.UserService
	bookService services.BookService
}

type UserData struct {
	Name            string
	FavouriteBooks  []models.BookBase
	FavouriteLabels []models.Label
	CreatedBooks    []models.BookBase
	CreatedLabels   []models.Label
}

func NewUserController(userService services.UserService, bookService services.BookService) UserController {
	return UserController{userService: userService, bookService: bookService}
}

func (uc *UserController) GetMe(ctx *gin.Context) {
	var err error
	currentUser := ctx.MustGet("currentUser").(*models.User)
	data := UserData{Name: currentUser.Name}

	if uc.bookService != nil {
		if data.CreatedBooks, data.CreatedLabels, err = uc.bookService.FindBooksByCreatorID(currentUser.ID); err != nil {
			ctx.Error(err)
		}
	}

	if data.FavouriteBooks, data.FavouriteLabels, err = uc.bookService.FindFavoriteBooksByUserID(currentUser.ID); err != nil {
		ctx.Error(err)
	}

	// Execute the template and write the output to the response writer
	if err := userPage.Execute(ctx.Writer, data); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}
