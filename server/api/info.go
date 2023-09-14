package api

import (
	"github.com/czQery/llg/tl"
	"github.com/gofiber/fiber/v2"
	"os"
	"strings"
)

type InfoSum struct {
	Build string   `json:"build"`
	Users []string `json:"users"`
}

func Info(c *fiber.Ctx) error {
	var (
		users []string
	)

	// get files list
	files, err := os.ReadDir(os.Args[1])
	if err != nil {
		tl.Log("API", "Data - readDir error: "+err.Error(), "error")
		return c.Status(500).JSON(Response{Message: "unexpected internal error"})
	}

	// read all files in specified folder
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		users = append(users, strings.Split(file.Name(), ".")[0])
	}

	return c.Status(200).JSON(Response{Data: InfoSum{Build: tl.Build, Users: users}})
}
