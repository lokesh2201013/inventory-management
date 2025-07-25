package utils

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var secretKey = []byte("itsmysecrectkey9310") 

func GenerateJWT(userID uuid.UUID, name string) (string, error) {
	claims := jwt.MapClaims{
		"exp":  time.Now().Add(24 * time.Hour).Unix(),
		"iss":  "invertory",
		"userID":  userID,
		"email": name,
	}
	fmt.Println("+++++++++++++++++++++++++++++++++>USERID",userID)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		log.Println("Error generating JWT:", err)
		return "", err
	}

	return signedToken, nil
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CheckPassword(hashedPassword, plainPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
}
