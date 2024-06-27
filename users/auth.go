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

func Register(context *fiber.Ctx) error {

	type NewUser struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}

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

	newUser := NewUser{
		Email:    user.Email,
		Username: user.Username,
	}

	return context.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "User created!",
		"data":    newUser,
	})

}

func Login(context *fiber.Ctx) error {

	type LoginInput struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	type UserData struct {
		ID       uint   `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	input := new(LoginInput)
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
		return context.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Invalid identity or password", "data": err})
	} else {
		ud = UserData{
			ID:       um.ID,
			Username: um.Username,
			Email:    um.Email,
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

// func validToken(t *jwt.Token, id string) bool {

// 	n, err := strconv.Atoi(id)
// 	if err != nil {
// 		return false
// 	}

// 	claims := t.Claims.(jwt.MapClaims)
// 	uid := int(claims["user_id"].(float64))

// 	return uid == n
// }

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
