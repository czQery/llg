package api

import (
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/czQery/llg/tl"
	"github.com/gofiber/fiber/v2"
)

type DataSum struct {
	Dates []int64    `json:"dates"`
	Items []DataItem `json:"items"`
}

type DataItem struct {
	Name     string            `json:"name"`
	Type     string            `json:"type"`
	Sessions []DataItemSession `json:"sessions"`
}

type DataItemSession struct {
	Date   int64  `json:"date,omitempty"`
	Time   []int  `json:"time,omitempty"`
	Detail string `json:"detail,omitempty"`
}

const timeLayout = "02.01.2006-15:04"
const timeLayoutParam = "2006-01"

func Data(c *fiber.Ctx) error {
	var (
		dataDates          []int64
		searchDateList     = make(map[string]int64)
		searchDateListKeys = make([]string, 0, len(searchDateList))

		dateParam    time.Time
		usersParam   []string
		devicesParam []string
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

	if c.Query("users") == "" && c.Query("devices") == "" {
		return c.Status(400).JSON(Response{Message: "no users or devices selected"})
	}

	if c.Query("users") != "" {
		usersDecoded, _ := url.QueryUnescape(c.Query("users"))
		usersParam = strings.Split(usersDecoded, ",")
		tl.CleanSlice(&usersParam)
	}

	if c.Query("devices") != "" {
		devicesDecoded, _ := url.QueryUnescape(c.Query("devices"))
		devicesParam = strings.Split(devicesDecoded, ",")
		tl.CleanSlice(&devicesParam)
	}

	// debug
	dbgStart := time.Now()
	dbgLines := 0
	tl.Log("api", "data - reading init!", "debug")
	tl.Log("api", "data - users: "+strings.Join(usersParam, ","), "debug")
	tl.Log("api", "data - devices: "+strings.Join(devicesParam, ","), "debug")

	dataUsers, usersErr := dataRead("user", tl.Config["users"].(string), usersParam, dateParam, &searchDateList, &dbgLines)
	dataDevices, devicesErr := dataRead("device", tl.Config["devices"].(string), devicesParam, dateParam, &searchDateList, &dbgLines)

	if usersErr != nil || devicesErr != nil {
		return c.Status(500).JSON(Response{Message: "unexpected internal error"})
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

	return c.Status(200).JSON(Response{Data: DataSum{Dates: dataDates, Items: append(dataUsers, dataDevices...)}})
}

func dataRead(tp string, folder string, items []string, date time.Time, dateList *map[string]int64, dbgLines *int) ([]DataItem, error) {
	var list []DataItem

	if len(items) == 0 {
		return list, nil
	}

	// get files list
	files, err := os.ReadDir(folder)
	if err != nil {
		tl.Log("api", "data - readDir error: "+err.Error(), "error")
		return list, err
	}

	// read all files in specified folder
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".log") {
			continue
		}

		var hit bool
		for _, u := range items {
			if strings.Contains(file.Name(), u) {
				hit = true
				break
			}
		}

		if !hit {
			continue
		}

		fileData, err := os.ReadFile(folder + file.Name())
		if err != nil {
			tl.Log("api", "data - readFile error: "+err.Error(), "error")
		}

		var (
			search            bool
			searchName        string
			searchLogin       []string
			searchLogoff      []string
			searchSessionList []DataItemSession
		)

		fileLines := strings.Split(string(fileData), "\n")

		// read all lines in file
		for fileLineIndex, fileLine := range fileLines {

			fileLine = strings.ReplaceAll(fileLine, "\r", "")
			fileP := strings.Split(fileLine, ";")
			*dbgLines = *dbgLines + 1

			if len(fileP) < 5 {
				// skip new line on end in order to have currently ongoing session
				if fileLineIndex == len(fileLines)-1 && fileLineIndex > 0 {
					fileLine = fileLines[fileLineIndex-1]
					fileLine = strings.ReplaceAll(fileLine, "\r", "")
					fileP = strings.Split(fileLine, ";")
				} else {
					continue
				}
			}

			// mark session start
			if fileP[0] == "login" {
				if !search && fileLineIndex == len(fileLines)-1 {
					search = true
					searchName = fileP[1]
					searchLogin = fileP
				}

				// fill login to midnight if (the closest login is the next day || this is the last login in log)
				if search && (searchLogin[3] != fileP[3] || (fileLineIndex == len(fileLines)-1 && time.Now().Format("02.01.2006") != searchLogin[3])) {

					timeStart, _ := time.Parse(timeLayout, "01.01.1970-"+searchLogin[4])
					dateStart, _ := time.Parse(timeLayout, searchLogin[3]+"-00:00")

					// time & date sanity check
					if timeStart.Unix() < 0 || dateStart.Unix() < 0 {
						tl.Log("api", "data - session: "+searchLogin[1]+" invalid date: "+searchLogin[3]+" or time: "+searchLogin[4], "warn")
					} else {
						// selected month check
						if dateStart.Year() == date.Year() && dateStart.Month() == date.Month() {
							searchSessionList = append(searchSessionList, DataItemSession{Date: dateStart.Unix() / 60 / 60 / 24, Detail: searchLogin[2], Time: []int{int(timeStart.Unix() / 60), 1440}}) // start to midnight
							(*dateList)[searchLogin[3]] = dateStart.Unix() / 60 / 60 / 24
						}
					}

					search = false
				}

				if !search {
					search = true
					searchName = fileP[1]
					searchLogin = fileP
					continue
				}
			}

			// mark session end
			if fileP[0] == "logoff" {

				var (
					timeStart time.Time
					timeEnd   time.Time

					dateStart time.Time
					dateEnd   time.Time
				)

				timeEnd, _ = time.Parse(timeLayout, "01.01.1970-"+fileP[4])
				dateEnd, _ = time.Parse(timeLayout, fileP[3]+"-00:00")

				if search { // login as start
					timeStart, _ = time.Parse(timeLayout, "01.01.1970-"+searchLogin[4])
					dateStart, _ = time.Parse(timeLayout, searchLogin[3]+"-00:00")
				} else if len(searchLogoff) > 4 && fileP[3] == searchLogoff[3] { // missing login => previous logoff in same day as login
					timeStart, _ = time.Parse(timeLayout, "01.01.1970-"+searchLogoff[4])
					dateStart, _ = time.Parse(timeLayout, searchLogoff[3]+"-00:00")
				} else { // missing login & logoff => 00:00 as login
					timeStart, _ = time.Parse(timeLayout, "01.01.1970-00:00")
					dateStart = dateEnd
				}

				// date sanity check
				if dateStart.Unix() < 0 || dateEnd.Unix() < 0 {
					tl.Log("api", "data - session: "+fileP[1]+" invalid date: "+searchLogin[3]+", "+fileP[3], "warn")
					search = false
					continue
				}

				// time sanity check
				if timeStart.Unix() < 0 || timeEnd.Unix() < 0 {
					tl.Log("api", "data - session: "+fileP[1]+" invalid time: "+searchLogin[4]+", "+fileP[4], "warn")
					search = false
					continue
				}

				// selected month check
				if dateStart.Year() != date.Year() || dateStart.Month() != date.Month() || dateEnd.Year() != date.Year() || dateEnd.Month() != date.Month() {
					search = false
					continue
				}

				searchLogoff = fileP

				// over midnight check
				if dateEnd.Unix() > dateStart.Unix() {
					(*dateList)[searchLogin[3]] = dateStart.Unix() / 60 / 60 / 24
					(*dateList)[fileP[3]] = dateEnd.Unix() / 60 / 60 / 24

					searchSessionList = append(searchSessionList, DataItemSession{Date: dateStart.Unix() / 60 / 60 / 24, Detail: searchLogin[2], Time: []int{int(timeStart.Unix() / 60), 1440}}) // start to midnight
					searchSessionList = append(searchSessionList, DataItemSession{Date: dateEnd.Unix() / 60 / 60 / 24, Detail: fileP[2], Time: []int{0, int(timeEnd.Unix() / 60)}})              // midnight to end

					(*dateList)[searchLogin[3]] = dateStart.Unix() / 60 / 60 / 24
					(*dateList)[fileP[3]] = dateEnd.Unix() / 60 / 60 / 24

					search = false
					continue
				}

				searchSessionList = append(searchSessionList, DataItemSession{Date: dateStart.Unix() / 60 / 60 / 24, Detail: fileP[2], Time: []int{int(timeStart.Unix() / 60), int(timeEnd.Unix() / 60)}})
				(*dateList)[searchLogin[3]] = dateStart.Unix() / 60 / 60 / 24
				search = false
			}
		}

		// add currently ongoing session
		if search && time.Now().Format("02.01.2006") == searchLogin[3] {
			dateStart, _ := time.Parse(timeLayout, searchLogin[3]+"-00:00")

			timeStart, _ := time.Parse(timeLayout, "01.01.1970-"+searchLogin[4])
			timeEnd, _ := time.Parse(timeLayout, "01.01.1970-"+time.Now().Format("15:04"))

			// time & date sanity check
			if timeStart.Unix() < 0 || dateStart.Unix() < 0 {
				tl.Log("api", "data - session: "+searchLogin[1]+" invalid date: "+searchLogin[3]+" or time: "+searchLogin[4], "warn")
			} else {
				// selected month check
				if dateStart.Year() == date.Year() && dateStart.Month() == date.Month() {
					searchSessionList = append(searchSessionList, DataItemSession{Date: dateStart.Unix() / 60 / 60 / 24, Detail: searchLogin[2], Time: []int{int(timeStart.Unix() / 60), int(timeEnd.Unix() / 60)}})
					(*dateList)[searchLogin[3]] = dateStart.Unix() / 60 / 60 / 24
				}
			}
		}

		if len(searchSessionList) > 0 {

			list = append(list, DataItem{Name: searchName, Type: tp, Sessions: searchSessionList})
		}
	}

	return list, nil
}
