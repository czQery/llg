package main

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/handlebars/v2"
	"os"
	"strings"
	"time"
)

type ApiResponse struct {
	Message string `json:"message,omitempty"`
	Data    any    `json:"data"`
}

type UserTime struct {
	From string `json:"from"`
	To   string `json:"to"`
}

func main() {
	Log("fiber", "starting...", "info")

	engine := handlebars.New("./src/hbs", "")
	engine.Reload(true)
	engine.Debug(false)

	r := fiber.New(fiber.Config{
		CaseSensitive:         false,
		DisableStartupMessage: true,
		GETOnly:               false,
		BodyLimit:             100 * 1024 * 1024,
		JSONEncoder:           json.Marshal,
		JSONDecoder:           json.Unmarshal,
		Views:                 engine,
	})

	// Robots
	r.Get("/robots.txt", func(c *fiber.Ctx) error {
		c.Context().SetContentType("text/plain")
		return c.SendString("User-agent: *\nDisallow: /")
	})

	// Public
	r.Static("/src/public", "./src/public")

	// Main
	r.Get("/", func(c *fiber.Ctx) error {
		bind := fiber.Map{"data": "{\":)\"}"}

		dir := os.Args[1]
		/*users := make(map[string]map[string][]UserTime)*/

		files, err := os.ReadDir(dir + "/login/user")
		if err != nil {
			Log("readDir", err.Error(), "error")
		}

		for _, file := range files {
			if !file.IsDir() {
				fileData, err := os.ReadFile(dir + "/login/user/" + file.Name())
				if err != nil {
					Log("readFile", err.Error(), "error")
				}

				for _, fileLine := range strings.Split(string(fileData), "\r\n") {
					fileP := strings.Split(fileLine, ";")
					/*users[fileP[0]][fileP[1]]*/
					fmt.Println(fileP[0], fileP[1])
				}
			}
		}

		c.Context().SetContentType("text/html")
		return c.Render("main.hbs", bind)
	})

	// Default
	r.Use(func(c *fiber.Ctx) error {
		return c.Status(404).JSON(ApiResponse{Message: "Unknown endpoint!"})
	})

	// Run
	err := r.Listen(":8893")
	Log("fiber", err.Error(), "error")
}

func Log(prefix string, message string, messageType string) {
	now := time.Now()

	messageType = "[" + messageType + "]"

	fmt.Println("[" + now.Format("02/01/2006") + " - " + now.Format("15:04:05") + "] " + messageType + " " + prefix + " - " + message)
}
