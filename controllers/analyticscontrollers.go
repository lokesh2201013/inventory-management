package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/lokesh2201013/database"
	"github.com/lokesh2201013/models"
)

// GetProductByID godoc
// @Summary      Get a product by ID
// @Description  Retrieves a single product based on the provided UUID in query parameter
// @Tags         Products
// @Produce      json
// @Param        product_id  query     string  true  "Product UUID"  example("d290f1ee-6c54-4b01-90e6-d701748f0851")
// @Success      200  {object}  models.Product
// @Failure      400  {object}  map[string]string  "Missing or invalid product_id"
// @Failure      404  {object}  map[string]string  "Product not found"
// @Router       /products/get [get]
func GetProductByID(c *fiber.Ctx) error {
	productIDParam := c.Query("product_id")

	if productIDParam == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "product_id query parameter is required"})
	}

	productID, err := uuid.Parse(productIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid product_id format"})
	}

	var product models.Product
	err = database.DB.First(&product, "id = ?", productID).Error
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Product not found"})
	}

	return c.JSON(product)
}

// GetProductByQuantityExtremes godoc
// @Summary      Get product with extreme quantity
// @Description  Fetches either the product with the highest or lowest quantity based on query parameter
// @Tags         Products
// @Produce      json
// @Param        most   query  bool  false  "Set to true to get product with highest quantity"   example(true)
// @Param        least  query  bool  false  "Set to true to get product with lowest quantity"    example(false)
// @Success      200  {object}  models.Product
// @Failure      400  {object}  map[string]string  "Missing or conflicting query parameters"
// @Failure      500  {object}  map[string]string  "Internal server error"
// @Router       /products/extreme [get]
func GetProductByQuantityExtremes(c *fiber.Ctx) error {
	most := c.QueryBool("most")
	least := c.QueryBool("least")

	var product models.Product
	var err error

	switch {
	case most:
		err = database.DB.Order("quantity DESC").First(&product).Error
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get product with highest quantity"})
		}
		return c.JSON(product)

	case least:
		err = database.DB.Order("quantity ASC").First(&product).Error
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get product with lowest quantity"})
		}
		return c.JSON(product)

	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Provide either 'most=true' or 'least=true' in query",
		})
	}
}



