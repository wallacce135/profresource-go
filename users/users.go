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

func GetOneUser(context *fiber.Ctx) error {
	user_id := context.Params("id")
	var user models.User

	database.DBConnection.First(&user, user_id)

	if user.Username == "" {
		return context.Status(400).JSON(fiber.Map{"status": "error", "meessage": "User not found"})
	}

	return context.JSON(fiber.Map{"status": "success", "message": "User found", "data": user})

}

func DeleteOneUser(context *fiber.Ctx) error {

	user_id := context.Params("id")
	var user models.User

	database.DBConnection.First(&user, user_id)

	if user.Username == "" {
		return context.Status(400).JSON(fiber.Map{"status": 400, "message": "User not found!"})
	}

	database.DBConnection.Delete(&user, user_id)
	return context.Status(200).JSON(fiber.Map{"status": "success", "message": "User successfully deleted from database"})

}

func CreateUser(context *fiber.Ctx) error {

	user := new(models.User)

	if err := context.BodyParser(user); err != nil {
		return context.Status(400).JSON(err.Error())
	}

	hash, err := hashPassword(user.Password)

	if err != nil {
		return context.Status(500).JSON(fiber.Map{"status": "500", "message": "Can't hash users's password", "data": err})
	}

	user.Password = hash
	database.DBConnection.Create(&user)

	return context.Status(200).JSON(user)

}
