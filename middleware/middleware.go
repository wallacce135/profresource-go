package middleware

import (
	"os"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

func Protected() fiber.Handler {

	secret := os.Getenv("APPLICATION_SECRET")

	return jwtware.New(jwtware.Config{
		SigningKey:   jwtware.SigningKey{Key: []byte(secret)},
		ErrorHandler: jwtError,
	})
}

func jwtError(context *fiber.Ctx, err error) error {
	if err.Error() == "Missing or mailformed JWT" {
		return context.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Missing or mailformed JWT", "data": nil})
	}

	return context.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Invalid or expired JWT", "data": nil})
}
