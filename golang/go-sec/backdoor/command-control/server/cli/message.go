package cli

import (
	"github.com/fatih/color"
	"time"
)

func message(level string, message string) {
	switch level {
	case "info":
		color.Cyan("[*]" + message)
	case "note":
		color.Yellow("[-]" + message)
	case "warn":
		color.Red("[!]" + message)
	case "debug":
		color.Green("[Debug]" + message)
	case "success":
		color.Green("[+]" + message)
	default:
		color.Red("[_-_]Invalid message level: " + message)
	}
}

func ListCmdResult() {
	beat := time.Tick(time.Second * 2)
	for range beat {
		DisplayCmdResult()
	}
}