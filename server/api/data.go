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

type DataSum struct {
	Dates          []int64    `json:"dates"`
	Users          []DataUser `json:"users"`
	SessionsPerDay int        `json:"sessionsPerDay"`
}

type DataUser struct {
	Name     string            `json:"name"`
	Sessions []DataUserSession `json:"sessions"`
}

type DataUserSession struct {
	Date   int64  `json:"date,omitempty"`
	Time   []int  `json:"time,omitempty"`
	Device string `json:"device,omitempty"`
}

func Data(c *fiber.Ctx) error {
	var (
		dir            = os.Args[1]
		users          []DataUser
		sessionsPerDay = 1
	)

	files, err := os.ReadDir(dir)
	if err != nil {
		tl.Log("readDir", err.Error(), "error")
	}

	// read all files in specified folder
	for _, file := range files {
		if !file.IsDir() {
			tl.Log("API", "Data - reading log: "+file.Name(), "debug")

			fileData, err := os.ReadFile(dir + file.Name())
			if err != nil {
				tl.Log("readFile", err.Error(), "error")
			}

			var (
				search               bool
				searchName           string
				searchLogin          []string
				searchList           []DataUserSession
				searchSessionsPerDay int
			)

			// read all lines in file
			for _, fileLine := range strings.Split(string(fileData), "\r\n") {
				fileP := strings.Split(fileLine, ";")

				// mark session start
				if fileP[0] == "login" && !search {
					search = true
					searchName = fileP[1]
					searchLogin = fileP
					searchSessionsPerDay = 1
				}

				// mark session end
				if fileP[0] == "logoff" && search {
					date, _ := time.Parse(timeLayout, searchLogin[3]+"-00:00")

					timeStart, _ := time.Parse(timeLayout, "01.01.1970-"+searchLogin[4])
					timeEnd, _ := time.Parse(timeLayout, "01.01.1970-"+fileP[4])

					tl.Log("API", "Data - session: "+fileP[1]+" from: "+strconv.Itoa(int(timeStart.Unix()/60))+" to: "+strconv.Itoa(int(timeEnd.Unix()/60)), "debug")

					// count number of session in same day for this day
					for i := len(searchList) - 1; i >= 0; i-- {
						if searchList[i].Date == date.Unix() {
							searchSessionsPerDay = searchSessionsPerDay + 1
						}
					}
					if sessionsPerDay < searchSessionsPerDay {
						sessionsPerDay = searchSessionsPerDay
					}

					searchList = append(searchList, DataUserSession{Date: date.Unix(), Device: fileP[2], Time: []int{int(timeStart.Unix() / 60), int(timeEnd.Unix() / 60)}})
					search = false
				}
			}
			users = append(users, DataUser{Name: searchName, Sessions: searchList})
		}
	}

	tl.Log("API", "Data - sessions per day: "+strconv.Itoa(sessionsPerDay), "debug")

	// post parse adjustments
	for u := 0; u < len(users); u++ {

		var (
			sDay      int64
			sDayCount int
		)

		for si := 0; si < len(users[u].Sessions); si++ {
			if sDay == 0 || sDayCount == sessionsPerDay {
				sDay = users[u].Sessions[si].Date
				sDayCount = 1

				// fill missing end
				if si == len(users[u].Sessions)-1 {
					for i := 0; i < sessionsPerDay-1; i++ {
						users[u].Sessions = append(users[u].Sessions, DataUserSession{})
					}
					break
				}
				continue
			}

			sDayCount = sDayCount + 1

			if sDay != users[u].Sessions[si].Date {
				// fill missing sessions in day
				if sDayCount <= sessionsPerDay {
					users[u].Sessions = append(users[u].Sessions, DataUserSession{})
					copy(users[u].Sessions[si+1:], users[u].Sessions[si:])
					users[u].Sessions[si] = DataUserSession{}
				} else {
					sDay = 0
				}
			}

			// fill missing end
			if si == len(users[u].Sessions)-1 && sDayCount != sessionsPerDay {
				for i := 0; i < sessionsPerDay-sDayCount; i++ {
					users[u].Sessions = append(users[u].Sessions, DataUserSession{})
				}
				break
			}
		}
	}

	return c.Status(200).JSON(Response{Data: DataSum{Dates: []int64{}, Users: users, SessionsPerDay: sessionsPerDay}})
}
