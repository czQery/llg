package tl

import (
	"fmt"
	"time"
)

func Log(prefix string, message string, messageType string) {
	now := time.Now()

	messageType = "[" + messageType + "]"

	fmt.Println("[" + now.Format("02/01/2006") + " - " + now.Format("15:04:05") + "] " + messageType + " " + prefix + " - " + message)
}
