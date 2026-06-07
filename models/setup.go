package models

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Database struct {
	*gorm.DB
}

var DB *gorm.DB

// Opening a database and save the reference to `Database` struct.
func OpenDbConnectionWithConfig(host, dbname, user, password string) (*gorm.DB, error) {

	dsn := fmt.Sprintf("host=%s dbname=%s user=%s password=%s",
		host, dbname, user, password)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: nil,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Database connection successful and migrated.")
	DB = db

	return db, nil
}

// Delete the database after running testing cases.
func RemoveDb(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		fmt.Println("Error getting sql.DB from gorm.DB:", err)
		return err
	}

	err = sqlDB.Close()
	return err
}

// Using this function to get a connection, you can create your connection pool here.
func GetDB() *gorm.DB {
	return DB
}
