package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/wallacce135/profresource/database"
	"github.com/wallacce135/profresource/env"
	"github.com/wallacce135/profresource/routes"
)

func main() {

	// fmt.Println(utils.GenerateIUserId())
	env.LoadEnvironments()
	database.ConnectToDatabase()

	app := fiber.New()
	app.Use(cors.New())

	routes.SetupRoutes(app)

	log.Fatal(app.Listen(":4000"))
}
