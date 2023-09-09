package main

import (
	"encoding/json"
	"github.com/czQery/llg/api"
	"github.com/czQery/llg/tl"
	"github.com/gofiber/fiber/v2"
)

func main() {
	tl.Log("fiber", "starting...", "info")

	r := fiber.New(fiber.Config{
		CaseSensitive:         false,
		DisableStartupMessage: true,
		GETOnly:               false,
		BodyLimit:             100 * 1024 * 1024,
		JSONEncoder:           json.Marshal,
		JSONDecoder:           json.Unmarshal,
	})

	// API
	r.Get("/api/data", api.Data)

	// Static files
	r.Static("/", "./dist")

	// Default
	r.Use(func(c *fiber.Ctx) error {
		return c.Status(404).JSON(api.Response{Message: "Unknown endpoint!"})
	})

	// Run
	err := r.Listen(":8893")
	tl.Log("fiber", err.Error(), "error")
}
