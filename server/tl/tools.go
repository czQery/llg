package tl

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"os"
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

func CleanSlice(slice *[]string) {
	var tmp []string

	for _, d := range *slice {
		if d != "" {
			tmp = append(tmp, d)
		}
	}

	*slice = tmp
}
