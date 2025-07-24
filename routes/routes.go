package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lokesh2201013/controllers"
	"github.com/lokesh2201013/utils"
	//"github.com/lokesh2201013/database"
)


func AuthRoutes(app *fiber.App) {
	// Public routes
	app.Post("/login", controllers.Login)
	app.Post("/register", controllers.Register)

	// Protected routes (grouped)
	protected := app.Group("/products", utils.AuthMiddleware())
	protected.Post("/", controllers.ProductInsert)
	protected.Put("/:id/quantity",controllers.UpdateQuantity)
	protected.Get("/all",controllers.GetAllUserProduct)
}
