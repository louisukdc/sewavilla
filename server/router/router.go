package router

import (
	"uas/controller"
	"uas/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

// Setup routing information
func SetupRoutes(app *fiber.App) {

	// Public routes (without JWT protection)
	app.Post("/login", controller.Login)
	app.Post("/register", controller.Register)

	// Protected routes (with JWT protection)
	protected := app.Group("/protected")
	protected.Use(middleware.Protected()) // Protect all routes in this group
	// protected.Get("/profile", controller.Profile)
	app.Get("/profile", middleware.Protected(), controller.GetProfile)
	protected.Post("/update", controller.UpdateProfile)
	app.Put("/update-password", middleware.Protected(), controller.UpdatePassword)

	// WebSocket endpoint
	app.Get("/ws", websocket.New(controller.WebSocketHandler))

	// Blog routes
	app.Get("/blog", controller.BlogList)          // List all blogs
	app.Get("/blog/:id", controller.BlogDetail)    // Get a specific blog by ID
	app.Post("/blog", controller.BlogCreate)       // Create a new blog
	app.Put("/blog/:id", controller.BlogUpdate)    // Update an existing blog by ID
	app.Delete("/blog/:id", controller.BlogDelete) // Delete a blog by ID

	// User routes
	app.Get("/user", controller.UserList)          // List all users
	app.Get("/user/:id", controller.UserDetail)    // Get a specific user by ID
	app.Post("/user", controller.UserCreate)       // Create a new user
	app.Put("/user/:id", controller.UserUpdate)    // Update an existing user by ID
	app.Delete("/user/:id", controller.UserDelete) // Delete a user by ID

	// Room routes
	app.Get("/room", controller.RoomList)          // List all users
	app.Get("/room/:id", controller.RoomDetail)    // Get a specific user by ID
	app.Post("/room", controller.RoomCreate)       // Create a new user
	app.Put("/room/:id", controller.RoomUpdate)    // Update an existing user by ID
	app.Delete("/room/:id", controller.RoomDelete) // Delete a user by ID

	// Reservation routes
	app.Get("/reservation", controller.ReservationList)          // List all users
	app.Get("/reservation/:id", controller.ReservationDetail)    // Get a specific user by ID
	app.Post("/reservation", controller.ReservationCreate)       // Create a new user
	app.Put("/reservation/:id", controller.ReservationUpdate)    // Update an existing user by ID
	app.Delete("/reservation/:id", controller.ReservationDelete) // Delete a user by ID

	// In router.go
	app.Use(func(c *fiber.Ctx) error {
		// Call your custom error handler here
		if err := middleware.ErrorHandler(c, nil); err != nil {
			return err
		}
		return c.Next()
	})

}
