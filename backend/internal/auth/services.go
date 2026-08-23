package auth

import (
	"context"
	"errors"
	"os"
	"time"

	"factorypulse/backend/internal/database"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(input RegisterInput) (int, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	var userID int
	query := `INSERT INTO users (name, email, password_hash, role) 
	          VALUES ($1, $2, $3, 'technician') RETURNING id`

	err = database.Pool.QueryRow(context.Background(), query, input.Name, input.Email, string(hashedPassword)).Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func AuthenticateUser(input LoginInput) (string, error) {
	var user User

	query := `SELECT id, name, email, password_hash, role FROM users WHERE email = $1`
	err := database.Pool.QueryRow(context.Background(), query, input.Email).Scan(
		&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.Role,
	)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	// Password check karo — diya gaya password aur DB mein saved hash match karta hai?
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password))
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	// JWT token banao
	token, err := generateToken(user)
	if err != nil {
		return "", err
	}

	return token, nil
}

func generateToken(user User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(), // 24 ghante mein expire hoga
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}
