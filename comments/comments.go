package comments

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/wallacce135/profresource/database"
	"github.com/wallacce135/profresource/models"
	"github.com/wallacce135/profresource/users"
)

func GetAllComments(context *fiber.Ctx) error {

	comments := []models.Comments{}
	if err := database.DBConnection.Find(&comments).Error; err != nil {

		return context.Status(400).JSON(fiber.Map{
			"status": "error",
			"error":  err,
		})

	}

	return context.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Comments processed successfully",
		"data":    comments,
	})

}

func PostNewComment(context *fiber.Ctx) error {

	type CommentInput struct {
		Text       string `json:"text"`
		Article_id string `json:"article_id"`
	}

	var ci CommentInput
	var comment models.Comments

	if err := context.BodyParser(&ci); err != nil {

		return context.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	art_id, err := strconv.ParseUint(ci.Article_id, 10, 32)

	if err != nil {
		return context.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Error while converting types!",
		})
	}

	comment.ArticleId = uint(art_id)
	comment.Text = ci.Text

	user_id, err := users.GetUserIdFromToken(context)

	if err != nil {
		return context.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "User not found!",
		})
	}

	comment.UserId = user_id

	if err := database.DBConnection.Create(&comment).Error; err != nil {
		return context.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Error while creating comment!",
		})
	}

	return context.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Comment created successfully",
		"data":    comment,
	})

}
