package controller

import (
	"log"

	"uas/database"
	"uas/model"

	"github.com/gofiber/fiber/v2"
)

// Get list of Reservations
func ReservationList(c *fiber.Ctx) error {
	var reservations []model.Reservation

	// Ambil semua reservasi beserta relasinya
	result := database.DBConn.Preload("User").Preload("Room").Preload("Blog").Find(&reservations)
	if result.Error != nil {
		log.Println("Error fetching reservations:", result.Error)
		return c.Status(500).JSON(fiber.Map{
			"statusText": "Error",
			"msg":        "Failed to fetch reservations.",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"statusText": "Ok",
		"msg":        "Reservation List",
		"data":       reservations,
	})
}

// Get Reservation by ID
func ReservationDetail(c *fiber.Ctx) error {
	id := c.Params("id")
	var reservation model.Reservation

	// Cari reservasi berdasarkan ID beserta relasinya
	result := database.DBConn.Preload("User").Preload("Room").Preload("Blog").First(&reservation, id)
	if result.Error != nil || reservation.ID == 0 {
		log.Println("Reservation not found.")
		return c.Status(404).JSON(fiber.Map{
			"statusText": "Error",
			"msg":        "Reservation not found.",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"statusText": "Ok",
		"msg":        "Reservation Detail",
		"data":       reservation,
	})
}

// Create a new Reservation
func ReservationCreate(c *fiber.Ctx) error {
	reservation := new(model.Reservation)

	// Parsing body JSON
	if err := c.BodyParser(reservation); err != nil {
		log.Println("Error parsing request body:", err)
		return c.Status(400).JSON(fiber.Map{
			"statusText": "Error",
			"msg":        "Invalid request body.",
		})
	}

	// Validasi User, Room, dan Blog
	if err := validateReservationRelationships(reservation); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"statusText": "Error",
			"msg":        err.Error(),
		})
	}

	// Simpan ke database
	result := database.DBConn.Create(reservation)
	if result.Error != nil {
		log.Println("Error creating reservation:", result.Error)
		return c.Status(500).JSON(fiber.Map{
			"statusText": "Error",
			"msg":        "Failed to create reservation.",
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"statusText": "Ok",
		"msg":        "Reservation created successfully.",
		"data":       reservation,
	})
}

// Update Reservation
func ReservationUpdate(c *fiber.Ctx) error {
	id := c.Params("id")
	var reservation model.Reservation

	// Cari reservasi berdasarkan ID
	result := database.DBConn.First(&reservation, id)
	if result.Error != nil || reservation.ID == 0 {
		log.Println("Reservation not found.")
		return c.Status(404).JSON(fiber.Map{
			"statusText": "Error",
			"msg":        "Reservation not found.",
		})
	}

	// Parsing body JSON
	if err := c.BodyParser(&reservation); err != nil {
		log.Println("Error parsing request body:", err)
		return c.Status(400).JSON(fiber.Map{
			"statusText": "Error",
			"msg":        "Invalid request body.",
		})
	}

	// Validasi User, Room, dan Blog
	if err := validateReservationRelationships(&reservation); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"statusText": "Error",
			"msg":        err.Error(),
		})
	}

	// Update ke database
	result = database.DBConn.Save(&reservation)
	if result.Error != nil {
		log.Println("Error updating reservation:", result.Error)
		return c.Status(500).JSON(fiber.Map{
			"statusText": "Error",
			"msg":        "Failed to update reservation.",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"statusText": "Ok",
		"msg":        "Reservation updated successfully.",
		"data":       reservation,
	})
}

// Delete Reservation
func ReservationDelete(c *fiber.Ctx) error {
	id := c.Params("id")
	var reservation model.Reservation

	// Cari reservasi berdasarkan ID
	result := database.DBConn.First(&reservation, id)
	if result.Error != nil || reservation.ID == 0 {
		log.Println("Reservation not found.")
		return c.Status(404).JSON(fiber.Map{
			"statusText": "Error",
			"msg":        "Reservation not found.",
		})
	}

	// Hapus reservasi
	result = database.DBConn.Delete(&reservation)
	if result.Error != nil {
		log.Println("Error deleting reservation:", result.Error)
		return c.Status(500).JSON(fiber.Map{
			"statusText": "Error",
			"msg":        "Failed to delete reservation.",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"statusText": "Ok",
		"msg":        "Reservation deleted successfully.",
	})
}

// Helper function to validate relationships
func validateReservationRelationships(reservation *model.Reservation) error {
	var user model.User
	var room model.Room
	var blog model.Blog

	// Validasi User
	if database.DBConn.First(&user, reservation.UserID).RowsAffected == 0 {
		return fiber.NewError(400, "Invalid user_id.")
	}

	// Validasi Room
	if database.DBConn.First(&room, reservation.RoomID).RowsAffected == 0 {
		return fiber.NewError(400, "Invalid room_id.")
	}

	// Validasi Blog
	if database.DBConn.First(&blog, reservation.BlogID).RowsAffected == 0 {
		return fiber.NewError(400, "Invalid blog_id.")
	}

	return nil
}
