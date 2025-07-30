package main

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	//"github.com/gofiber/fiber/v2/middleware/limiter"
    //"github.com/joho/godotenv"
	logger "github.com/lokesh2201013/Logger"
	"github.com/lokesh2201013/database"
	"github.com/lokesh2201013/routes"

	"github.com/lokesh2201013/docs"             
	fiberSwagger "github.com/swaggo/fiber-swagger" 
	_ "github.com/swaggo/files"  
	"sync"      
	"math"        
	"os"   
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
    docs.SwaggerInfo.Title = "Product API"
    docs.SwaggerInfo.Description = "API for managing products with JWT authentication"
    docs.SwaggerInfo.Version = "1.0"
    docs.SwaggerInfo.Host = "localhost:8080"
    docs.SwaggerInfo.BasePath = "/"
    docs.SwaggerInfo.Schemes = []string{"http"}
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	app.Use(tokenBucketMiddleware)
	routes.AuthRoutes(app)

port := ":" + os.Getenv("PORT")
if port == ":" {
    port = ":8080" 
}

	log.Printf("Server is running on http://localhost%s\n", port)
	log.Fatal(app.Listen(port))
}


type Bucket struct {
	Tokens         int
	LastRefillTime time.Time
}

var buckets = make(map[string]*Bucket)
var mu sync.Mutex

const maxTokens = 100
const refillRate = 1 

func tokenBucketMiddleware(c *fiber.Ctx) error {
	ip := c.IP()

	mu.Lock()
	bucket, exists := buckets[ip]
	if !exists {
		bucket = &Bucket{Tokens: maxTokens, LastRefillTime: time.Now()}
		buckets[ip] = bucket
	}

	now := time.Now()
	elapsed := now.Sub(bucket.LastRefillTime).Seconds()
	newTokens := int(elapsed * float64(refillRate))

	if newTokens > 0 {
		bucket.Tokens = int(math.Min(float64(maxTokens), float64(bucket.Tokens+newTokens)))
		bucket.LastRefillTime = now
	}

	if bucket.Tokens > 0 {
		bucket.Tokens--
		mu.Unlock()
		return c.Next()
	}

	mu.Unlock()
	return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
		"error": "Too many requests. Please slow down.",
	})
}
