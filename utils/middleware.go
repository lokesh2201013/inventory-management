package utils

import (
	"fmt"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenHeader := c.Get("Authorization")
		if tokenHeader == "" || !strings.HasPrefix(tokenHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}

		tokenString := strings.TrimPrefix(tokenHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid signing method")
			}
			return secretKey, nil
		})

		if err != nil || !token.Valid {
			log.Println("Token validation error:", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		claims := token.Claims.(jwt.MapClaims)
		c.Locals("userID", claims["userID"])
		c.Locals("email", claims["name"])
		fmt.Println("claims-=-=-=-=-=-=-=-=-=-=-=-=-=->",claims["userID"],claims["name"])

		return c.Next()
	}
}
