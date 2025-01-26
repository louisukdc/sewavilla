package controller

import (
	"fmt"
	"log"
	"strings"
	"time"
	"uas/database"
	"uas/model"
	"uas/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Secret key for JWT (you can store this in an environment variable for better security)
var jwtSecret = []byte("your_jwt_secret_key")

// Simulasi kesalahan login (misalnya user tidak ditemukan atau password salah)
func someLoginLogic() error {
	return fmt.Errorf("Login failed due to invalid credentials")
}

func GetProfile(c *fiber.Ctx) error {
	// Ambil data user dari context yang sudah di-set di middleware
	user := c.Locals("user").(model.User)

	// Kirim data user sebagai response
	return c.JSON(fiber.Map{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
	})
}

// Profile handles fetching the user profile for authenticated users
func Profile(c *fiber.Ctx) error {
	// Extract token from Authorization header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Authorization header missing"})
	}

	// Split the token from the "Bearer " prefix
	token := strings.Replace(authHeader, "Bearer ", "", 1)

	// Validate the token
	userID, err := utils.ValidateToken(token)
	if err != nil {
		log.Println("Error validating token:", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
	}

	// Fetch user details from the database using userID
	var user model.User
	if err := database.DBConn.First(&user, userID).Error; err != nil {
		log.Println("Error fetching user:", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	// Respond with user profile data (excluding password)
	return c.JSON(fiber.Map{
		"id":    user.ID,
		"email": user.Email,
	})
}

// UpdateProfile handles updating the user profile
func UpdateProfile(c *fiber.Ctx) error {
	// Extract token from Authorization header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Authorization header missing"})
	}

	// Split the token from the "Bearer " prefix
	token := strings.Replace(authHeader, "Bearer ", "", 1)

	// Validate the token
	userID, err := utils.ValidateToken(token)
	if err != nil {
		log.Println("Error validating token:", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
	}

	// Fetch the current user details from the database
	var user model.User
	if err := database.DBConn.First(&user, userID).Error; err != nil {
		log.Println("Error fetching user:", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	// Bind the request body to the user struct for update
	if err := c.BodyParser(&user); err != nil {
		log.Println("Error parsing request body:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Update the user's profile in the database
	if err := database.DBConn.Save(&user).Error; err != nil {
		log.Println("Error updating user:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update user"})
	}

	// Return the updated user profile data
	return c.JSON(fiber.Map{
		"id":    user.ID,
		"email": user.Email,
	})
}

// Register handles user registration and generates a JWT token
func Register(c *fiber.Ctx) error {
	var newUser model.User
	if err := c.BodyParser(&newUser); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Hash the password before saving to database
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error hashing password:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not hash password"})
	}
	newUser.Password = string(hashedPassword)

	// Save the new user to the database
	if err := database.DBConn.Create(&newUser).Error; err != nil {
		log.Println("Error saving user to the database:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not save user"})
	}

	// Generate a JWT token after successful registration
	token, err := utils.GenerateToken(newUser.ID)
	if err != nil {
		log.Println("Error generating token:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate token"})
	}

	// Respond with the token
	return c.JSON(fiber.Map{"token": token})
}

// Login function that generates JWT token
func Login(c *fiber.Ctx) error {
	// Parse the login request body (example: username and password)
	var reqBody struct {
		// Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&reqBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
		})
	}

	// Declare a user struct
	var user struct {
		ID       int    `json:"id"`
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password"` // This should be the hashed password
	}

	// Explicitly specify the table name if necessary
	if err := database.DBConn.Table("users").Where("username = ?", reqBody.Username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// If no user found, return Unauthorized status
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid credentials",
			})
		}
		// If there is another error, return Internal Server Error
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Internal Server Error",
			"details": err.Error(),
		})
	}

	// Compare the provided password with the stored hashed password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqBody.Password))
	if err != nil {
		// Password doesn't match, return Unauthorized
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid credentials",
		})
	}

	// Create the JWT token
	claims := jwt.MapClaims{
		"user_id": user.ID, // You can store user ID or other user info here
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Could not generate token",
			"details": err.Error(),
		})
	}

	// If successful, return the JWT token
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Login successful",
		"token":   tokenString,
	})
}

func UpdatePassword(c *fiber.Ctx) error {
	// Parse request body
	type UpdatePasswordRequest struct {
		OldPassword string `json:"old_password"` // Opsional: Untuk validasi password lama
		NewPassword string `json:"new_password"`
	}

	var request UpdatePasswordRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
		})
	}

	// Validate input
	if request.NewPassword == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "New password is required",
		})
	}

	// Get user from context (set by middleware.Protected())
	user := c.Locals("user").(model.User)

	// Optional: Verify old password
	if request.OldPassword != "" {
		if !utils.CheckPasswordHash(request.OldPassword, user.Password) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "Old password is incorrect",
			})
		}
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(request.NewPassword)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Could not hash password",
		})
	}

	// Update password in database
	user.Password = hashedPassword
	if err := database.DBConn.Save(&user).Error; err != nil {
		log.Printf("Error updating password: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Could not update password",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Password updated successfully",
	})
}
