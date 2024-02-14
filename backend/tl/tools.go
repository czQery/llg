package tl

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

var Build = "dev"

var Config map[string]interface{}

func LoadConfig() {
	configFile, err1 := os.ReadFile("config.json")
	err2 := json.Unmarshal(configFile, &Config)
	if err1 != nil || err2 != nil {
		log.WithFields(log.Fields{
			"read_error": err1,
			"json_error": err2,
		}).Panic("config - load failed")
	}

	log.Info("config - successfully loaded")
}

var Dist bool

func LoadDist() {
	_, err := os.ReadDir("./dist")
	if err != nil {
		Dist = false
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Warn("dist - load failed")
		return
	}

	Dist = true
	log.Info("dist - successfully loaded")
}

// CleanSlice removes empty items
func CleanSlice(slice *[]string) {
	var tmp []string

	for _, d := range *slice {
		if d != "" {
			tmp = append(tmp, d)
		}
	}

	*slice = tmp
}

func ParseTime(tm string) time.Time {

	// for time format without extra 0. for example 01.02.2024 -> 1.2.2024
	offset := 0

	day := (int(tm[0])-'0')*10 + int(tm[1]) - '0'
	if tm[1] == '.' {
		day = int(tm[0]) - '0'
		offset = 1
	}

	month := time.Month((int(tm[3-offset])-'0')*10 + int(tm[4-offset]) - '0')
	if tm[4-offset] == '.' {
		month = time.Month(int(tm[3-offset]) - '0')
		offset = offset + 1
	}

	year := (((int(tm[6-offset])-'0')*10+int(tm[7-offset])-'0')*10+int(tm[8-offset])-'0')*10 + int(tm[9-offset]) - '0'

	hour := (int(tm[11-offset])-'0')*10 + int(tm[12-offset]) - '0'
	if tm[12-offset] == ':' {
		hour = int(tm[11-offset]) - '0'
		offset = offset + 1
	}

	minute := (int(tm[14-offset])-'0')*10 + int(tm[15-offset]) - '0'
	if tm[12-offset] == ':' {
		minute = int(tm[14-offset]) - '0'
	}

	return time.Date(year, month, day, hour, minute, 0, 0, time.UTC)
}
