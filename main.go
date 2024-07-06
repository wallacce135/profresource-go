package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/wallacce135/profresource/database"
	"github.com/wallacce135/profresource/env"
	"github.com/wallacce135/profresource/routes"
)

// @title Profresource go Application
// @version 1.0
// description skill demontation API for profresource application

// @BasePath /
// @host localhost:4000

// @securityDefinitions.apiKey ApiKeyAuth
// @in header
// @name Authorization
func main() {

	// fmt.Println(utils.GenerateIUserId())
	env.LoadEnvironments()
	database.ConnectToDatabase()

	app := fiber.New()

	routes.SetupRoutes(app)
	app.Use(cors.New())

	log.Fatal(app.Listen(":4000"))
}
