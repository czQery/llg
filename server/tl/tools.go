package tl

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

var Build = "dev"

func Log(prefix string, message string, messageType string) {
	if messageType == "debug" && Build != "dev" {
		return
	}

	now := time.Now()
	messageType = "[" + messageType + "]"
	fmt.Println("[" + now.Format("02/01/2006") + " - " + now.Format("15:04:05") + "] " + messageType + " " + prefix + " - " + message)
}

var Config map[string]interface{}

func LoadConfig() {
	configFile, err1 := os.ReadFile("config.json")
	err2 := json.Unmarshal(configFile, &Config)
	if err1 != nil || err2 != nil {
		Log("config", "load error: "+err1.Error()+" & "+err2.Error(), "error")
		os.Exit(1)
	}

	Log("config", "successfully loaded!", "info")
}

var Dist bool

func LoadDist() {
	_, err := os.ReadDir("./dist")
	if err != nil {
		Dist = false
		Log("dist", "load failed: "+err.Error(), "warn")
		return
	}

	Dist = true
	Log("dist", "successfully loaded!", "info")
}
