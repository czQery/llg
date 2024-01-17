package main

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"

	"github.com/czQery/llg/backend/api"
	"github.com/czQery/llg/backend/tl"
	"github.com/gofiber/fiber/v2"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		ForceColors:     true,
		DisableColors:   false,
		FullTimestamp:   true,
		TimestampFormat: "02/01/2006 - 15:04:05",
	})
	log.SetOutput(os.Stdout)

	if tl.Build == "dev" {
		log.SetLevel(log.DebugLevel)
	}
}

func main() {
	log.Info("main - Login/Logoff Graph")
	log.Info("main - by: Štěpán Aubrecht")
	log.Info("main - build: " + tl.Build)

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

	log.Info("fiber - started")

	// Run
	err := r.Listen(tl.Config["address"].(string))

	log.WithFields(log.Fields{
		"error": err.Error(),
	}).Panic("fiber - server failed")
}
