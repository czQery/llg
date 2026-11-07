package api

import "github.com/gofiber/fiber/v3"

type Response struct {
	Message string `json:"message,omitempty"`
	Data    any    `json:"data"`
}

func Health(c fiber.Ctx) error {
	return c.Status(200).JSON(Response{Message: "up & running"})
}
