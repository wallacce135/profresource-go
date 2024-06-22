package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/wallacce135/profresource/articles"
	"github.com/wallacce135/profresource/users"
)

func SetupRoutes(app *fiber.App) {
	articlesRouter := app.Group("/api", logger.New())
	articlesRouter.Get("/articles", articles.GetAllAricles)

	usersRouter := app.Group("/users", logger.New())
	usersRouter.Get("/", users.GetAllUsers)
}
