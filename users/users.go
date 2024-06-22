package users

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wallacce135/profresource/database"
	"github.com/wallacce135/profresource/models"
)

func GetAllUsers(context *fiber.Ctx) error {

	users := []models.User{}
	database.DBConnection.Find(&users)

	return context.Status(200).JSON(users)
}
