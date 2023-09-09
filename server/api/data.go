package api

import (
	"github.com/czQery/llg/tl"
	"github.com/gofiber/fiber/v2"
	"os"
	"strconv"
	"strings"
	"time"
)

const timeLayout = "02.01.2006-15:04"

type DataUser struct {
	Name     string            `json:"name"`
	Sessions []DataUserSession `json:"sessions"`
}

type DataUserSession struct {
	Date   int64  `json:"date"`
	Device string `json:"device"`
	From   int    `json:"from"`
	To     int    `json:"to"`
}

func Data(c *fiber.Ctx) error {
	var (
		dir   = os.Args[1]
		users []DataUser
	)

	files, err := os.ReadDir(dir)
	if err != nil {
		tl.Log("readDir", err.Error(), "error")
	}

	for _, file := range files {
		if !file.IsDir() {
			tl.Log("API", "Data - reading log: "+file.Name(), "debug")

			fileData, err := os.ReadFile(dir + file.Name())
			if err != nil {
				tl.Log("readFile", err.Error(), "error")
			}

			var (
				search      bool
				searchName  string
				searchLogin []string
				searchList  []DataUserSession
			)

			for _, fileLine := range strings.Split(string(fileData), "\r\n") {
				fileP := strings.Split(fileLine, ";")

				if fileP[0] == "login" && !search {
					search = true
					searchName = fileP[1]
					searchLogin = fileP
				}

				if fileP[0] == "logoff" && search {
					date, _ := time.Parse(timeLayout, searchLogin[3]+"-00:00")

					timeStart, _ := time.Parse(timeLayout, "01.01.1970-"+searchLogin[4])
					timeEnd, _ := time.Parse(timeLayout, "01.01.1970-"+fileP[4])

					tl.Log("API", "Data - session: "+fileP[1]+" from: "+strconv.Itoa(int(timeStart.Unix()/60))+" to: "+strconv.Itoa(int(timeEnd.Unix()/60)), "debug")

					searchList = append(searchList, DataUserSession{Date: date.Unix(), Device: fileP[2], From: int(timeStart.Unix() / 60), To: int(timeEnd.Unix() / 60)})
					search = false
				}
			}

			users = append(users, DataUser{Name: searchName, Sessions: searchList})
		}
	}

	return c.Status(200).JSON(Response{Data: users})
}
