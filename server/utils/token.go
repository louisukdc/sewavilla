package utils

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

// Secret key used for signing the JWT token
var jwtKey = []byte("your_secret_key")

// Claims struct represents the structure of JWT claims
type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

// ValidateToken verifies the JWT token and returns the userID if valid
func ValidateToken(tokenString string) (uint, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		// Pastikan subject adalah userID yang valid
		userID, err := strconv.ParseUint(claims.Subject, 10, 64)
		if err != nil {
			return 0, errors.New("invalid user ID in token")
		}
		return uint(userID), nil
	}
	return 0, errors.New("invalid token")
}

// Function to parse and validate JWT
func ParseToken(tokenString string) (*jwt.Token, error) {
	// Get the secret key from environment variables
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		return nil, errors.New("JWT secret key is not set")
	}

	// Parse the token using the secret key
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the token's signing method (HMAC)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secretKey), nil
	})

	return token, err
}

// GenerateToken creates a JWT token for a given user
func GenerateToken(userID uint) (string, error) {
	claims := &jwt.RegisteredClaims{
		Issuer:    "uas-app",
		Subject:   string(userID),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(48 * time.Hour)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

// HashPassword hashes the password using bcrypt
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// CheckPasswordHash compares a hashed password with its plain text equivalent
func CheckPasswordHash(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
