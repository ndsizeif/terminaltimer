package tmuxmenu

import (
	"fmt"
	"os/exec"
	"errors"
)

var format = []string{
	"#[align=left]",
	"#[align=right]",
	"#[align=centre]",
	"#[align=absolute-centre]",
}

var height = []string {
	"0",
	"#{window_height}",
}

var width = []string{
	"0",
	"#{window_width}",
}

type Menu struct {
	Cmd string
	OpenCmd string
	CloseCmd []string
	Padding []string
	Separator []string
	Message []string
	Format string
	X string
	Y string
}

func Initialize() *Menu{
	return &Menu{
		Cmd: "tmux",
		OpenCmd: "display-menu",
		CloseCmd: []string{"display-popup", "-C"},
		Padding: []string{"-#[nodim]", "", ""},
		Separator: []string{"", "", ""},
		Message: []string{"test message", "", ""},
		Format: format[3],
		X: width[1],
		Y: height[1],
	}
}

// NOTE when Tmux is running with sessions: output will be > 0 err will be nil
// when Tmux is not running: output will be 0 error will be !nil
// Available == len(output) > 0 && err == nil

// check if tmux is available before issuing commands
func (m *Menu) Available() bool {
	out, err := exec.Command("tmux", "list-sessions").Output()
	return len(out) > 0 && err == nil
}
// please do not close popup if the user is not expecting it
func (m *Menu) Close() error {
	if !m.Available() {
		return errors.New("tmux instance not available")
	}
	err := CommandRun(m.Cmd, m.CloseCmd)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
// generate string array needed for menu command, minus the actual message
func (m *Menu) Build(args ...string) ([]string) {
	var menu []string
	for _, v := range args {
		menu = append (menu, v)
	}
	return menu
}

func (m *Menu) Open(title, message string) error {
	if !m.Available() {
		return errors.New("tmux instance not available")
	}
	cmd := m.Build(m.OpenCmd, "-T", title, "-x", m.X, "-y", m.Y)
	formattedMessage := fmt.Sprintf("%v%v", m.Format, message)
	msg := append (cmd, m.Separator...)
	msg = append (msg, m.Padding...)
	msg = append (msg, m.Padding...)
	msg = append (msg, []string{formattedMessage, "", ""}...)
	msg = append (msg, m.Padding...)
	msg = append (msg, m.Padding...)
	msg = append (msg, m.Separator...)

	err := CommandRun(m.Cmd, msg)
	if err != nil {
		fmt.Println("menu error")
		return err
	}
	return nil
}

func CommandRun(cmd string, args []string) (err error){
	c := exec.Command(cmd, args...)
	err = c.Run()
	if err != nil {
		return err
	}
	return
}
