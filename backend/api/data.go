package api

import (
	"maps"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/czQery/llg/backend/tl"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

type DataSum struct {
	Dates []int64    `json:"dates"`
	Items []DataItem `json:"items"`

	Lines int   `json:"lines"`
	Time  int64 `json:"time"`
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

type dataSource struct {
	tp       string // type: user || device
	folder   string
	items    []string
	date     time.Time
	dateList map[string]int64
	itemList []DataItem
	lines    int
	mu       *sync.Mutex
}

type dataSearch struct {
	active      bool
	name        string
	login       []string
	logoff      []string
	sessionList []DataItemSession
	lines       int
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
		execTime     = time.Now()

		wg sync.WaitGroup
	)

	// query params
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

	var (
		sourceUsers = dataSource{
			tp:       "user",
			folder:   tl.Config["path_users"].(string),
			items:    usersParam,
			date:     dateParam,
			dateList: make(map[string]int64),
			itemList: []DataItem{},
			lines:    0,
			mu:       &sync.Mutex{},
		}

		sourceDevices = dataSource{
			tp:       "device",
			folder:   tl.Config["path_devices"].(string),
			items:    devicesParam,
			date:     dateParam,
			dateList: make(map[string]int64),
			itemList: []DataItem{},
			lines:    0,
			mu:       &sync.Mutex{},
		}

		usersErr   error
		devicesErr error
	)

	// parse logs
	wg.Add(2)

	go func() {
		usersErr = sourceUsers.readFolder()
		wg.Done()
	}()

	go func() {
		devicesErr = sourceDevices.readFolder()
		wg.Done()
	}()

	wg.Wait()

	if usersErr != nil || devicesErr != nil {
		return c.Status(500).JSON(Response{Message: "unexpected internal error"})
	}

	// merge date lists
	maps.Copy(searchDateList, sourceUsers.dateList)
	maps.Copy(searchDateList, sourceDevices.dateList)

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

	log.WithFields(log.Fields{
		"users":   strings.Join(usersParam, ","),
		"devices": strings.Join(usersParam, ","),
		"time":    time.Since(execTime).String(),
		"lines":   strconv.Itoa(sourceUsers.lines + sourceDevices.lines),
	}).Debug("api - data: reading done")

	if len(dataDates) == 0 {
		return c.Status(404).JSON(Response{Message: "no data"})
	}

	return c.Status(200).JSON(Response{Data: DataSum{Dates: dataDates, Items: append(sourceUsers.itemList, sourceDevices.itemList...), Lines: sourceUsers.lines + sourceDevices.lines, Time: time.Since(execTime).Nanoseconds()}})
}

func (src *dataSource) readFolder() error {

	var (
		files []os.DirEntry
		wg    sync.WaitGroup

		err error
	)

	if len(src.items) == 0 {
		return nil
	}

	// get files list
	files, err = os.ReadDir(src.folder)
	if err != nil {

		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("api - data: readDir")

		return err
	}

	// read all files in specified folder
	for _, file := range files {
		wg.Add(1)

		go func() {
			src.readFile(&file)
			wg.Done()
		}()
	}

	wg.Wait()

	return nil
}

func (src *dataSource) readFile(file *os.DirEntry) {
	var (
		fileData  []byte
		fileLines []string

		search = dataSearch{
			active:      false,
			name:        "",
			login:       []string{},
			logoff:      []string{},
			sessionList: []DataItemSession{},
			lines:       0,
		}

		err error
	)

	// check file
	if (*file).IsDir() || !strings.HasSuffix((*file).Name(), ".log") {
		return
	}

	// check file name
	for i, u := range src.items {
		if strings.Contains((*file).Name(), u) {
			break
		}

		// return on last item if file name doesn't contain requested name
		if i == len(src.items)-1 {
			return
		}
	}

	// open file
	fileData, err = os.ReadFile(src.folder + (*file).Name())
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("api - data: readFile")
		return
	}

	// split file data into lines
	fileLines = strings.Split(string(fileData), "\n")

	// read all lines in file
	for fileLineIndex, fileLine := range fileLines {
		src.readFileLine(&search, &fileLineIndex, &fileLine, &fileLines)
	}

	src.addLines(search.lines)

	// add currently ongoing session
	src.addCurrentSession(&search)

	// bundle sessions if found
	if len(search.sessionList) > 0 {
		if src.tp == "device" {
			search.name = strings.Split(search.name, " ")[0]
		}

		src.addItem(DataItem{Name: search.name, Type: src.tp, Sessions: search.sessionList})
	}

	return
}

func (src *dataSource) readFileLine(search *dataSearch, fileLineIndex *int, fileLine *string, fileLines *[]string) {

	*fileLine = strings.ReplaceAll(*fileLine, "\r", "")
	fileLineSegments := strings.Split(*fileLine, ";")

	search.lines = search.lines + 1

	if len(fileLineSegments) < 5 {
		// skip new line on end in order to have currently ongoing session
		if *fileLineIndex == len(*fileLines)-1 && *fileLineIndex > 0 {
			*fileLine = (*fileLines)[*fileLineIndex-1]
			*fileLine = strings.ReplaceAll(*fileLine, "\r", "")
			fileLineSegments = strings.Split(*fileLine, ";")
		} else {
			return
		}
	}

	switch fileLineSegments[0] {
	case "login":
		// mark session start
		src.readFileLineLogin(search, fileLineIndex, fileLines, &fileLineSegments)
	case "logoff":
		// mark session end
		src.readFileLineLogoff(search, &fileLineSegments)
	}
}

func (src *dataSource) readFileLineLogin(search *dataSearch, fileLineIndex *int, fileLines *[]string, fileLineSegments *[]string) {
	if !search.active && *fileLineIndex == len(*fileLines)-1 {
		search.active = true
		search.name = (*fileLineSegments)[1]
		search.login = *fileLineSegments
	}

	// fill login to midnight if (the closest login is the next day || this is the last login in log)
	if search.active && (search.login[3] != (*fileLineSegments)[3] || (*fileLineIndex == len(*fileLines)-1 && time.Now().Format("02.01.2006") != search.login[3])) {

		timeStart := tl.ParseTime("01.01.1970-" + search.login[4])
		dateStart := tl.ParseTime(search.login[3] + "-00:00")

		// time & date sanity check
		if timeStart.Unix() < 0 || dateStart.Unix() < 0 {
			log.WithFields(log.Fields{
				"session": search.login[1],
				"date":    search.login[3],
				"time":    search.login[4],
			}).Warn("api - data: session invalid date or time")
		} else {
			// selected month check
			if dateStart.Year() == src.date.Year() && dateStart.Month() == src.date.Month() {
				search.sessionList = append(search.sessionList, DataItemSession{Date: dateStart.Unix() / 60 / 60 / 24, Detail: src.getDetail(search.login[1], search.login[2]), Time: []int{int(timeStart.Unix() / 60), 1440}}) // start to midnight
				src.addDate(&search.login[3], dateStart.Unix()/60/60/24)
			}
		}

		search.active = false
	}

	if !search.active {
		search.active = true
		search.name = (*fileLineSegments)[1]
		search.login = *fileLineSegments
		return
	}
}

func (src *dataSource) readFileLineLogoff(search *dataSearch, fileLineSegments *[]string) {
	var (
		timeStart time.Time
		timeEnd   time.Time

		dateStart time.Time
		dateEnd   time.Time
	)

	timeEnd = tl.ParseTime("01.01.1970-" + (*fileLineSegments)[4])
	dateEnd = tl.ParseTime((*fileLineSegments)[3] + "-00:00")

	if search.active { // login as start
		timeStart = tl.ParseTime("01.01.1970-" + search.login[4])
		dateStart = tl.ParseTime(search.login[3] + "-00:00")
	} else if len(search.logoff) > 4 && (*fileLineSegments)[3] == search.logoff[3] { // missing login => previous logoff in same day as login
		timeStart = tl.ParseTime("01.01.1970-" + search.logoff[4])
		dateStart = tl.ParseTime(search.logoff[3] + "-00:00")
	} else { // missing login & logoff => 00:00 as login
		timeStart = tl.ParseTime("01.01.1970-00:00")
		dateStart = dateEnd
	}

	// date sanity check
	if dateStart.Unix() < 0 || dateEnd.Unix() < 0 {
		log.WithFields(log.Fields{
			"session": (*fileLineSegments)[1],
			"date":    search.login[3] + "," + (*fileLineSegments)[3],
		}).Warn("api - data: session invalid date")
		search.active = false
		return
	}

	// time sanity check
	if timeStart.Unix() < 0 || timeEnd.Unix() < 0 {
		log.WithFields(log.Fields{
			"session": (*fileLineSegments)[1],
			"time":    search.login[4] + "," + (*fileLineSegments)[4],
		}).Warn("api - data: session invalid time")
		search.active = false
		return
	}

	// selected month check
	if dateStart.Year() != src.date.Year() || dateStart.Month() != src.date.Month() || dateEnd.Year() != src.date.Year() || dateEnd.Month() != src.date.Month() {
		search.active = false
		return
	}

	search.logoff = *fileLineSegments

	// over midnight check
	if dateEnd.Unix() > dateStart.Unix() {
		search.sessionList = append(search.sessionList, DataItemSession{Date: dateStart.Unix() / 60 / 60 / 24, Detail: src.getDetail(search.login[1], search.login[2]), Time: []int{int(timeStart.Unix() / 60), 1440}})        // start to midnight
		search.sessionList = append(search.sessionList, DataItemSession{Date: dateEnd.Unix() / 60 / 60 / 24, Detail: src.getDetail((*fileLineSegments)[1], (*fileLineSegments)[2]), Time: []int{0, int(timeEnd.Unix() / 60)}}) // midnight to end

		src.addDate(&search.login[3], dateStart.Unix()/60/60/24)
		src.addDate(&(*fileLineSegments)[3], dateEnd.Unix()/60/60/24)

		search.active = false
		return
	}

	search.sessionList = append(search.sessionList, DataItemSession{Date: dateStart.Unix() / 60 / 60 / 24, Detail: src.getDetail((*fileLineSegments)[1], (*fileLineSegments)[2]), Time: []int{int(timeStart.Unix() / 60), int(timeEnd.Unix() / 60)}})
	src.addDate(&search.login[3], dateStart.Unix()/60/60/24)

	search.active = false
}

func (src *dataSource) addCurrentSession(search *dataSearch) {
	if search.active && time.Now().Format("02.01.2006") == search.login[3] {
		dateStart := tl.ParseTime(search.login[3] + "-00:00")

		timeStart := tl.ParseTime("01.01.1970-" + search.login[4])
		timeEnd := tl.ParseTime("01.01.1970-" + time.Now().Format("15:04"))

		// time & date sanity check
		if timeStart.Unix() < 0 || dateStart.Unix() < 0 {
			log.WithFields(log.Fields{
				"session": search.login[1],
				"date":    search.login[3],
				"time":    search.login[4],
			}).Warn("api - data: session invalid date or time")
		} else {
			// selected month check
			if dateStart.Year() == src.date.Year() && dateStart.Month() == src.date.Month() {
				search.sessionList = append(search.sessionList, DataItemSession{Date: dateStart.Unix() / 60 / 60 / 24, Detail: src.getDetail(search.login[1], search.login[2]), Time: []int{int(timeStart.Unix() / 60), int(timeEnd.Unix() / 60)}})
				src.addDate(&search.login[3], dateStart.Unix()/60/60/24)
			}
		}
	}
}

// addDate (thread safe)
func (src *dataSource) addDate(dateString *string, date int64) {
	src.mu.Lock()

	// set date only if there isn't already value for that date
	if _, e := src.dateList[*dateString]; !e {
		src.dateList[*dateString] = date
	}

	src.mu.Unlock()
}

// addDate (thread safe)
func (src *dataSource) addItem(item DataItem) {
	src.mu.Lock()

	src.itemList = append(src.itemList, item)

	src.mu.Unlock()
}

// addLines (thread safe)
func (src *dataSource) addLines(lines int) {
	src.mu.Lock()

	src.lines = src.lines + lines

	src.mu.Unlock()
}

// dataDetail: return proper description based on item type
// pos1: user/device
// pos2: device/user
func (src *dataSource) getDetail(pos1 string, pos2 string) string {
	switch src.tp {
	case "user":
		return pos2
	case "device":

		split := strings.Split(pos1, " ")
		if len(split) < 2 {
			return pos2
		}

		return pos2 + " " + split[1]
	default:
		return "unknown"
	}
}
