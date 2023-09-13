package api

import (
	"github.com/czQery/llg/tl"
	"github.com/gofiber/fiber/v2"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

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

const timeLayout = "02.01.2006-15:04"
const timeLayoutParam = "2006-01"

func Data(c *fiber.Ctx) error {
	var (
		dir                = os.Args[1]
		dataDates          []int64
		dataUsers          []DataUser
		dataSessionsPerDay = 1
		searchDateList     = make(map[string]int64)
		searchDateListKeys = make([]string, 0, len(searchDateList))

		dateParam time.Time
	)

	if c.Query("date") != "" {
		var err error
		dateParam, err = time.Parse(timeLayoutParam, c.Query("date"))
		if err != nil {
			dateParam = time.Now()
		}

	} else {
		dateParam = time.Now()
	}

	// debug
	dbgStart := time.Now()
	dbgLines := 0
	tl.Log("API", "Data - reading init!", "debug")

	// get files list
	files, err := os.ReadDir(dir)
	if err != nil {
		tl.Log("readDir", err.Error(), "error")
	}

	// read all files in specified folder
	for _, file := range files {
		if !file.IsDir() {
			fileData, err := os.ReadFile(dir + file.Name())
			if err != nil {
				tl.Log("readFile", err.Error(), "error")
			}

			var (
				search               bool
				searchName           string
				searchLogin          []string
				searchSessionList    []DataUserSession
				searchSessionsPerDay int
			)

			// read all lines in file
			for _, fileLine := range strings.Split(string(fileData), "\n") {

				fileLine = strings.ReplaceAll(fileLine, "\r", "")
				fileP := strings.Split(fileLine, ";")
				dbgLines = dbgLines + 1

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
					dateOver, _ := time.Parse(timeLayout, fileP[3]+"-00:00")

					timeStart, _ := time.Parse(timeLayout, "01.01.1970-"+searchLogin[4])
					timeEnd, _ := time.Parse(timeLayout, "01.01.1970-"+fileP[4])

					// date sanity check
					if date.Unix() < 0 {
						tl.Log("API", "Data - session: "+fileP[1]+" invalid date: "+searchLogin[3], "warn")
						search = false
						continue
					}

					// selected month check
					if date.Year() != dateParam.Year() || date.Month() != dateParam.Month() || dateOver.Year() != dateParam.Year() || dateOver.Month() != dateParam.Month() {
						search = false
						continue
					}

					// count number of session in same day for this day
					for i := len(searchSessionList) - 1; i >= 0; i-- {
						if searchSessionList[i].Date == date.Unix() {
							searchSessionsPerDay = searchSessionsPerDay + 1
						}
					}
					if dataSessionsPerDay < searchSessionsPerDay {
						dataSessionsPerDay = searchSessionsPerDay
					}

					// over midnight check
					if timeStart.Unix() > timeEnd.Unix() {
						searchDateList[searchLogin[3]] = date.Unix()
						searchDateList[fileP[3]] = dateOver.Unix()

						searchSessionList = append(searchSessionList, DataUserSession{Date: date.Unix(), Device: fileP[2], Time: []int{int(timeStart.Unix() / 60), 1440}}) // start to midnight
						searchSessionList = append(searchSessionList, DataUserSession{Date: dateOver.Unix(), Device: fileP[2], Time: []int{0, int(timeEnd.Unix() / 60)}})  // midnight to end

						search = false
						continue
					}

					searchSessionList = append(searchSessionList, DataUserSession{Date: date.Unix(), Device: fileP[2], Time: []int{int(timeStart.Unix() / 60), int(timeEnd.Unix() / 60)}})
					searchDateList[searchLogin[3]] = date.Unix()
					search = false
				}
			}
			if len(searchSessionList) > 0 {
				dataUsers = append(dataUsers, DataUser{Name: searchName, Sessions: searchSessionList})
			}
		}
	}

	tl.Log("API", "Data - sessions per day: "+strconv.Itoa(dataSessionsPerDay), "debug")

	// fill blank sessions
	for u := 0; u < len(dataUsers); u++ {

		var (
			sDay      int64
			sDayCount int
		)

		for si := 0; si < len(dataUsers[u].Sessions); si++ {
			if sDay == 0 || sDayCount == dataSessionsPerDay {
				sDay = dataUsers[u].Sessions[si].Date
				sDayCount = 1

				// fill missing end
				if si == len(dataUsers[u].Sessions)-1 {
					for i := 0; i < dataSessionsPerDay-1; i++ {
						dataUsers[u].Sessions = append(dataUsers[u].Sessions, DataUserSession{})
					}
					break
				}
				continue
			}

			sDayCount = sDayCount + 1

			if sDay != dataUsers[u].Sessions[si].Date {
				// fill missing sessions in day
				if sDayCount <= dataSessionsPerDay {
					dataUsers[u].Sessions = append(dataUsers[u].Sessions, DataUserSession{})
					copy(dataUsers[u].Sessions[si+1:], dataUsers[u].Sessions[si:])
					dataUsers[u].Sessions[si] = DataUserSession{Date: sDay, Time: []int{0, 0}}
				} else {
					sDay = 0
				}
			}

			// fill missing end
			if si == len(dataUsers[u].Sessions)-1 && sDayCount != dataSessionsPerDay {
				for i := 0; i < dataSessionsPerDay-sDayCount; i++ {
					dataUsers[u].Sessions = append(dataUsers[u].Sessions, DataUserSession{})
				}
				break
			}
		}
	}

	// format dates
	for key := range searchDateList {
		searchDateListKeys = append(searchDateListKeys, key)
	}
	sort.SliceStable(searchDateListKeys, func(i, j int) bool {
		return searchDateList[searchDateListKeys[i]] < searchDateList[searchDateListKeys[j]]
	})

	for _, k := range searchDateListKeys {
		for i := 0; i < dataSessionsPerDay; i++ {
			dataDates = append(dataDates, searchDateList[k])
		}
	}

	tl.Log("API", "Data - reading done: "+time.Since(dbgStart).String()+" lines: "+strconv.Itoa(dbgLines), "debug")

	return c.Status(200).JSON(Response{Data: DataSum{Dates: dataDates, Users: dataUsers, SessionsPerDay: dataSessionsPerDay}})
}
