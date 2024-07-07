package users

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wallacce135/profresource/database"
	"github.com/wallacce135/profresource/models"
)

// GetAllUsers Returns all users in the database
// @Description Get all users in the database
// @Summary Get all users in the database
// @Tags Users
// @Produce json
// @Success 200 {object} []models.User{}
// @Failure 401 {object} error
// @Router /users/all [get]
// @Security Bearer
func GetAllUsers(context *fiber.Ctx) error {

	users := []models.User{}
	database.DBConnection.Find(&users)

	return context.Status(200).JSON(models.HttpResponse{
		Status:  "success",
		Message: "Users returned successfully",
		Data:    users,
	})
}

// GetAllActiveUsers Returns all active users in the database
// @Description Get all active users in the database
// @Summary Get all users in the database
// @Tags Users
// @Produce json
// @Success 200 {object} []models.User{}
// @Failure 401 {object} error
// @Router /users [get]
// @Security Bearer
func GetAllActiveUsers(context *fiber.Ctx) error {

	users := []models.User{}
	database.DBConnection.Where("is_removed = 0").Find(&users)

	return context.Status(200).JSON(models.HttpResponse{
		Status:  "success",
		Message: "Active users returned successfully",
		Data:    users,
	})

}

// GetOne Returns one user from database
// @Description Returns one user from database
// @Summary Returns one user from database
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} []models.User{Data=models.User}
// @Failure 404 {object} models.HttpResponse
// @Failure 401 {object} models.HttpResponse
// @Router /users/:id [get]
// @Security Bearer
func GetOneUser(context *fiber.Ctx) error {
	user_id := context.Params("id")
	var user models.User

	database.DBConnection.First(&user, user_id)

	if user.Username == "" {
		return context.Status(404).JSON(models.HttpResponse{
			Status:  "error",
			Message: "User not found",
			Data:    nil,
		})
	}

	return context.Status(200).JSON(models.HttpResponse{
		Status:  "success",
		Message: "User found",
		Data:    user,
	})
}

// DeleteOneUser Deletes user from database
// @Description Deletes user from database(change flag isRemoved to 1)
// @Summary Deletes user from database
// @Tags Users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} models.HttpResponse{Data=models.User}
// @Failure 400 {object} models.HttpResponse{Data=nil}
// @Failure 401 {object} models.HttpResponse{Data=nil}
// @Router /users/:id [delete]
// @Security Bearer
func DeleteOneUser(context *fiber.Ctx) error {

	user_id := context.Params("id")
	var user models.User

	database.DBConnection.First(&user, user_id)

	if user.Username == "" {
		return context.Status(400).JSON(models.HttpResponse{
			Status:  "error",
			Message: "User not found",
			Data:    nil,
		})
	}

	database.DBConnection.Model(&user).Update("IsRemoved", 1)
	return context.Status(200).JSON(models.HttpResponse{
		Status:  "success",
		Message: "User successfully deleted from database",
		Data:    user,
	})

}

// GetUserBack Restores user in database
// @Description Restores user in database(change flag isRemoved to 0)
// @Summary Restores user in database
// @Tags Users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} models.HttpResponse{Data=models.User}
// @Failure 400 {object} models.HttpResponse{Data=nil}
// @Failure 401 {object} models.HttpResponse{Data=nil}
// @Router /users/restore/:id [get]
// @Security Bearer
func GetUserBack(context *fiber.Ctx) error {

	user_id := context.Params("id")
	var user models.User

	database.DBConnection.First(&user, user_id)

	if user.Username == "" {
		return context.Status(400).JSON(models.HttpResponse{
			Status:  "error",
			Message: "User not found",
			Data:    nil,
		})
	}

	database.DBConnection.Model(&user).Update("IsRemoved", 0)
	return context.Status(200).JSON(models.HttpResponse{
		Status:  "success",
		Message: "User successfully restored",
		Data:    user,
	})

}

// CreateUser Allows loggined user to create a new user
// @Description Allows to create a new user in database for existing user
// @Summary Allows loggined user to create a new user
// @Tags Users
// @Produce json
// @Param request body models.NewUser true "Body information"
// @Success 200 {object} models.HttpResponse{Data=models.User}
// @Failure 400 {object} models.HttpResponse{Data=error}
// @Failure 401 {object} models.HttpResponse{Data=error}
// @Router /users/create [post]
// @Security Bearer
func CreateUser(context *fiber.Ctx) error {

	user := new(models.User)

	if err := context.BodyParser(user); err != nil {
		return context.Status(400).JSON(models.HttpResponse{
			Status:  "error",
			Message: "Error while parsing JSON body",
			Data:    err,
		})
	}

	hash, err := hashPassword(user.Password)

	if err != nil {
		return context.Status(500).JSON(models.HttpResponse{
			Status:  "error",
			Message: "Can't hash users's password",
			Data:    err,
		})
	}

	user.Password = hash
	database.DBConnection.Create(&user)

	return context.Status(200).JSON(models.HttpResponse{
		Status:  "success",
		Message: "User successfully created",
		Data:    user,
	})

}
