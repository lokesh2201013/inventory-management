package controllers

import (
	"go.uber.org/zap"
	"github.com/gofiber/fiber/v2"
     "github.com/lokesh2201013/Logger"
	"github.com/lokesh2201013/database"
	"github.com/lokesh2201013/models"
	"github.com/lokesh2201013/utils"
)


// Register godoc
// @Summary      Register new user
// @Description  Creates a new user in the database with hashed password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user  body  models.User  true  "User data"
// @Success      201  {object}  models.User
// @Failure      400  {object}  map[string]string
// @Failure      409  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /register [post]

func Register(c *fiber.Ctx)error{

	var user models.User
	if err:= c.BodyParser(&user); err!=nil{
	  logger.Log.Error(" Package controllers File Authcontroller", zap.Error(err), zap.String("Message", "Failed to decrypt request body"))
      return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error":"Invalid input"})
	}
    
	var existingUser models.User

	if err:=database.DB.Where("email=?",user.Username).First(&existingUser).Error;err==nil{
		logger.Log.Error(" Package controllers File Authcontroller", zap.Error(err), zap.String("Message", "Similar Name"))
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Email already in use"})
	}

	if err:=database.DB.Where("email=?",user.Email).First(&existingUser).Error;err==nil{
			  logger.Log.Error(" Package controllers File Authcontroller", zap.Error(err), zap.String("Message", "Similar Emails"))
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Email already in use"})
	}

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
// @Summary      Login user
// @Description  Authenticates a user and returns JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        login  body  models.LoginRequest  true  "Login credentials"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string

func Login(c *fiber.Ctx)error{
	var loginData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&loginData); err != nil {
		logger.Log.Error(" Package controllers File Authcontroller", zap.Error(err), zap.String("Message", "Failed to parse Request body"))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error":"invald input"})
	}

	var user models.User

	if err:= database.DB.Where("email=?",loginData.Email).First(&user).Error;err!=nil{
		logger.Log.Error(" Package controllers File Authcontroller", zap.Error(err), zap.String("Message", "Failed to Create database input"))
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error":"User not found"})
	}

	if err := utils.CheckPassword(user.Password, loginData.Password); err != nil {
		logger.Log.Error(" Package controllers File Authcontroller", zap.Error(err), zap.String("Message", "Failed to Create database input"))
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	token, err := utils.GenerateJWT(user.UserID, user.Email)
	if err != nil {
		logger.Log.Error(" Package controllers File Authcontroller", zap.Error(err), zap.String("Message", "Failed to Create database input"))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}
    
	logger.Log.Info(" Package controllers File Authcontroller", zap.Error(err), zap.String("Message", "token created"))
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"access_token":   token,
	}) 
}