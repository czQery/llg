package main

import (
	"encoding/json"
	"strings"

	"github.com/czQery/llg/api"
	"github.com/czQery/llg/tl"
	"github.com/gofiber/fiber/v2"
)

func main() {
	tl.Log("main", "Login/Logoff Graph", "info")
	tl.Log("main", "by: Štěpán Aubrecht", "info")
	tl.Log("main", "build: "+tl.Build, "info")

	tl.LoadConfig()
	tl.LoadDist()

	r := fiber.New(fiber.Config{
		CaseSensitive:         false,
		DisableStartupMessage: true,
		GETOnly:               false,
		BodyLimit:             100 * 1024 * 1024,
		JSONEncoder:           json.Marshal,
		JSONDecoder:           json.Unmarshal,
		ServerHeader:          "llg",
	})

	// API
	r.Get("/api/info", api.Info)
	r.Get("/api/data", api.Data)

	// Static files
	r.Static("/", "./dist")

	// Default
	r.Use(func(c *fiber.Ctx) error {
		if strings.HasPrefix(string(c.Request().URI().Path()), "/api") {
			return c.Status(404).JSON(api.Response{Message: "unknown endpoint"})
		} else if !tl.Dist {
			return c.Status(503).JSON(api.Response{Message: "front-end unavailable"})
		} else {
			return c.Redirect("/", 307)
		}
	})

	tl.Log("fiber", "started!", "info")

	// Run
	err := r.Listen(tl.Config["address"].(string))
	tl.Log("fiber", err.Error(), "error")
}
