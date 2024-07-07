package users

import (
	"errors"
	"net/mail"
	"os"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/wallacce135/profresource/database"
	"github.com/wallacce135/profresource/models"
	"golang.org/x/crypto/bcrypt"
)

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func getUserByEmail(e string) (*models.User, error) {
	var user models.User
	if err := database.DBConnection.Where(&models.User{Email: e}).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func getUserByUsername(u string) (*models.User, error) {
	var user models.User
	if err := database.DBConnection.Where(&models.User{Username: u}).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func valid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

// Register Allows user to create new user in the system
// @Description Allows you to create a new user
// @Summary Allows you to create a new user
// @Tags Users
// @Accept json
// @Produce json
// @Param request body models.NewUser true "Body information to create a new user"
// @Success 200 {object} models.HttpResponse{data=models.NewUser}
// @Failure 500 {object} models.HttpResponse{data=nil}
// @Router /users/register [post]
func Register(context *fiber.Ctx) error {

	user := new(models.User)
	if err := context.BodyParser(user); err != nil {
		return context.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "error",
			"message": "Error while parsing user data!",
		})
	}

	validate := validator.New()

	if err := validate.Struct(user); err != nil {
		return context.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "error",
			"message": "Error while validating user data struct",
		})
	}

	hash, err := hashPassword(user.Password)

	if err != nil {
		return context.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "error",
			"message": "Error while user password hashing",
		})
	}

	user.Password = hash

	if err := database.DBConnection.Create(&user).Error; err != nil {
		return context.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "error",
			"message": "Error while adding user in database",
		})
	}

	newUser := models.NewUser{
		Email:    user.Email,
		Username: user.Username,
		Password: "password hashed!",
	}

	return context.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "User created!",
		"data":    newUser,
	})

}

// Login Allows user to login in the system
// @Description Allows user to login in the system
// @Summary Allows user to login in the system
// @Tags Users
// @Accept json
// @Produce json
// @Param request body models.LoginInput true "Body information for login"
// @Success 200 {object} models.HttpResponse
// @Failure 401 {object} models.HttpResponse
// @Failure 500 {object} models.HttpResponse
// @Router /users/login [post]
func Login(context *fiber.Ctx) error {

	type UserData struct {
		ID       uint   `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	input := new(models.LoginInput)
	var ud UserData

	if err := context.BodyParser(input); err != nil {
		return context.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Error on login request", "errors": err.Error()})
	}

	userInput := input.Username
	pass := input.Password
	um, err := new(models.User), *new(error)

	if valid(userInput) {
		um, err = getUserByEmail(userInput)
	} else {
		um, err = getUserByUsername(userInput)
	}

	if err != nil {
		return context.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Internal Server Error", "data": err})
	} else if um == nil {
		CheckPasswordHash(pass, "")
		return context.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Invalid username or password", "data": err})
	} else {
		ud = UserData{
			ID:       um.ID,
			Username: um.Username,
			Password: um.Password,
		}
	}

	if !CheckPasswordHash(pass, ud.Password) {
		return context.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid credentials",
		})
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = ud.Username
	claims["user_id"] = ud.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	secret := os.Getenv("APPLICATION_SECRET")

	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return context.SendStatus(fiber.StatusInternalServerError)
	}

	return context.JSON(fiber.Map{"status": "success", "message": "Success login", "data": t})
}

func GetUserIdFromToken(context *fiber.Ctx) (uint, error) {

	tokenString := strings.Split(context.GetReqHeaders()["Authorization"][0], "Bearer ")[1]
	claims := jwt.MapClaims{}
	access := os.Getenv("APPLICATION_SECRET")

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(access), nil
	})

	if err != nil {
		context.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Error while parsing token",
		})
	}

	if !token.Valid {
		context.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid token",
		})
	}

	user_id := claims["user_id"]
	return uint(user_id.(float64)), nil

}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
