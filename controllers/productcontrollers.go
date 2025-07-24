package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/lokesh2201013/Logger"
	"github.com/lokesh2201013/database"
	"github.com/lokesh2201013/models"
	"go.uber.org/zap"
)

// ProductInsert godoc
// @Summary      Insert a new product
// @Description  Create a new product linked to the logged-in user
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        product  body      models.Product  true  "Product Info"
// @Success      201  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Router       /products/ [post]

func ProductInsert(c *fiber.Ctx) error {
	const file = "ProductController"
	var product models.Product

	if err := c.BodyParser(&product); err != nil {
		logger.Log.Error("Package controllers File "+file,zap.String("Function", "ProductInsert"),zap.String("Message", "Failed to parse request body"),zap.Error(err),
		)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	if product.Name == "" || product.SKU == "" || product.Quantity < 0 || product.Price < 0 {
		logger.Log.Error("Package controllers File "+file,zap.String("Function", "ProductInsert"),zap.String("Message", "Invalid product fields"),zap.Any("product", product),
		)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid product fields"})
	}

	userID, ok := c.Locals("userID").(uuid.UUID)
	if !ok {
		logger.Log.Error("Package controllers File "+file,zap.String("Function", "ProductInsert"),zap.String("Message", "Invalid user ID in context"),)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user ID"})
	}
	product.UserID = userID

	if err := database.DB.Create(&product).Error; err != nil {
		logger.Log.Error("Package controllers File "+file,zap.String("Function", "ProductInsert"),zap.String("Message", "Database error while creating product"),zap.Error(err),)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error saving product"})
	}

	logger.Log.Info("Package controllers File "+file,zap.String("Function", "ProductInsert"),zap.String("Message", "Product inserted successfully"),zap.String("product_id", product.ID.String()),
	)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":    "Product inserted successfully",
		"product_id": product.ID,
	})
}

// UpdateQuantity godoc
// @Summary      Update product quantity
// @Description  Update the quantity of an existing product by ID
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id      path      string              true  "Product ID"
// @Param        input   body      map[string]int      true  "Quantity Input"
// @Success      200     {object}  models.Product
// @Failure      400     {object}  map[string]string
// @Failure      404     {object}  map[string]string
// @Router       /products/{id}/quantity [put]

func UpdateQuantity(c *fiber.Ctx) error {
	const file = "ProductController"
	productID := c.Params("id")

	var input struct {
		Quantity int `json:"quantity"`
	}

	if err := c.BodyParser(&input); err != nil {
		logger.Log.Error("Package controllers File "+file,zap.String("Function", "UpdateQuantity"),zap.String("Message", "Failed to parse input"),zap.Error(err),)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	var product models.Product
	if err := database.DB.First(&product, "id = ?", productID).Error; err != nil {
		logger.Log.Error("Package controllers File "+file,zap.String("Function", "UpdateQuantity"),zap.String("Message", "Product not found"),zap.String("product_id", productID),)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Product not found"})
	}

	product.Quantity = input.Quantity

	if err := database.DB.Save(&product).Error; err != nil {
		logger.Log.Error("Package controllers File "+file,
			zap.String("Function", "UpdateQuantity"),
			zap.String("Message", "Failed to update product"),
			zap.Error(err),
		)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update product"})
	}

	logger.Log.Info("Package controllers File "+file,zap.String("Function", "UpdateQuantity"),zap.String("Message", "Product quantity updated"),zap.String("product_id", productID),zap.Int("new_quantity", input.Quantity),
	)

	return c.Status(fiber.StatusOK).JSON(product)
}

// GetAllUserProduct godoc
// @Summary      Get all user products
// @Description  Get paginated list of products created by the user
// @Tags         Products
// @Produce      json
// @Security     BearerAuth
// @Param        pagenum  query     int  false  "Page Number"
// @Param        limit    query     int  false  "Limit"
// @Success      200      {array}   models.Product
// @Failure      401      {object}  map[string]string
// @Router       /products/all [get]
func GetAllUserProduct(c *fiber.Ctx) error {
	const file = "ProductController"
	userID, ok := c.Locals("userID").(uuid.UUID)
	if !ok {
		logger.Log.Error("Package controllers File "+file,zap.String("Function", "GetAllUserProduct"),zap.String("Message", "Invalid user ID in context"),
		)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	pageNumber := c.QueryInt("pagenum", 1)
	if pageNumber <= 0 {
		pageNumber = 1
	}
	limit := c.QueryInt("limit", 10)
	offset := (pageNumber - 1) * limit

	var products []models.Product
	if err := database.DB.Where("user_id = ?", userID).
		Limit(limit).Offset(offset).
		Find(&products).Error; err != nil {
		logger.Log.Error("Package controllers File "+file,zap.String("Function", "GetAllUserProduct"),zap.String("Message", "Error retrieving products"),zap.String("user_id", userID.String()),zap.Error(err),)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error retrieving products"})
	}

	logger.Log.Info("Package controllers File "+file,zap.String("Function", "GetAllUserProduct"),zap.String("Message", "Products retrieved successfully"),zap.String("user_id", userID.String()),zap.Int("count", len(products)),
	)
	return c.JSON(products)
}
