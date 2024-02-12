package api

import (
	"github.com/czQery/llg/backend/tl"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"os"
	"sort"
	"strings"
)

type InfoSum struct {
	Build         string   `json:"build"`
	SelectedUsers int      `json:"selected_items"`
	Users         []string `json:"users"`
	Devices       []string `json:"devices"`
}

func Info(c *fiber.Ctx) error {

	users, usersErr := infoRead(tl.Config["path_users"].(string))
	devices, devicesErr := infoRead(tl.Config["path_devices"].(string))

	if usersErr != nil || devicesErr != nil {
		return c.Status(500).JSON(Response{Message: "unexpected internal error"})
	}

	return c.Status(200).JSON(Response{Data: InfoSum{Build: tl.Build, SelectedUsers: int(tl.Config["selected_items"].(float64)), Users: users, Devices: devices}})
}

func infoRead(folder string) ([]string, error) {
	var list []string

	// get folder files list
	files, err := os.ReadDir(folder)

	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("api - infoRead: readDir")
		return list, err
	}

	var name string

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".log") {
			continue
		}

		name = file.Name()

		name = strings.TrimSuffix(name, ".log")
		name = strings.TrimSuffix(name, "-login")

		list = append(list, name)
	}

	// alphabetically sort list
	sort.Slice(list, func(i, j int) bool {
		return strings.ToLower(list[i]) < strings.ToLower(list[j])
	})

	return list, nil
}
