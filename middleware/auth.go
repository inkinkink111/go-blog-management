package middleware

import (
	"inkinkink111/go-blog-management/utils"

	"github.com/gofiber/fiber/v2"
)

func Authenticate(c *fiber.Ctx) error {
	header := c.Get("Authorization")

	if header == "" || header == "Bearer " {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized"})
	}

	token := header[7:]

	userId, err := utils.VerifyToken(token)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized"})
	}

	if userId == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized"})
	}

	c.Locals("userId", userId)

	return c.Next()
}
