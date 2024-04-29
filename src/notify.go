package main

import (
	"os/exec"
)

func NotifySend(message string) error {
	notifyCmd := "notify-send"
	notifyApp := "--app-name=" + programName
	notifyIcon := "--icon=clock"

	cmd := exec.Command(notifyCmd, message, notifyApp, notifyIcon)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
