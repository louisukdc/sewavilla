package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"uas/controller"
	"uas/database"
	"uas/router"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	database.ConnectDB()
}

func main() {
	sqlDB, err := database.DBConn.DB()
	if err != nil {
		log.Fatal("Error initializing SQL database connection:", err)
	}
	defer func() {
		if err := sqlDB.Close(); err != nil {
			log.Println("Error closing database connection:", err)
		}
	}()

	app := fiber.New()

	app.Static("/static", "./static")
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))
	app.Use(logger.New())

	router.SetupRoutes(app)
	go controller.HandleMessages()

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8001"
	}

	// Graceful shutdown on SIGINT or SIGTERM
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		log.Println("Shutting down server gracefully...")
		if err := app.Shutdown(); err != nil {
			log.Fatal("Error shutting down server:", err)
		}
	}()

	log.Printf("Server is running on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}

	fmt.Println("Program selesai setelah delay")
	time.Sleep(2 * time.Second)
}
