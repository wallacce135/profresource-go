package articles

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/wallacce135/profresource/database"
	"github.com/wallacce135/profresource/models"
	"github.com/wallacce135/profresource/users"
)

func GetAllAricles(context *fiber.Ctx) error {
	articles := []models.Articles{}
	database.DBConnection.Find(&articles)
	return context.JSON(fiber.Map{"status": "success", "message": "Articles found successfully", "data": articles})
}

func GetArticleById(context *fiber.Ctx) error {

	var article models.Articles

	article_id := context.Params("id")

	database.DBConnection.First(&article, article_id)

	if article.ID == 0 && article.Title == "" {
		return context.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Article does not exist",
			"data":    nil,
		})
	}

	return context.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Article found successfully",
		"data":    article,
	})

}

func PostNewArticle(context *fiber.Ctx) error {

	article := new(models.Articles)

	if err := context.BodyParser(article); err != nil {
		return context.Status(400).JSON(err.Error())
	}

	user_id, err := users.GetUserIdFromToken(context)

	if err != nil {
		context.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "User with this token does not exist",
		})
	}

	article.UserId = user_id

	database.DBConnection.Create(&article)

	return context.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Article successfully created",
		"data":    article,
	})

}

func UpdateArticle(context *fiber.Ctx) error {

	type ArtInput struct {
		Title string `json:"title"`
		Text  string `json:"text"`
	}

	var artInput ArtInput

	var article models.Articles
	article_id := context.Params("id")

	if err := context.BodyParser(artInput); err != nil {
		return context.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Error while parsing article data!",
		})
	}

	user_id, err := users.GetUserIdFromToken(context)

	if err != nil {
		context.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "User with this token does not exist",
		})
	}

	database.DBConnection.First(&article, article_id)

	if article.UserId != user_id {
		return context.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "You unable to update this article!",
		})
	}

	article.Title = artInput.Title
	article.Text = artInput.Text
	article.UpdatedAt = time.Now()

	database.DBConnection.Save(&article)

	return context.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Article updated successfully",
	})

}

func DeleteArticle(context *fiber.Ctx) error {

	article_id := context.Params("id")

	var article models.Articles

	database.DBConnection.First(&article, article_id)

	if article.Title == "" {
		return context.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Article not found",
			"data":    nil,
		})
	}

	user_id, err := users.GetUserIdFromToken(context)

	if err != nil {
		context.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "User with this token does not exist",
		})
	}

	if article.UserId != user_id {
		return context.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "You unable to delete this article!",
		})
	}

	database.DBConnection.Delete(&article, article_id)

	return context.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Article deleted successfully",
	})

}
