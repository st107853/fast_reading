package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/st107853/fast_reading/config"
	"github.com/st107853/fast_reading/controllers"
	"github.com/st107853/fast_reading/logger"
	"github.com/st107853/fast_reading/models"
	"github.com/st107853/fast_reading/routes"
	"github.com/st107853/fast_reading/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	server *gin.Engine
	ctx    context.Context
	conf   config.Config

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
	var err error

	fmt.Println("Starting Fast Reading API...")
	conf, err = config.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load config", err)
	}

	// Initialize GORM via models helper (returns *gorm.DB)
	gdb, err := models.OpenDbConnectionWithConfig(conf.Host, conf.DBname, conf.DBuser, conf.DBpassword)
	if err != nil {
		logger.Log.Error("database connection failed", slog.Any("error", err))
		os.Exit(1)
	}
	if gdb == nil {
		logger.Log.Error("models.OpenDbConnection returned nil *gorm.DB")
	}

	// Auto-migrate core models (safe no-op if tables exist)
	if err := gdb.AutoMigrate(&models.Book{}, &models.User{}, &models.ReadingProgress{}); err != nil {
		logger.Log.Error("failed to migrate models", slog.Any("error", err))
		os.Exit(1)
	}

	// Wire services with GORM-backed implementations
	userService = services.NewUserServiceImpl(gdb, ctx)
	authService = services.NewAuthService(gdb, ctx)
	bookService = services.NewBookService(gdb, ctx)

	// Create controllers and route controllers
	AuthController = controllers.NewAuthController(authService, userService, conf)
	AuthRouteController = routes.NewAuthRouteController(AuthController)

	UserController = controllers.NewUserController(userService, bookService)
	UserRouteController = routes.NewRouteUserController(UserController)

	BookController := controllers.NewBookController(bookService, userService)
	BookRouteController = routes.NewBookRouteController(BookController)

	var level slog.Level
	if err := level.UnmarshalText([]byte(conf.LogLevel)); err != nil {
		level = slog.LevelInfo // safe default
	}

	logger.Init(logger.Config{
		Level:  level,
		Format: conf.LogFormat,
	})

	logger.Log.Info("starting Fast Reading API",
		slog.String("log_level", conf.LogLevel),
		slog.String("log_format", conf.LogFormat),
	)

	server = gin.New()
	server.Use(gin.Logger())   // Add Logger middleware explicitly
	server.Use(gin.Recovery()) // Add Recovery middleware explicitly
}

func main() {

	defer models.RemoveDb(db)

	server.Static("/static", "./static")
	server.Static("/covers", "./covers")

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{conf.CORS1, conf.CORS2}
	corsConfig.AllowCredentials = true

	server.Use(cors.New(corsConfig))

	router := server.Group("/library")

	AuthRouteController.AuthRoute(router, userService, conf)
	UserRouteController.UserRoute(router, userService, conf)
	BookRouteController.BookRoute(router, bookService, userService, conf)

	log.Println("Registered routes:")
	for _, route := range server.Routes() {
		logger.Log.Info("Route registered", slog.String("method", route.Method), slog.String("path", route.Path))
	}

	log.Fatal(server.Run(":" + conf.Port))
}
