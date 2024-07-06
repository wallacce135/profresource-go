package articles

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/wallacce135/profresource/database"
	"github.com/wallacce135/profresource/models"
	"github.com/wallacce135/profresource/users"
)

// GetAllArticles returns all articles
// @Description Get all articles
// @Summary Get all articles
// @Tags Articles
// @Accept json
// @Produce json
// @Success 200  {object}  []models.Articles
// @Failure 401 {object} error
// @Router /articles [get]
func GetAllAricles(context *fiber.Ctx) error {
	articles := []models.Articles{}
	database.DBConnection.Find(&articles)
	return context.JSON(fiber.Map{"status": "success", "message": "Articles found successfully", "data": articles})
}

// GetArticleByID return one article with provided ID
// @Description Get one article with provided ID
// @Summary Get one article with provided ID
// @Tags Article
// @Accept json
// @Produce json
// @Param id path int true "Article ID" " "
// @Success 200 {object} models.Articles
// @Failure 401 {object} error
// @Failure 404 {string} string "Article not found"
// @Router /articles/:id [get]
func GetArticleById(context *fiber.Ctx) error {

	var article models.Articles

	article_id := context.Params("id")

	database.DBConnection.First(&article, article_id)

	if article.ID == 0 && article.Title == "" {
		return context.Status(404).JSON(fiber.Map{
			"message": "Article does not exist",
		})
	}

	return context.Status(200).JSON(article)

}

type ArticleBody struct {
	Title string `json:"title"`
	Text  string `json:"text:"`
}

// PostNewArticle creates a new article
// @Description Creating a new article
// @Summary Creating a new article
// @Tags Article
// @Accept json
// @Produce json
// @Param request body ArticleBody true "Body information for article creation"
// @Success 200 {object} models.Articles
// @Failure 401 {object} error
// @Failure 404 {string} string "Article not found"
// @Router /articles/create [post]
func PostNewArticle(context *fiber.Ctx) error {

	article := new(models.Articles)

	if err := context.BodyParser(article); err != nil {
		return context.Status(400).JSON(err.Error())
	}

	user_id, err := users.GetUserIdFromToken(context)

	if err != nil {
		context.Status(401).JSON(fiber.Map{
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

// UpdateArticle Update information about an article
// @Description Update information about an article
// @Summary Update information about an article
// @Tags Article
// @Accept json
// @Produce json
// @Param request body ArticleBody true "Body information for article update"
// @Param id path int true "Article ID"
// @Success 200 {object} models.Articles
// @Failure 401 {object} error
// @Failure 404 {string} string "Article not found"
// @Router /articles/:id [put]
func UpdateArticle(context *fiber.Ctx) error {

	var artInput ArticleBody

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

// DeleteArticle Deletes one article with provided ID
// @Description Delete one article with provided ID
// @Summary Delete one article with provided ID from database
// @Tags Article
// @Accept json
// @Produce json
// @Param id path int true "Article ID"
// @Success 200 {string} string "Article deleted successfully"
// @Failure 401 {object} error
// @Failure 404 {string} string "Article not found"
// @Failure 500 {string} string "You unable to delete this article!"
// @Router /articles/:id [delete]
func DeleteArticle(context *fiber.Ctx) error {

	article_id := context.Params("id")

	var article models.Articles

	database.DBConnection.First(&article, article_id)

	if article.Title == "" {
		return context.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Article not found",
			"data":    nil,
		})
	}

	user_id, err := users.GetUserIdFromToken(context)

	if err != nil {
		context.Status(401).JSON(fiber.Map{
			"status":  "error",
			"message": "User with this token does not exist",
		})
	}

	if article.UserId != user_id {
		return context.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "You unable to delete this article!",
		})
	}

	// database.DBConnection.Delete(&article, article_id)
	database.DBConnection.Model(&article).Update("IsRemoved", 1)

	return context.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Article deleted successfully",
	})

}
