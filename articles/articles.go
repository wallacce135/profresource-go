package articles

import (
	"github.com/gofiber/fiber/v2"
)

func GetAllAricles(context *fiber.Ctx) error {
	return context.JSON(fiber.Map{"status": "success", "message": "All products", "data": "data after database"})
}
