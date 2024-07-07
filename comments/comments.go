package comments

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/wallacce135/profresource/database"
	"github.com/wallacce135/profresource/models"
	"github.com/wallacce135/profresource/users"
)

func GetAllComments(context *fiber.Ctx) error {

	comments := []models.Comments{}
	if err := database.DBConnection.Find(&comments).Error; err != nil {

		return context.Status(400).JSON(models.HttpResponse{
			Status:  "error",
			Message: "can't find comments",
			Data:    err,
		})

	}

	return context.Status(200).JSON(models.HttpResponse{
		Status:  "success",
		Message: "Comments processed successfully",
		Data:    comments,
	})

}

func PostNewComment(context *fiber.Ctx) error {

	article_id := context.Query("article_id")
	fmt.Println(article_id)

	type CommentInput struct {
		Text string `json:"text"`
	}

	var ci CommentInput
	var comment models.Comments

	if err := context.BodyParser(&ci); err != nil {

		return context.Status(400).JSON(models.HttpResponse{
			Status:  "error",
			Message: "can't parse comment JSON body",
			Data:    err,
		})
	}

	comment.Text = ci.Text

	article_ui64, err := strconv.ParseUint(article_id, 10, 64)
	if err != nil {
		panic(err)
	}

	comment.ArticleId = uint(article_ui64)

	user_id, err := users.GetUserIdFromToken(context)

	if err != nil {
		return context.Status(400).JSON(models.HttpResponse{
			Status:  "error",
			Message: "User not found!",
			Data:    err,
		})
	}

	comment.UserId = user_id

	if err := database.DBConnection.Create(&comment).Error; err != nil {
		return context.Status(400).JSON(models.HttpResponse{
			Status:  "error",
			Message: "Error while creating comment!",
			Data:    nil,
		})
	}

	return context.Status(200).JSON(models.HttpResponse{
		Status:  "success",
		Message: "Comment created successfully",
		Data:    comment,
	})

}

func DeleteOneComment(context *fiber.Ctx) error {

	comment_id := context.Params("id")
	var comment models.Comments

	database.DBConnection.First(&comment, comment_id)

	if comment.ID == 0 && comment.Text == "" {
		return context.Status(400).JSON(models.HttpResponse{
			Status:  "error",
			Message: "Comment does not exist",
			Data:    nil,
		})
	}

	user_id, err := users.GetUserIdFromToken(context)

	if err != nil {
		context.Status(400).JSON(models.HttpResponse{
			Status:  "error",
			Message: "User with this token doee not exist",
			Data:    nil,
		})
	}

	if comment.UserId != user_id {
		return context.Status(400).JSON(models.HttpResponse{
			Status:  "error",
			Message: "You unable to delete this article!",
			Data:    nil,
		})
	}

	database.DBConnection.Model(&comment).Update("IsRemoved", 1)

	return context.Status(200).JSON(models.HttpResponse{
		Status:  "success",
		Message: "Comment deleted successfully",
		Data:    comment,
	})
}
