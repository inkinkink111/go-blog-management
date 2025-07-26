package config

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
)

func NewFiberConfig() fiber.Config {
	return fiber.Config{
		AppName:     "Go Blog Management",
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	}
}
