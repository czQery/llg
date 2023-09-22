package api

import (
	"github.com/czQery/llg/tl"
	"github.com/gofiber/fiber/v2"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type DataSum struct {
	Dates []int64    `json:"dates"`
	Users []DataUser `json:"users"`
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
		dir                = tl.Config["path"].(string)
		dataDates          []int64
		dataUsers          []DataUser
		searchDateList     = make(map[string]int64)
		searchDateListKeys = make([]string, 0, len(searchDateList))

		dateParam  time.Time
		usersParam []string
	)

	if c.Query("date") != "" {
		var err error
		dateDecoded, _ := url.QueryUnescape(c.Query("date"))
		dateParam, err = time.Parse(timeLayoutParam, dateDecoded)
		if err != nil {
			dateParam = time.Now()
		}

	} else {
		dateParam = time.Now()
	}

	if c.Query("users") != "" {
		usersDecoded, _ := url.QueryUnescape(c.Query("users"))
		usersParam = strings.Split(usersDecoded, ",")
	} else {
		return c.Status(400).JSON(Response{Message: "no users selected"})
	}

	// debug
	dbgStart := time.Now()
	dbgLines := 0
	tl.Log("api", "data - reading init!", "debug")

	// get files list
	files, err := os.ReadDir(dir)
	if err != nil {
		tl.Log("api", "data - readDir error: "+err.Error(), "error")
		return c.Status(500).JSON(Response{Message: "unexpected internal error"})
	}

	// read all files in specified folder
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		var hit bool
		for _, u := range usersParam {
			if strings.Contains(file.Name(), u) {
				hit = true
				break
			}
		}

		if !hit {
			continue
		}

		fileData, err := os.ReadFile(dir + file.Name())
		if err != nil {
			tl.Log("api", "data - readFile error: "+err.Error(), "error")
		}

		var (
			search            bool
			searchName        string
			searchLogin       []string
			searchSessionList []DataUserSession
		)

		// read all lines in file
		for fileLineIndex, fileLine := range strings.Split(string(fileData), "\n") {

			fileLine = strings.ReplaceAll(fileLine, "\r", "")
			fileP := strings.Split(fileLine, ";")
			dbgLines = dbgLines + 1

			// mark session start
			if fileP[0] == "login" {
				if !search || fileLineIndex == len(strings.Split(string(fileData), "\n"))-1 {
					search = true
					searchName = fileP[1]
					searchLogin = fileP
					continue
				}
			}

			// mark session end
			if fileP[0] == "logoff" && search {
				date, _ := time.Parse(timeLayout, searchLogin[3]+"-00:00")
				dateOver, _ := time.Parse(timeLayout, fileP[3]+"-00:00")

				timeStart, _ := time.Parse(timeLayout, "01.01.1970-"+searchLogin[4])
				timeEnd, _ := time.Parse(timeLayout, "01.01.1970-"+fileP[4])

				// date sanity check
				if date.Unix() < 0 {
					tl.Log("api", "data - session: "+fileP[1]+" invalid date: "+searchLogin[3], "warn")
					search = false
					continue
				}

				// selected month check
				if date.Year() != dateParam.Year() || date.Month() != dateParam.Month() || dateOver.Year() != dateParam.Year() || dateOver.Month() != dateParam.Month() {
					search = false
					continue
				}

				// over midnight check
				if dateOver.Unix() > date.Unix() {
					searchDateList[searchLogin[3]] = date.Unix() / 60 / 60 / 24
					searchDateList[fileP[3]] = dateOver.Unix() / 60 / 60 / 24

					searchSessionList = append(searchSessionList, DataUserSession{Date: date.Unix() / 60 / 60 / 24, Device: fileP[2], Time: []int{int(timeStart.Unix() / 60), 1440}}) // start to midnight
					searchSessionList = append(searchSessionList, DataUserSession{Date: dateOver.Unix() / 60 / 60 / 24, Device: fileP[2], Time: []int{0, int(timeEnd.Unix() / 60)}})  // midnight to end

					search = false
					continue
				}

				searchSessionList = append(searchSessionList, DataUserSession{Date: date.Unix() / 60 / 60 / 24, Device: fileP[2], Time: []int{int(timeStart.Unix() / 60), int(timeEnd.Unix() / 60)}})
				searchDateList[searchLogin[3]] = date.Unix() / 60 / 60 / 24
				search = false
			}
		}

		// add currently ongoing session
		if search {
			date, _ := time.Parse(timeLayout, searchLogin[3]+"-00:00")

			timeStart, _ := time.Parse(timeLayout, "01.01.1970-"+searchLogin[4])
			timeEnd, _ := time.Parse(timeLayout, "01.01.1970-"+time.Now().Format("15:04"))

			searchSessionList = append(searchSessionList, DataUserSession{Date: date.Unix() / 60 / 60 / 24, Device: searchLogin[2], Time: []int{int(timeStart.Unix() / 60), int(timeEnd.Unix() / 60)}})
			searchDateList[searchLogin[3]] = date.Unix() / 60 / 60 / 24
		}

		if len(searchSessionList) > 0 {
			dataUsers = append(dataUsers, DataUser{Name: searchName, Sessions: searchSessionList})
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
		dataDates = append(dataDates, searchDateList[k])
	}

	tl.Log("api", "data - reading done: "+time.Since(dbgStart).String()+" lines: "+strconv.Itoa(dbgLines), "debug")

	if len(dataDates) == 0 {
		return c.Status(404).JSON(Response{Message: "no data"})
	}

	return c.Status(200).JSON(Response{Data: DataSum{Dates: dataDates, Users: dataUsers}})
}
