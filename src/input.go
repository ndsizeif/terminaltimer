package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

var (
	programName  = "terminalTimer"
	zeroDuration = time.Duration(0 * time.Minute)
	programHelp  bool

	setTimer time.Duration
	setBreak time.Duration
	setAlert time.Duration
)

func ValidateFlags() {
	args := len(os.Args)
	programCmd := flag.NewFlagSet(programName, flag.ExitOnError)
	programCmd.BoolVar(&programHelp, "help", false, UsageString["help"])
	programCmd.BoolVar(&programHelp, "h", false, UsageString["help"]+" (shorthand)")

	taskCmd := flag.NewFlagSet(programName+" task", flag.ExitOnError)
	taskString := taskCmd.String("task", "", UsageString["task"])

	setCmd := flag.NewFlagSet(programName+" set", flag.ExitOnError)
	setCmd.DurationVar(&setTimer, "timer", zeroDuration, UsageString["setTimer"])
	setCmd.DurationVar(&setTimer, "t", zeroDuration, UsageString["setTimer"]+" (shorthand)")
	setCmd.DurationVar(&setBreak, "break", zeroDuration, UsageString["setBreak"])
	setCmd.DurationVar(&setBreak, "b", zeroDuration, UsageString["setBreak"]+" (shorthand)")
	setCmd.DurationVar(&setAlert, "alert", zeroDuration, UsageString["setAlert"])
	setCmd.DurationVar(&setAlert, "a", zeroDuration, UsageString["setAlert"]+" (shorthand)")

	styleCmd := flag.NewFlagSet(programName+" style", flag.ExitOnError)
	styleWidth := styleCmd.Int("width", 0, UsageString["styleWidth"])
	styleBar := styleCmd.Int("bar", 0, UsageString["styleBar"])
	styleIcon := styleCmd.Int("icon", 0, UsageString["styleIcon"])

	toggleCmd := flag.NewFlagSet(programName+" toggle", flag.ExitOnError)
	toggleBar := toggleCmd.Bool("bar", false, UsageString["toggleBar"])
	toggleBell := toggleCmd.Bool("bell", false, UsageString["toggleBell"])
	toggleClock := toggleCmd.Bool("clock", false, UsageString["toggleClock"])
	toggleIcon := toggleCmd.Bool("icon", false, UsageString["toggleIcon"])
	togglePercent := toggleCmd.Bool("percent", false, UsageString["togglePercent"])
	toggleNotify := toggleCmd.Bool("notify", false, UsageString["toggleNotify"])
	toggleTmux := toggleCmd.Bool("tmux", false, UsageString["toggleTmux"])
	toggleRestart := toggleCmd.Bool("restart", false, UsageString["toggleRestart"])
	toggleReverse := toggleCmd.Bool("reverse", false, UsageString["toggleReverse"])

	// if program is invoked with one argument, and it is --help
	if args == 2 {
		programCmd.Parse(os.Args[1:])
		if programHelp {
			PrintUsage()
			fmt.Println()
			setCmd.Usage()
			fmt.Println()
			styleCmd.Usage()
			fmt.Println()
			toggleCmd.Usage()
			os.Exit(0)
		}
	}
	// invoking with program name by itself will render static output and exit
	if args < 2 {
		t, _ := InitializeTimer()
		t.GetTime()
		t.Render()
		os.Exit(0)
	}
	// check if first argument is a non-flag, and matches one of the listed strings
	switch os.Args[1] {
	case "clear":
		HandleProgramCmd("clear")
	case "task":
		HandleTaskCmd(taskCmd, taskString)
	case "start":
		HandleProgramCmd("start")
	case "stop":
		HandleProgramCmd("stop")
	case "pause":
		HandleProgramCmd("pause")
	case "resume":
		HandleProgramCmd("resume")
	case "break":
		HandleProgramCmd("break")
	case "run":
		HandleProgramCmd("run")
		HandleProgramRun()
	case "info":
		HandleInfo()
	case "status":
		HandleStatus()
	case "clean":
		HandleClean()
	case "set":
		HandleSetCmd(setCmd, setTimer, setBreak, setAlert)
	case "style":
		HandleStyleCmd(styleCmd, styleWidth, styleBar, styleIcon)
	case "toggle":
		HandleToggleCmd(toggleCmd, toggleBar, toggleBell, toggleClock, toggleIcon,
			toggleNotify, togglePercent, toggleRestart, toggleReverse, toggleTmux)
	case "help":
		PrintUsage()
		fmt.Println()
		setCmd.Usage()
		fmt.Println()
		styleCmd.Usage()
		fmt.Println()
		toggleCmd.Usage()
	default:
		PrintUsage()
	}
	os.Exit(0)
}

func ValidateTaskCmd(taskCmd *flag.FlagSet, task *string) {
	taskCmd.Parse(os.Args[2:])
	if len(taskCmd.Args()) == 0 {
		taskCmd.Usage()
		os.Exit(0)
	}
}

func HandleTaskCmd(taskCmd *flag.FlagSet, taskStr *string) {
	ValidateTaskCmd(taskCmd, taskStr)
	t, _ := InitializeTimer()
	cmd := t.Display["task"]
	str := strings.Join(taskCmd.Args(), " ")
	cmd(str)
	t.Message("changed task to")
	t.State.Save()
}

func HandleProgramCmd(command string) {
	t, _ := InitializeTimer()
	cmd := t.ExecuteCommand(command)
	if cmd == nil {
		t.State.Debug.Print("invalid command:", command)
		return
	}
	cmd()
}

func HandleProgramRun() {
	t, _ := InitializeTimer()
	cmd := t.ExecuteCommand("run")
	cmd()
}

func HandleInfo() {
	t, _ := InitializeTimer()
	c := t.Info()
	fmt.Printf(c)
}

func HandleStatus() {
	t, _ := InitializeTimer()
	c := t.GetState()
	fmt.Printf(c)
}

func HandleClean() {
	err := RemoveLogFiles()
	if err != nil && !os.IsNotExist(err) { // exclude "can't be found" error
		fmt.Printf("%v clean error: %v\n", programName, err)
		return
	}
	fmt.Printf("%v log files removed\n", programName)
}

// handle negative durations being passed
func ValidateSetCmd(setCmd *flag.FlagSet, timeInterval, breakInterval, alertInterval time.Duration) {
	setCmd.Parse(os.Args[2:])
	if len(os.Args) < 3 {
		setCmd.Usage()
		os.Exit(0)
	}
	if timeInterval < zeroDuration || breakInterval < zeroDuration || alertInterval < zeroDuration {
		setCmd.Usage()
		os.Exit(0)
	}
	if len(setCmd.Args()) > 0 {
		args := setCmd.Args()
		fmt.Printf("'%s' invalid: flags accept one duration value\n", strings.Join(args[0:], " "))
		setCmd.Usage()
		os.Exit(0)
	}
}
func HandleSetCmd(setCmd *flag.FlagSet, timeInterval, breakInterval, alertInterval time.Duration) {
	setCmd.Parse(os.Args[2:])
	ValidateSetCmd(setCmd, timeInterval, breakInterval, alertInterval)

	t, _ := InitializeTimer()
	if timeInterval > zeroDuration {
		cmd := t.SetDuration("timer")
		cmd(timeInterval)
	}
	if breakInterval > zeroDuration {
		cmd := t.SetDuration("break")
		cmd(breakInterval)
	}
	if alertInterval > zeroDuration {
		cmd := t.SetDuration("alert")
		cmd(alertInterval)
	}
	t.State.Save()
}

func ValidateStyleCmd(styleCmd *flag.FlagSet, width, bar, icon *int) {
	if len(os.Args) < 3 {
		styleCmd.Usage()
		os.Exit(0)
	}
	if *width < 0 || *bar < 0 || *icon < 0 {
		styleCmd.Usage()
		os.Exit(0)
	}
	if *width < 0 {
		*width = 0
	}
	if *width > 300 {
		*width = 300
	}
}

func HandleStyleCmd(styleCmd *flag.FlagSet, width, bar, icon *int) {
	styleCmd.Parse(os.Args[2:])
	ValidateStyleCmd(styleCmd, width, bar, icon)

	t, _ := InitializeTimer()
	if *width != 0 {
		cmd := t.SetOption("size")
		cmd(*width)
	}
	if *bar > 0 {
		cmd := t.SetOption("style")
		cmd(*bar)
	}
	if *icon > 0 {
		cmd := t.SetOption("symbol")
		cmd(*icon)
	}
	err := t.Config.Save()
	if err != nil {
		t.State.Debug.Print(t.State.Debug.Trace(), err)
	}
}

func HandleToggleCmd(toggleCmd *flag.FlagSet, bar, bell, clock, icon, notify, percent, restart, reverse, tmux *bool) {
	toggleCmd.Parse(os.Args[2:])
	var cmd []func(bool)
	var config []bool
	t, _ := InitializeTimer()
	if *bar {
		cmd = append(cmd, t.ToggleOption("bar"))
		config = append(config, t.Config.HideBar)
	}
	if *bell {
		cmd = append(cmd, t.ToggleOption("bell"))
		config = append(config, t.Config.Bell)
	}
	if *clock {
		cmd = append(cmd, t.ToggleOption("clock"))
		config = append(config, t.Config.HideTime)
	}
	if *icon {
		cmd = append(cmd, t.ToggleOption("icon"))
		config = append(config, t.Config.HideIcon)
	}
	if *notify {
		cmd = append(cmd, t.ToggleOption("notify"))
		config = append(config, t.Config.Notify)
	}
	if *percent {
		cmd = append(cmd, t.ToggleOption("percent"))
		config = append(config, t.Config.Percent)
	}
	if *restart {
		cmd = append(cmd, t.ToggleOption("restart"))
		config = append(config, t.Config.Restart)
	}
	if *reverse {
		cmd = append(cmd, t.ToggleOption("reverse"))
		config = append(config, t.Config.ReverseTime)
	}
	if *tmux {
		cmd = append(cmd, t.ToggleOption("tmux"))
		config = append(config, t.Config.NotifyTmux)
	}
	if len(cmd) == 0 {
		toggleCmd.Usage()
		os.Exit(0)
	}
	for k, function := range cmd { // call function with the current config value flipped
		function(!config[k])
	}
	err := t.Config.Save()
	if err != nil {
		t.State.Debug.Print(t.State.Debug.Trace(), err)
		os.Exit(1)
	}
}

var UsageString = map[string]string{
	"help":          "display full help",
	"start":         "start timer",
	"stop":          "stop timer",
	"pause":         "pause timer",
	"resume":        "start timer if paused",
	"break":         "start break",
	"run":           "display timer inline inside terminal",
	"clean":         "delete timer log file",
	"clear":         "clear the string for current task",
	"status":        "return current timer status",
	"setTimer":      "set timer interval",
	"setBreak":      "set break interval",
	"setAlert":      "set threshold to start alert",
	"styleWidth":    "style progress bar width",
	"styleBar":      "style progress bar appearance",
	"styleIcon":     "style icon appearance",
	"task":          "set the string for current task",
	"toggleBar":     "turn progress bar on/off",
	"toggleBell":    "turn terminal bell on/off",
	"toggleClock":   "display time on/off",
	"toggleIcon":    "display icons on/off",
	"toggleNotify":  "turn timer notifications on/off for notify-send",
	"togglePercent": "display interval percentage on/off",
	"toggleRestart": "turn automatic timer restart on/off",
	"toggleReverse": "timer displays time descending/ascending",
	"toggleTmux":    "turn timer notifications on/off for tmux",
}

func PrintUsage() {
	fmt.Printf("Usage of %v:\n", programName)
	fmt.Printf("  start\n\t%v\n", UsageString["start"])
	fmt.Printf("  stop\n\t%v\n", UsageString["stop"])
	fmt.Printf("  pause\n\t%v\n", UsageString["pause"])
	fmt.Printf("  resume\n\t%v\n", UsageString["resume"])
	fmt.Printf("  break\n\t%v\n", UsageString["break"])
	fmt.Printf("  run\n\t%v\n", UsageString["run"])
	fmt.Printf("  clean\n\t%v\n", UsageString["clean"])
	fmt.Printf("  clear\n\t%v\n", UsageString["clear"])
	fmt.Printf("  help\n\t%v\n", UsageString["help"])
}
