package main

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"

	logger "github.com/lokesh2201013/Logger"
	"github.com/lokesh2201013/database"
	"github.com/lokesh2201013/routes"

	_ "github.com/lokesh2201013/docs"             // docs is imported for Swagger setup
	fiberSwagger "github.com/swaggo/fiber-swagger" // fiber-swagger handler
	_ "github.com/swaggo/files"                   // swagger embed files (used internally)
)

// @title           Product API
// @version         1.0
// @description     API for managing products with JWT authentication
// @termsOfService  http://swagger.io/terms/

// @contact.name   Lokesh Choraria
// @contact.email  lokeshchoraria2@email.com

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization


func main() {
	app := fiber.New()

	logger.InitLogger()
	app.Use(logger.ZapLogger())

	database.ConnectDB()

	// Swagger route
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// CORS Middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// Rate Limiting
	app.Use(limiter.New(limiter.Config{
		Max:        100,
		Expiration: 60 * time.Second,
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Too many requests. Please try again later.",
			})
		},
	}))

	// Register routes
	routes.AuthRoutes(app)

	port := ":8080"
	log.Printf("Server is running on http://localhost%s\n", port)
	log.Fatal(app.Listen(port))
}
