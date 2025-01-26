package database

import (
	"log"
	"os"

	"uas/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DBConn *gorm.DB

func ConnectDB() {
	// Load database credentials from environment variables
	host := os.Getenv("db_host")
	user := os.Getenv("db_user")
	password := os.Getenv("db_password")
	dbname := os.Getenv("db_name")

	// Construct DSN
	var dsn string
	if password == "" {
		dsn = user + "@tcp(" + host + ":3306)/" + dbname + "?charset=utf8mb4&parseTime=True&loc=Local"
	} else {
		dsn = user + ":" + password + "@tcp(" + host + ":3306)/" + dbname + "?charset=utf8mb4&parseTime=True&loc=Local"
	}

	// Open database connection
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connection successful.")

	// Automatically migrate the database schema
	// db.AutoMigrate(&model.Blog{})

	// Automigrate models
	db.AutoMigrate(&model.Blog{}, &model.User{}, &model.Room{}, &model.Reservation{})

	DBConn = db
}
