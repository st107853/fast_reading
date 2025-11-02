package models

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"os"
)

type Database struct {
	*gorm.DB
}

var DB *gorm.DB

// Opening a database and save the reference to `Database` struct.
func OpenDbConnection() (*gorm.DB, error) {

	dsn := fmt.Sprintf("host=%s dbname=%s user=%s password=%s",
		os.Getenv("HOST"), os.Getenv("DBNAME"), os.Getenv("DBUSER"), os.Getenv("DBPASS"))

	// 2. Настраиваем GORM Logger (лучше, чем fmt.Println)
	// newLogger := logger.New(
	// 	log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
	// 	logger.Config{
	// 		SlowThreshold: time.Second, // Порог медленного SQL-запроса
	// 		LogLevel:      logger.Info, // Уровень логгирования
	// 		Colorful:      true,
	// 	},
	// )

	// 3. Открываем соединение
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: nil, // Используем настроенный логгер
	})

	if err != nil {
		// Ошибка подключения - возвращаем ее
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// 4. Запускаем AutoMigrate
	// (Используем новую модель User)

	log.Println("Database connection successful and migrated.")
	DB = db

	// 5. Возвращаем 'db', а не ошибку
	return db, nil
}

// Delete the database after running testing cases.
func RemoveDb(db *gorm.DB) error {
	sqlDB, err := db.DB() // Получаем *sql.DB из *gorm.DB
	if err != nil {
		fmt.Println("Error getting sql.DB from gorm.DB:", err)
		return err
	}

	// Закрываем соединение с базой
	err = sqlDB.Close()
	return err
}

// Using this function to get a connection, you can create your connection pool here.
func GetDB() *gorm.DB {
	return DB
}
