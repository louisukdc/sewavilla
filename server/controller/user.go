package controller

import (
	"log"
	"time"

	"uas/database"
	"uas/model"

	"github.com/gofiber/fiber/v2"
)

// Get list of Users
func UserList(c *fiber.Ctx) error {
	context := fiber.Map{
		"statusText": "Ok",
		"msg":        "User List",
	}

	// Mengatur delay untuk simulasi response lebih lama
	time.Sleep(time.Millisecond * 1500)

	// Inisialisasi koneksi database
	db := database.DBConn

	var users []model.User

	// Mengambil semua user dari database
	result := db.Find(&users)

	// Jika terjadi error saat mengambil data
	if result.Error != nil {
		log.Println("Error fetching users:", result.Error)
		context["statusText"] = "Error"
		context["msg"] = "Something went wrong while fetching users."
		c.Status(500) // Internal Server Error
		return c.JSON(context)
	}

	// Mengembalikan daftar user
	context["data"] = users

	c.Status(200) // OK
	return c.JSON(context)
}

// Blog detail page
func UserDetail(c *fiber.Ctx) error {
	c.Status(400)
	context := fiber.Map{
		"statusText": "",
		"msg":        "",
	}

	id := c.Params("id")

	var record model.User

	database.DBConn.First(&record, id)

	if record.ID == 0 {
		log.Println("Record not Found.")
		context["msg"] = "Record not Found."

		c.Status(404)
		return c.JSON(context)
	}

	context["record"] = record
	context["statusText"] = "Ok"
	context["msg"] = "User Detail"
	c.Status(200)
	return c.JSON(context)
}

// Register a new user
func UserCreate(c *fiber.Ctx) error {
	context := fiber.Map{
		"statusText": "Ok",
		"msg":        "User Created",
	}

	user := new(model.User)

	// Parsing request body untuk mengambil data user
	if err := c.BodyParser(user); err != nil {
		log.Println("Error in parsing request.")
		context["statusText"] = ""
		context["msg"] = "Something went wrong."
		c.Status(400)
		return c.JSON(context)
	}

	// Simpan data user ke database
	result := database.DBConn.Create(user)

	// Jika ada error dalam penyimpanan data
	if result.Error != nil {
		log.Println("Error in saving user.")
		context["statusText"] = ""
		context["msg"] = "Something went wrong."
		c.Status(400)
		return c.JSON(context)
	}

	// Berikan respons sukses
	context["msg"] = "User registered successfully."
	context["data"] = user
	c.Status(201)
	return c.JSON(context)
}

// Update a User
func UserUpdate(c *fiber.Ctx) error {

	context := fiber.Map{
		"statusText": "Ok",
		"msg":        "Update User",
	}

	//http://localhost:8000/user/2

	id := c.Params("id")

	var user model.User

	// Menemukan user berdasarkan id
	database.DBConn.First(&user, id)

	if user.ID == 0 {
		log.Println("User not found.")

		context["statusText"] = ""
		context["msg"] = "User not found."
		c.Status(400)
		return c.JSON(context)
	}

	// Parsing data yang diterima dari request body
	if err := c.BodyParser(&user); err != nil {
		log.Println("Error in parsing request.")

		context["msg"] = "Something went wrong."
		c.Status(400)
		return c.JSON(context)
	}

	// Menyimpan perubahan pada database
	result := database.DBConn.Save(&user)

	if result.Error != nil {
		log.Println("Error in saving data.")

		context["msg"] = "Error in saving data."
		c.Status(400)
		return c.JSON(context)
	}

	context["msg"] = "User updated successfully."
	context["data"] = user

	c.Status(200)
	return c.JSON(context)
}

// Delete a User
func UserDelete(c *fiber.Ctx) error {
	context := fiber.Map{
		"statusText": "",
		"msg":        "",
	}

	// Mendapatkan ID user dari parameter URL
	id := c.Params("id")

	var user model.User

	// Mencari user berdasarkan ID
	database.DBConn.First(&user, id)

	// Jika user tidak ditemukan
	if user.ID == 0 {
		log.Println("User not Found.")
		context["msg"] = "User not found."
		c.Status(400)
		return c.JSON(context)
	}

	// Menghapus user dari database
	result := database.DBConn.Delete(&user)

	// Jika ada error dalam penghapusan
	if result.Error != nil {
		log.Println("Error in deleting user.", result.Error)
		context["msg"] = "Something went wrong."
		c.Status(400)
		return c.JSON(context)
	}

	// Berikan respons sukses
	context["statusText"] = "Ok"
	context["msg"] = "User deleted successfully."
	c.Status(200)
	return c.JSON(context)
}
