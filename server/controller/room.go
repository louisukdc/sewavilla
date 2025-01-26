package controller

import (
	"log"
	"strconv"
	"time"

	"uas/database"
	"uas/model"

	"github.com/gofiber/fiber/v2"
)

// Get list of Users
func RoomList(c *fiber.Ctx) error {
	context := fiber.Map{
		"statusText": "Ok",
		"msg":        "Room List",
	}

	// Mengatur delay untuk simulasi response lebih lama
	time.Sleep(time.Millisecond * 1500)

	// Inisialisasi koneksi database
	db := database.DBConn

	var room []model.Room

	// Mengambil semua user dari database
	result := db.Find(&room)

	// Jika terjadi error saat mengambil data
	if result.Error != nil {
		log.Println("Error fetching room:", result.Error)
		context["statusText"] = "Error"
		context["msg"] = "Something went wrong while fetching room."
		c.Status(500) // Internal Server Error
		return c.JSON(context)
	}

	// Mengembalikan daftar user
	context["data"] = room

	c.Status(200) // OK
	return c.JSON(context)
}

// Blog detail page
func RoomDetail(c *fiber.Ctx) error {
	c.Status(400)
	context := fiber.Map{
		"statusText": "",
		"msg":        "",
	}

	id := c.Params("id")

	var record model.Room

	database.DBConn.First(&record, id)

	if record.ID == 0 {
		log.Println("Record not Found.")
		context["msg"] = "Record not Found."

		c.Status(404)
		return c.JSON(context)
	}

	context["record"] = record
	context["statusText"] = "Ok"
	context["msg"] = "Room Detail"
	c.Status(200)
	return c.JSON(context)
}

// Register a new user
func RoomCreate(c *fiber.Ctx) error {
	context := fiber.Map{
		"statusText": "Ok",
		"msg":        "Room Created",
	}

	room := new(model.Room)

	// Parsing request body untuk mengambil data user
	if err := c.BodyParser(room); err != nil {
		log.Println("Error in parsing request.")
		context["statusText"] = ""
		context["msg"] = "Something went wrong."
		c.Status(400)
		return c.JSON(context)
	}

	// Simpan data user ke database
	result := database.DBConn.Create(room)

	// Jika ada error dalam penyimpanan data
	if result.Error != nil {
		log.Println("Error in saving room.")
		context["statusText"] = ""
		context["msg"] = "Something went wrong."
		c.Status(400)
		return c.JSON(context)
	}

	// Berikan respons sukses
	context["msg"] = "Room registered successfully."
	context["data"] = room
	c.Status(201)
	return c.JSON(context)
}

// Update a User
func RoomUpdate(c *fiber.Ctx) error {

	context := fiber.Map{
		"statusText": "Ok",
		"msg":        "Update Room",
	}

	//http://localhost:8000/user/2

	id := c.Params("id")

	var room model.Room

	// Menemukan user berdasarkan id
	database.DBConn.First(&room, id)

	if room.ID == 0 {
		log.Println("Room not found.")

		context["statusText"] = ""
		context["msg"] = "Room not found."
		c.Status(400)
		return c.JSON(context)
	}

	// Parsing data yang diterima dari request body
	if err := c.BodyParser(&room); err != nil {
		log.Println("Error in parsing request.")

		context["msg"] = "Something went wrong."
		c.Status(400)
		return c.JSON(context)
	}

	// Menyimpan perubahan pada database
	result := database.DBConn.Save(&room)

	if result.Error != nil {
		log.Println("Error in saving data.")

		context["msg"] = "Error in saving data."
		c.Status(400)
		return c.JSON(context)
	}

	context["msg"] = "Room updated successfully."
	context["data"] = room

	c.Status(200)
	return c.JSON(context)
}

// Delete a User
func RoomDelete(c *fiber.Ctx) error {
	context := fiber.Map{
		"statusText": "",
		"msg":        "",
	}

	// Mendapatkan ID user dari parameter URL
	id := c.Params("id")

	var room model.Room

	// Mencari user berdasarkan ID
	database.DBConn.First(&room, id)

	// Jika user tidak ditemukan
	if room.ID == 0 {
		log.Println("Room not Found.")
		context["msg"] = "Room not found."
		c.Status(400)
		return c.JSON(context)
	}

	// Menghapus user dari database
	result := database.DBConn.Delete(&room)

	// Jika ada error dalam penghapusan
	if result.Error != nil {
		log.Println("Error in deleting room.", result.Error)
		context["msg"] = "Something went wrong."
		c.Status(400)
		return c.JSON(context)
	}

	// Berikan respons sukses
	context["statusText"] = "Ok"
	context["msg"] = "Room deleted successfully."
	c.Status(200)
	return c.JSON(context)
}

// UpdateRoomStatus untuk memperbarui status kamar
func UpdateRoomStatus(roomID int, status string) {
	var room model.Room
	if err := database.DBConn.First(&room, roomID).Error; err != nil {
		log.Println("Room not found:", err)
		return
	}

	room.Status = status
	if err := database.DBConn.Save(&room).Error; err != nil {
		log.Println("Error updating room status:", err)
		return
	}

	// Kirim pembaruan status kamar ke semua klien melalui WebSocket
	message := `{"room_id": ` + strconv.Itoa(roomID) + `, "status": "` + status + `"}`
	broadcast <- message
}

// BookRoom untuk menangani pemesanan kamar dan memperbarui status
func BookRoom(c *fiber.Ctx) error {
	roomIDStr := c.Params("id") // Mendapatkan id sebagai string

	// Konversi id kamar ke tipe int
	roomID, err := strconv.Atoi(roomIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid room ID",
		})
	}

	// Proses pemesanan kamar
	UpdateRoomStatus(roomID, "booked")

	context := fiber.Map{
		"status":  "success",
		"message": "Room booked successfully.",
	}
	return c.JSON(context)
}
