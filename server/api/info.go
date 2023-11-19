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
	SelectedUsers int      `json:"selected_users"`
	Users         []string `json:"users"`
	Devices       []string `json:"devices"`
}

func Info(c *fiber.Ctx) error {

	users, usersErr := infoRead(tl.Config["users"].(string))
	devices, devicesErr := infoRead(tl.Config["devices"].(string))

	if usersErr != nil || devicesErr != nil {
		return c.Status(500).JSON(Response{Message: "unexpected internal error"})
	}

	return c.Status(200).JSON(Response{Data: InfoSum{Build: tl.Build, SelectedUsers: int(tl.Config["selected_users"].(float64)), Users: users, Devices: devices}})
}

func infoRead(folder string) ([]string, error) {
	var list []string

	// get folder files list
	files, err := os.ReadDir(folder)

	if err != nil {
		tl.Log("api", "info - infoRead: readDir error: "+err.Error(), "error")
		return list, err
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".log") {
			continue
		}

		name := ""

		split := strings.Split(file.Name(), ".")
		for d := 0; d < len(split)-1; d++ {
			if name == "" {
				name = split[d]
			} else {
				name = name + "." + split[d]
			}
		}

		name = strings.ReplaceAll(name, "-login", "")

		list = append(list, name)
	}

	sort.Slice(list, func(i, j int) bool {
		return strings.ToLower(list[i]) < strings.ToLower(list[j])
	})

	return list, nil
}
