package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/st107853/fast_reading/config"
	"github.com/st107853/fast_reading/controllers"
	"github.com/st107853/fast_reading/models"
	"github.com/st107853/fast_reading/routes"
	"github.com/st107853/fast_reading/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	server *gin.Engine
	ctx    context.Context

	userService         services.UserService
	UserController      controllers.UserController
	UserRouteController routes.UserRouteController

	authService         services.AuthService
	AuthController      controllers.AuthController
	AuthRouteController routes.AuthRouteController

	bookService         services.BookService
	db                  *gorm.DB
	BookRouteController routes.BookRouteController
)

func init() {
	ctx = context.TODO()

	// Initialize GORM via models helper (returns *gorm.DB)
	gdb, err := models.OpenDbConnection()
	if err != nil {
		log.Fatalf("models.OpenDbConnection failed: %v", err)
	}
	if gdb == nil {
		log.Fatal("models.OpenDbConnection returned nil *gorm.DB")
	}

	// Auto-migrate core models (safe no-op if tables exist)
	if err := gdb.AutoMigrate(&models.Book{}, &models.User{}); err != nil {
		log.Fatalf("Failed to migrate models: %v", err)
	}

	// Wire services with GORM-backed implementations
	userService = services.NewUserServiceImpl(gdb, ctx)
	authService = services.NewAuthService(gdb, ctx)
	bookService = services.NewBookService(gdb, ctx)

	// Create controllers and route controllers
	AuthController = controllers.NewAuthController(authService, userService)
	AuthRouteController = routes.NewAuthRouteController(AuthController)

	UserController = controllers.NewUserController(userService, bookService)
	UserRouteController = routes.NewRouteUserController(UserController)

	BookController := controllers.NewBookController(bookService, userService)
	BookRouteController = routes.NewBookRouteController(BookController)

	server = gin.New()
	server.Use(gin.Logger())   // Add Logger middleware explicitly
	server.Use(gin.Recovery()) // Add Recovery middleware explicitly
}

func main() {
	fmt.Println("Starting Fast Reading API...")
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load config", err)
	}

	fmt.Println("config loaded")

	defer models.RemoveDb(db)

	server.Static("/static", "./static")

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:8080", "http://localhost:3000"}
	corsConfig.AllowCredentials = true

	server.Use(cors.New(corsConfig))

	router := server.Group("/library")

	AuthRouteController.AuthRoute(router, userService)
	UserRouteController.UserRoute(router, userService)
	BookRouteController.BookRoute(router, bookService)

	log.Println("Registered routes:")
	for _, route := range server.Routes() {
		log.Printf("Method: %s, Path: %s", route.Method, route.Path)
	}

	log.Fatal(server.Run(":" + config.Port))
}
