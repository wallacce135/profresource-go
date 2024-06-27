package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/wallacce135/profresource/articles"
	"github.com/wallacce135/profresource/comments"
	"github.com/wallacce135/profresource/middleware"
	"github.com/wallacce135/profresource/users"
)

func SetupRoutes(app *fiber.App) {
	articlesRouter := app.Group("/articles", logger.New())
	articlesRouter.Get("/", middleware.Protected(), articles.GetAllAricles)
	articlesRouter.Get("/:id", middleware.Protected(), articles.GetArticleById)
	articlesRouter.Post("/create", middleware.Protected(), articles.PostNewArticle)
	articlesRouter.Delete("/:id", middleware.Protected(), articles.DeleteArticle)
	articlesRouter.Put("/:id", middleware.Protected(), articles.UpdateArticle)

	usersRouter := app.Group("/users", logger.New())
	usersRouter.Get("/", middleware.Protected(), users.GetAllUsers)
	usersRouter.Post("/create", middleware.Protected(), users.CreateUser)
	usersRouter.Get("/:id", middleware.Protected(), users.GetOneUser)
	usersRouter.Delete("/:id", middleware.Protected(), users.DeleteOneUser)
	usersRouter.Post("/register", users.Register)
	usersRouter.Post("/login", users.Login)

	commentsRouter := app.Group("/comments", logger.New())
	commentsRouter.Get("/", middleware.Protected(), comments.GetAllComments)
	commentsRouter.Post("/create", middleware.Protected(), comments.PostNewComment)
}
