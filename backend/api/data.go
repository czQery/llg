package api

import (
	"bufio"
	"bytes"
	"maps"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/czQery/llg/backend/tl"
	"github.com/gofiber/fiber/v3"
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
	active      bool // searching for logoff
	name        string
	login       dataLine
	logoff      dataLine
	sessionList []DataItemSession
	lines       int
}

type dataLine struct {
	action string
	name   string
	detail string
	date   string
	time   string
}

const timeLayout = "02.01.2006-15:04"
const timeLayoutParam = "2006-01"

func Data(c fiber.Ctx) error {
	var (
		dataDates          = make([]int64, 0, 30)
		searchDateList     = make(map[string]int64)
		searchDateListKeys = make([]string, 0, 30)

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
		"devices": strings.Join(devicesParam, ","),
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
		fileOpen *os.File
		scanner  *bufio.Scanner

		lineRaw      []byte
		lineRawSplit [5][]byte
		line         = dataLine{
			action: "",
			name:   "",
			detail: "",
			date:   "",
			time:   "",
		}

		lineDate time.Time

		search = dataSearch{
			active:      false,
			name:        "",
			login:       dataLine{},
			logoff:      dataLine{},
			sessionList: make([]DataItemSession, 0, 30),
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
	fileOpen, err = os.Open(src.folder + (*file).Name())
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("api - data: open")
		return
	}
	defer fileOpen.Close()

	// read all lines in file
	scanner = bufio.NewScanner(fileOpen)
	for scanner.Scan() {
		lineRaw = scanner.Bytes()
		search.lines = search.lines + 1

		if !getSegment(&(lineRawSplit[0]), &lineRaw) ||
			!getSegment(&(lineRawSplit[1]), &lineRaw) ||
			!getSegment(&(lineRawSplit[2]), &lineRaw) ||
			!getSegment(&(lineRawSplit[3]), &lineRaw) {
			continue
		}
		lineRawSplit[4] = lineRaw

		// skip lines outside of selected time range
		line.date = string(lineRawSplit[3])
		lineDate = tl.ParseTime(line.date + "-00:00")
		if lineDate.Year() != src.date.Year() || lineDate.Month() != src.date.Month() {
			continue
		}

		// cast the rest to string
		line.action = string(lineRawSplit[0])
		line.name = string(lineRawSplit[1])
		line.detail = string(lineRawSplit[2])
		line.time = string(lineRawSplit[4])

		src.readFileLine(&search, &line)
	}

	if err = scanner.Err(); err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("api - data: scanner")
		return
	}

	// add line count
	src.addLines(search.lines)

	// add currently ongoing session || fix fill last session in log to midnight
	src.addLastSession(&search)

	// bundle sessions if found
	if len(search.sessionList) > 0 {
		if src.tp == "device" {
			search.name = strings.Split(search.name, " ")[0]
		}

		src.addItem(DataItem{Name: search.name, Type: src.tp, Sessions: search.sessionList})
	}

	return
}

func (src *dataSource) readFileLine(search *dataSearch, fileLine *dataLine) {
	switch fileLine.action {
	case "login":
		// mark session start
		src.readFileLineLogin(search, fileLine)
	case "logoff":
		// mark session end
		src.readFileLineLogoff(search, fileLine)
	}
}

func (src *dataSource) readFileLineLogin(search *dataSearch, fileLine *dataLine) {
	// fill the prevous active session to midnight if this login is from a different day than the previous login
	if search.active && search.login.date != fileLine.date {

		timeStart := tl.ParseTime("01.01.1970-" + search.login.time)
		dateStart := tl.ParseTime(search.login.date + "-00:00")

		// time & date sanity check
		if timeStart.Unix() < 0 || dateStart.Unix() < 0 {
			log.WithFields(log.Fields{
				"session": search.login.name,
				"date":    search.login.date,
				"time":    search.login.time,
			}).Warn("api - data: session invalid date or time")
		} else {
			search.sessionList = append(search.sessionList, DataItemSession{Date: dateStart.Unix() / 60 / 60 / 24, Detail: src.getDetail(search.login.name, search.login.detail), Time: []int{int(timeStart.Unix() / 60), 1440}}) // start to midnight
			src.addDate(&search.login.date, dateStart.Unix()/60/60/24)
		}

		search.active = false
	}

	if !search.active {
		search.active = true
		search.name = fileLine.name
		search.login = *fileLine
		return
	}
}

func (src *dataSource) readFileLineLogoff(search *dataSearch, fileLine *dataLine) {
	var (
		timeStart time.Time
		timeEnd   time.Time

		dateStart time.Time
		dateEnd   time.Time
	)

	// hack for logs starting with logoff
	if search.login.action == "" {
		search.login = *fileLine
	}

	timeEnd = tl.ParseTime("01.01.1970-" + fileLine.time)
	dateEnd = tl.ParseTime(fileLine.date + "-00:00")

	timeStart = tl.ParseTime("01.01.1970-" + search.login.time)
	dateStart = tl.ParseTime(search.login.date + "-00:00")

	if !search.active { // missing login
		if search.logoff.action != "" && search.logoff.date == fileLine.date { // previous logoff in same day => filling session from the previous logoff
			timeStart = tl.ParseTime("01.01.1970-" + search.logoff.time)
			dateStart = tl.ParseTime(search.logoff.date + "-00:00")
		} else { // first logoff => filling session from 00:00
			timeStart = tl.ParseTime("01.01.1970-00:00")
			dateStart = dateEnd
		}

		// date & time sanity check
		if timeStart.Unix() < 0 || dateStart.Unix() < 0 || timeEnd.Unix() < 0 {
			log.WithFields(log.Fields{
				"session": fileLine.name,
				"date":    fileLine.date,
				"time":    fileLine.time,
			}).Warn("api - data: session invalid date or time")
		} else {
			search.sessionList = append(search.sessionList, DataItemSession{Date: dateStart.Unix() / 60 / 60 / 24, Detail: src.getDetail(fileLine.name, fileLine.detail), Time: []int{int(timeStart.Unix() / 60), int(timeEnd.Unix() / 60)}})
			src.addDate(&fileLine.date, dateStart.Unix()/60/60/24)
		}

		search.logoff = *fileLine
		return
	}

	// date sanity check
	if dateStart.Unix() < 0 || dateEnd.Unix() < 0 {
		log.WithFields(log.Fields{
			"session": fileLine.name,
			"date":    search.login.date + "," + fileLine.date,
		}).Warn("api - data: session invalid date")
		return
	}

	// time sanity check
	if timeStart.Unix() < 0 || timeEnd.Unix() < 0 {
		log.WithFields(log.Fields{
			"session": fileLine.name,
			"time":    search.login.time + "," + fileLine.time,
		}).Warn("api - data: session invalid time")
		return
	}

	search.active = false
	search.logoff = *fileLine

	// over midnight check
	if dateEnd.Unix() > dateStart.Unix() {
		search.sessionList = append(search.sessionList, DataItemSession{Date: dateStart.Unix() / 60 / 60 / 24, Detail: src.getDetail(search.login.name, search.login.detail), Time: []int{int(timeStart.Unix() / 60), 1440}}) // start to midnight
		search.sessionList = append(search.sessionList, DataItemSession{Date: dateEnd.Unix() / 60 / 60 / 24, Detail: src.getDetail(fileLine.name, fileLine.detail), Time: []int{0, int(timeEnd.Unix() / 60)}})                // midnight to end

		src.addDate(&search.login.date, dateStart.Unix()/60/60/24)
		src.addDate(&fileLine.date, dateEnd.Unix()/60/60/24)
		return
	}

	search.sessionList = append(search.sessionList, DataItemSession{Date: dateStart.Unix() / 60 / 60 / 24, Detail: src.getDetail(fileLine.name, fileLine.detail), Time: []int{int(timeStart.Unix() / 60), int(timeEnd.Unix() / 60)}})
	src.addDate(&search.login.date, dateStart.Unix()/60/60/24)
}

func (src *dataSource) addLastSession(search *dataSearch) {
	if !search.active {
		return
	}

	if time.Now().Format("02.01.2006") == search.login.date {
		dateStart := tl.ParseTime(search.login.date + "-00:00")

		timeStart := tl.ParseTime("01.01.1970-" + search.login.time)
		timeEnd := tl.ParseTime("01.01.1970-" + time.Now().Format("15:04"))

		// time & date sanity check
		if timeStart.Unix() < 0 || dateStart.Unix() < 0 || timeEnd.Unix() < 0 {
			log.WithFields(log.Fields{
				"session": search.login.name,
				"date":    search.login.date,
				"time":    search.login.time,
			}).Warn("api - data: session invalid date or time")
		} else {
			search.sessionList = append(search.sessionList, DataItemSession{Date: dateStart.Unix() / 60 / 60 / 24, Detail: src.getDetail(search.login.name, search.login.detail), Time: []int{int(timeStart.Unix() / 60), int(timeEnd.Unix() / 60)}}) // start to current time
			src.addDate(&search.login.date, dateStart.Unix()/60/60/24)
		}
	} else {
		timeStart := tl.ParseTime("01.01.1970-" + search.login.time)
		dateStart := tl.ParseTime(search.login.date + "-00:00")

		// time & date sanity check
		if timeStart.Unix() < 0 || dateStart.Unix() < 0 {
			log.WithFields(log.Fields{
				"session": search.login.name,
				"date":    search.login.date,
				"time":    search.login.time,
			}).Warn("api - data: session invalid date or time")
		} else {
			search.sessionList = append(search.sessionList, DataItemSession{Date: dateStart.Unix() / 60 / 60 / 24, Detail: src.getDetail(search.login.name, search.login.detail), Time: []int{int(timeStart.Unix() / 60), 1440}}) // start to midnight
			src.addDate(&search.login.date, dateStart.Unix()/60/60/24)
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

// getSegment: gets value from byte array separated by semicolon
func getSegment(dst *[]byte, raw *[]byte) bool {
	idx := bytes.IndexByte(*raw, ';')
	if idx == -1 {
		return false
	}
	*dst = (*raw)[:idx]
	*raw = (*raw)[idx+1:]
	return true
}
