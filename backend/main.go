package main

import (
	"os"
	"strings"

	"net/http"
	_ "net/http/pprof"

	"github.com/czQery/llg/backend/api"
	"github.com/czQery/llg/backend/tl"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/gofiber/fiber/v3/middleware/static"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		ForceColors:     true,
		DisableColors:   false,
		FullTimestamp:   true,
		TimestampFormat: "02.01.2006;15:04:05",
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

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	tl.LoadConfig()
	tl.LoadDist()

	// Fiber
	r := fiber.New(fiber.Config{
		CaseSensitive: false,
		GETOnly:       false,
		BodyLimit:     100 * 1024 * 1024,
		JSONEncoder:   json.Marshal,
		JSONDecoder:   json.Unmarshal,
		ServerHeader:  "llg",
	})

	// Compression
	r.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	// API
	r.Get("/api/health", api.Health)
	r.Get("/api/info", api.Info)
	r.Get("/api/data", api.Data)

	// Static files
	r.Get("/*", static.New("./dist"))

	// Default
	r.Use(func(c fiber.Ctx) error {
		if strings.HasPrefix(string(c.Request().URI().Path()), "/api") {
			return c.Status(404).JSON(api.Response{Message: "unknown endpoint"})
		} else if !tl.Dist {
			return c.Status(503).JSON(api.Response{Message: "front-end unavailable"})
		} else {
			return c.Redirect().Status(307).To("/")
		}
	})

	log.Info("fiber - started")

	// Run
	err := r.Listen(tl.Config["address"].(string), fiber.ListenConfig{DisableStartupMessage: true})

	log.WithFields(log.Fields{
		"error": err.Error(),
	}).Panic("fiber - server failed")
}
