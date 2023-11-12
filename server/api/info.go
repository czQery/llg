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
	var (
		files []os.DirEntry
		err   error

		tmp     []string
		users   []string
		devices []string
		split   []string
		name    string
	)

	// read all files in specified folder
	for i := 0; i < 2; i++ {
		if i == 0 {
			// get users files list
			files, err = os.ReadDir(tl.Config["users"].(string))
		} else {
			// get devices files list
			files, err = os.ReadDir(tl.Config["devices"].(string))
		}

		if err != nil {
			tl.Log("api", "data - readDir error: "+err.Error(), "error")
			return c.Status(500).JSON(Response{Message: "unexpected internal error"})
		}

		for _, file := range files {
			if file.IsDir() || !strings.HasSuffix(file.Name(), ".log") {
				continue
			}

			name = ""

			split = strings.Split(file.Name(), ".")
			for d := 0; d < len(split)-1; d++ {
				if name == "" {
					name = split[d]
				} else {
					name = name + "." + split[d]
				}
			}

			name = strings.ReplaceAll(name, "-login", "")

			tmp = append(tmp, name)
		}

		sort.Slice(tmp, func(i, j int) bool {
			return strings.ToLower(tmp[i]) < strings.ToLower(tmp[j])
		})

		if i == 0 {
			// save users list
			users = tmp
		} else {
			// save devices list
			devices = tmp
		}

		tmp = nil
	}

	return c.Status(200).JSON(Response{Data: InfoSum{Build: tl.Build, SelectedUsers: int(tl.Config["selected_users"].(float64)), Users: users, Devices: devices}})
}
