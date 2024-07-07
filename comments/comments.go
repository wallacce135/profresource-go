package comments

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/wallacce135/profresource/database"
	"github.com/wallacce135/profresource/models"
	"github.com/wallacce135/profresource/users"
)

// GetAllComments returns all comments
// @Description Returns all comments from database
// @Summary Returns all comments from database
// @Tags Comments
// @Produce json
// @Success 200 {object} models.HttpResponse{Data=[]models.Comments}
// @Failure 400 {object} models.HttpResponse{}
// @Failure 401 {object} models.HttpResponse{}
// @Router /comments [get]
// @Security Bearer
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

type CommentInput struct {
	Text string `json:"text"`
}

// PostNewComment Creates a new comment
// @Description Creates a new comment in database, user id provided from JWT token
// @Summary Creates a new comment in database
// @Tags Comments
// @Accept json
// @Produce json
// @Param request body CommentInput true "Body information for new comment"
// @Param article_id query int false "article ID for comment creation"
// @Success 200 {object} models.HttpResponse{Data=models.Comments}
// @Failure 400 {object} models.HttpResponse{}
// @Failure 401 {object} models.HttpResponse{}
// @Router /comments/create [post]
// @Security Bearer
func PostNewComment(context *fiber.Ctx) error {

	article_id := context.Query("article_id")
	fmt.Println(article_id)

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

// DeleteOneComment Deletes a comment from database
// @Description Deletes user's comment from database(change flag isRemoved to 1)
// @Summary Deletes user's comment from database
// @Tags Comments
// @Produce json
// @Param id path int true "Comment ID"
// @Success 200 {object} models.HttpResponse{Data=models.Comments}
// @Failure 400 {object} models.HttpResponse{Data=nil}
// @Failure 403 {object} models.HttpResponse{Data=nil}
// @Failure 404 {object} models.HttpResponse{Data=nil}
// @Router /comments/:id [delete]
// @Security Bearer
func DeleteOneComment(context *fiber.Ctx) error {

	comment_id := context.Params("id")
	var comment models.Comments

	database.DBConnection.First(&comment, comment_id)

	if comment.ID == 0 && comment.Text == "" {
		return context.Status(404).JSON(models.HttpResponse{
			Status:  "error",
			Message: "Comment does not exist",
			Data:    nil,
		})
	}

	user_id, err := users.GetUserIdFromToken(context)

	if err != nil {
		context.Status(400).JSON(models.HttpResponse{
			Status:  "error",
			Message: "User with this token does not exist",
			Data:    nil,
		})
	}

	if comment.UserId != user_id {
		return context.Status(403).JSON(models.HttpResponse{
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
