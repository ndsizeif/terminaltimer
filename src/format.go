package main

import (
	"fmt"
	"os"
	"os/exec"
)

const (
	clearLine         = "\033[2K"
	carriageReturn    = "\r"
	cursorHide        = "\x1B[?25l"
	cursorShow        = "\x1B[?25h"
	cursorPrevLine    = "\033[1A"
	fileDescriptor    = "/proc/self/fd/0"
	altFileDescriptor = "/dev/fd/0"

	messageBreak = "time for a break"
	messageWork  = "time to work"
	messageDone  = "time complete"
)

// run an stty command using the constant fileDescriptor path or alternate path
func Stty(cmd string) error {
	var fileFlag string // bsd uses -F, linux -f

	tty, err := os.Readlink(fileDescriptor)
	fileFlag = "-F"
	if err != nil {
		tty, err = os.Readlink(altFileDescriptor)
		fileFlag = "-f"
		if err != nil {
			return err
		}
	}
	err = exec.Command("stty", fileFlag, tty, cmd).Run()
	if err != nil {
		return err
	}
	return nil
}

func FormatPercent(i int) (result string) {
	result = fmt.Sprintf("%4d", i)
	return
}

var stateString = map[string]string{
	"on":      "running",
	"paused":  "paused",
	"stopped": "stopped",
	"expired": "complete",
	"break":   "break",
	"breakp":  "break paused",
	"notify":  "notify",
	"tmux":    "tmux",
}

// 󱎫 󱫌 󱫔 󱎬 󱫒 󰀦 󰙚 󱫞 󰀄 󱅞 󰍩 
var icon_solid = map[string]string{
	"on":      "󱎫",
	"warning": "󱫌",
	"paused":  "󱫔",
	"stopped": "󱎬",
	"edit":    "󱫒",
	"expired": "󰀦",
	"break":   "󰀄",
	"breakp":  "󱅞",
	"notify":  "󰍩",
	"tmux":    "",
	"restart": "󰜉",
}

// 󰔛 󱫍 󱫗 󰔞 󱫓 󰀪 󰌦 󱫟 󰀓 󱅟 󰍪 
var icon_trace = map[string]string{
	"on":      "󰔛",
	"warning": "󱫍",
	"paused":  "󱫗",
	"stopped": "󰔞",
	"edit":    "󱫓",
	"expired": "󰀪",
	"break":   "󰀓",
	"breakp":  "󱅟",
	"notify":  "󰍪",
	"tmux":    "",
	"restart": "󰜉",
}
var icon_ascii = map[string]string{
	"on":      ">",
	"warning": "!",
	"paused":  "~",
	"stopped": ":",
	"edit":    "^",
	"expired": "-",
	"break":   "+",
	"breakp":  "+",
	"notify":  "*",
	"tmux":    "t",
	"restart": "r",
}
var bar_solid = map[string]string{
	"done": "█",
	"todo": "░",
	"stop": " ",
}
var bar_solid_rev = map[string]string{
	"done": "░",
	"todo": "█",
	"stop": " ",
}
var bar_shade = map[string]string{
	"done": "▒",
	"todo": "░",
	"stop": " ",
}
var bar_shade_rev = map[string]string{
	"done": "░",
	"todo": "▒",
	"stop": " ",
}
var bar_ascii = map[string]string{
	"done": "@",
	"todo": "-",
	"stop": " ",
}
