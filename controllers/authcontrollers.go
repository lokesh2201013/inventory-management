package controllers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/lokesh2201013/Logger"
	"github.com/lokesh2201013/database"
	"github.com/lokesh2201013/models"
	"github.com/lokesh2201013/utils"
	"go.uber.org/zap"
)

// Register godoc
// @Summary      Register a new user
// @Description  Creates a new user in the system with a unique username and email. The password is securely hashed before storage.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        user  body      models.User  true  "User registration payload"
// @Success      201   {object}  models.User  "User successfully registered"
// @Failure      400   {object}  map[string]string  "Bad Request – Invalid JSON or missing fields"
// @Failure      409   {object}  map[string]string  "Conflict – Username or email already exists"
// @Failure      500   {object}  map[string]string  "Internal server error"
// @Router       /register [post]
func Register(c *fiber.Ctx)error{

	var user models.User
	if err:= c.BodyParser(&user); err!=nil{
	  logger.Log.Error(" Package controllers File Authcontroller", zap.Error(err), zap.String("Message", "Failed to decrypt request body"))
      return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error":"Invalid input"})
	}
    
	var existingUser models.User

	if err:=database.DB.Where("username=?",user.Username).First(&existingUser).Error;err==nil{
		logger.Log.Error(" Package controllers File Authcontroller", zap.Error(err), zap.String("Message", "Similar Name"))
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Email already in use"})
	}

	/*if err:=database.DB.Where("email=?",user.Email).First(&existingUser).Error;err==nil{
			  logger.Log.Error(" Package controllers File Authcontroller", zap.Error(err), zap.String("Message", "Similar Emails"))
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Email already in use"})
	}*/

	hashpassword, err:=utils.HashPassword(user.Password)

	if err != nil {
			  logger.Log.Error(" Package controllers File Authcontroller", zap.Error(err), zap.String("Message", "Failed to Hash request body password"))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}
	user.Password = hashpassword

	if err:=database.DB.Create(&user).Error;err!=nil{
			  logger.Log.Error(" Package controllers File Authcontroller", zap.Error(err), zap.String("Message", "Failed to Create database input"))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}

// Login godoc
// @Summary      User login
// @Description  Authenticates a registered user using username and password. Returns a JWT token upon success.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        login  body      models.LoginRequest  true  "Login credentials"
// @Success      200    {object}  map[string]string     "Authentication successful – JWT access token returned"
// @Failure      400    {object}  map[string]string     "Bad Request – Invalid JSON or missing fields"
// @Failure      401    {object}  map[string]string     "Unauthorized – Incorrect password"
// @Failure      404    {object}  map[string]string     "Not Found – User does not exist"
// @Failure      500    {object}  map[string]string     "Internal server error"
// @Router       /login [post]
func Login(c *fiber.Ctx)error{
	var loginData struct {
		Name    string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&loginData); err != nil {
		logger.Log.Error(" Package controllers File Authcontroller", zap.Error(err), zap.String("Message", "Failed to parse Request body"))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error":"invald input"})
	}

	var user models.User

	if err:= database.DB.Where("username=?",loginData.Name).First(&user).Error;err!=nil{
		logger.Log.Error(" Package controllers File Authcontroller", zap.Error(err), zap.String("Message", "Failed to Create database input"))
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error":"User not found"})
	}

	if err := utils.CheckPassword(user.Password, loginData.Password); err != nil {
		logger.Log.Error(" Package controllers File Authcontroller", zap.Error(err), zap.String("Message", "Failed to Create database input"))
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	token, err := utils.GenerateJWT(user.UserID, user.Username)
	fmt.Println("<<<<<<<<<<<<<<<<<<<<<<>TOKEN",token)

	if err != nil {
		logger.Log.Error(" Package controllers File Authcontroller", zap.Error(err), zap.String("Message", "Failed to Create database input"))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}
    
	logger.Log.Info(" Package controllers File Authcontroller", zap.Error(err), zap.String("Message", "token created"))
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"access_token":   token,
	}) 
}