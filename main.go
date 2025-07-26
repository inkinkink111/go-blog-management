package main

import (
	"log"

	"inkinkink111/go-blog-management/config"
	"inkinkink111/go-blog-management/db"
	"inkinkink111/go-blog-management/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	//
	app := fiber.New(config.NewFiberConfig())
	app.Use(logger.New(logger.Config{
		TimeFormat: "02-Jan-2006 15:04:05",
		TimeZone:   "Asia/Bangkok",
	}))
	app.Use(recover.New())
	app.Use(cors.New())

	db.ConnectMongo()
	db.ConnectRedis()

	// mongoClient := db.NewMongoClient(10)
	// userRepo := repositories.NewUsersDB(mongoClient)

	routes.SetupRoutes(app)

	app.Listen(":3000")
}
