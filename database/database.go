package database

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/lokesh2201013/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found")
	}

	dsn := "host=" + os.Getenv("DB_HOST") +
		" user=" + os.Getenv("DB_USER") +
		" password=" + os.Getenv("DB_PASSWORD") +
		" dbname=" + os.Getenv("DB_NAME") +
		" port=" + os.Getenv("DB_PORT") +
		" sslmode=disable"
		log.Println("Connecting with DSN:", dsn)


	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to db: %v\n", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get generic database object: %v\n", err)
	}


	sqlDB.SetMaxOpenConns(25)              
	sqlDB.SetMaxIdleConns(10)                
	sqlDB.SetConnMaxLifetime(5 * time.Minute)  

	db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)
	if err := db.AutoMigrate(&models.User{}, &models.Product{}); err != nil {
		log.Fatalf("AutoMigrate failed: %v\n", err)
	}

	DB = db
	log.Println("Connected to PostgreSQL with connection pooling enabled")
}
