package tl

import (
	"fmt"
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
