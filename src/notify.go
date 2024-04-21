package main

import (
	"os/exec"
)

var (
	notifyCmd    = "notify-send"
	notifyApp    = "--app-name=" + programName
	notifyIcon   = "--icon=clock"
)

func NotifySend(message string) error {
	cmd := exec.Command("notify-send", message, notifyApp, notifyIcon)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
