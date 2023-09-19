package api

import (
	"github.com/czQery/llg/tl"
	"github.com/gofiber/fiber/v2"
	"os"
	"sort"
	"strings"
)

type InfoSum struct {
	Build         string   `json:"build"`
	Users         []string `json:"users"`
	SelectedUsers int      `json:"selected_users"`
}

func Info(c *fiber.Ctx) error {
	var (
		users []string
		split []string
		name  string
	)

	// get files list
	files, err := os.ReadDir(tl.Config["path"].(string))
	if err != nil {
		tl.Log("api", "data - readDir error: "+err.Error(), "error")
		return c.Status(500).JSON(Response{Message: "unexpected internal error"})
	}

	// read all files in specified folder
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		name = ""

		split = strings.Split(file.Name(), ".")
		for i := 0; i < len(split)-1; i++ {
			if name == "" {
				name = split[i]
			} else {
				name = name + "." + split[i]
			}
		}

		name = strings.ReplaceAll(name, "-login", "")

		users = append(users, name)
	}

	sort.Slice(users, func(i, j int) bool {
		return strings.ToLower(users[i]) < strings.ToLower(users[j])
	})

	return c.Status(200).JSON(Response{Data: InfoSum{Build: tl.Build, Users: users, SelectedUsers: int(tl.Config["selected_users"].(float64))}})
}
