package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/lokesh2201013/Logger"
	"github.com/lokesh2201013/database"
	"github.com/lokesh2201013/models"
	"go.uber.org/zap"
	"fmt"
)

// ProductInsert godoc
// @Summary      Add a new product
// @Description  Adds a new product to the database
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        product body models.Product true "Product Info"
// @Success      201 {object} map[string]interface{} "Created Successfully"
// @Failure      400 {object} map[string]string "Invalid input or fields"
// @Failure      401 {object} map[string]string "Unauthorized"
// @Failure      500 {object} map[string]string "Internal server error"
// @Security     BearerAuth
// @Router       /products [post]
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
	userIDStr, ok := c.Locals("userID").(string)
if !ok {
	logger.Log.Error("Package controllers File "+file,
		zap.String("Function", "ProductInsert"),
		zap.String("Message", "userID not found or not a string"),
	)
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user ID"})
}

userID, err := uuid.Parse(userIDStr)
if err != nil {
	logger.Log.Error("Package controllers File "+file,
		zap.String("Function", "ProductInsert"),
		zap.String("Message", "Invalid UUID format in userID"),
		zap.Error(err),
	)
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user ID format"})
}

  
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
// @Param        id     path      string                     true  "Product ID (UUID)"
// @Param        input  body      models.QuantityUpdateRequest true  "Quantity Update Payload"
// @Success      200    {object}  models.Product
// @Failure      400    {object}  map[string]string "Invalid input"
// @Failure      404    {object}  map[string]string "Product not found"
// @Failure      500    {object}  map[string]string "Internal server error"
// @Security     BearerAuth
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
 if input.Quantity<0 {
    	logger.Log.Error("Package controllers File "+file,zap.String("Function", "UpdateQuantity"),zap.String("Message", "Quantity is lees < 0"),zap.String("product_id", productID),)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Quantity invalid"})
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
// @Description  Get paginated list of products created by the authenticated user
// @Tags         Products
// @Produce      json
// @Param        pagenum  query     int  false  "Page number (default: 1)"
// @Param        limit    query     int  false  "Items per page (default: 10)"
// @Success      200      {array}   models.Product
// @Failure      401      {object}  map[string]string "Unauthorized"
// @Failure      500      {object}  map[string]string "Internal server error"
// @Security     BearerAuth
// @Router       /products [get]
func GetAllUserProduct(c *fiber.Ctx) error {
	const file = "ProductController"
	userIDStr, ok := c.Locals("userID").(string)

if !ok {
	logger.Log.Error("Package controllers File "+file,
		zap.String("Function", "ProductInsert"),
		zap.String("Message", "userID not found or not a string"),
	)
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user ID"})
}

userID, err := uuid.Parse(userIDStr)

if err != nil {
	logger.Log.Error("Package controllers File "+file,
		zap.String("Function", "ProductInsert"),
		zap.String("Message", "Invalid UUID format in userID"),
		zap.Error(err),
	)
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user ID format"})
}

userido:=c.Locals("userID")
	//userID, ok := c.Locals("userID").(uuid.UUID)
	   fmt.Println("---------------> USERIDO",userido)
	      fmt.Println("---------------> USERID",userID)

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
	
	fmt.Println(products)
	return c.JSON(products)
}
