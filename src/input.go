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

	styleWidth int
	styleBar   int
	styleIcon  int

	toggleProgress bool
	toggleBell     bool
	toggleClock    bool
	toggleSymbol   bool
	togglePercent  bool
	toggleNotify   bool
	toggleTmux     bool
	toggleRestart  bool
	toggleReverse  bool
)

func ValidateFlags() {
	args := len(os.Args)
	programCmd := flag.NewFlagSet(programName, flag.ExitOnError)
	programCmd.BoolVar(&programHelp, "help", false, UsageString["help"])
	programCmd.BoolVar(&programHelp, "h", false, UsageString["help"])

	programCmd.Usage = func() {
		writer := flag.CommandLine.Output()
		fmt.Fprintf(writer, "%s\n\r", UsageString["programCmd"])
		f := programCmd.Lookup("help")
		fmt.Printf("  -%v, -%v\n", Shorthand[f.Name], f.Name)
		fmt.Printf("\n")
	}

	taskCmd := flag.NewFlagSet(programName+" task", flag.ExitOnError)
	taskString := taskCmd.String("task", "", UsageString["task"])

	setCmd := flag.NewFlagSet(programName+" set", flag.ExitOnError)
	setCmd.DurationVar(&setTimer, "timer", zeroDuration, UsageString["setTimer"])
	setCmd.DurationVar(&setTimer, "t", zeroDuration, UsageString["setTimer"])
	setCmd.DurationVar(&setBreak, "break", zeroDuration, UsageString["setBreak"])
	setCmd.DurationVar(&setBreak, "k", zeroDuration, UsageString["setBreak"])
	setCmd.DurationVar(&setAlert, "alert", zeroDuration, UsageString["setAlert"])
	setCmd.DurationVar(&setAlert, "a", zeroDuration, UsageString["setAlert"])

	setCmd.Usage = func() {
		writer := flag.CommandLine.Output()
		fmt.Fprintf(writer, "%s\n\r", UsageString["setCmd"])
		order := []string{"timer", "break", "alert"}
		for _, name := range order {
			f := setCmd.Lookup(name)
			fmt.Printf("  -%v, -%v\n", Shorthand[f.Name], f.Name)
			fmt.Printf("\t%s\n", f.Usage)
		}
		fmt.Printf("\n")
	}

	styleCmd := flag.NewFlagSet(programName+" style", flag.ExitOnError)
	styleCmd.IntVar(&styleWidth, "width", 0, UsageString["styleWidth"])
	styleCmd.IntVar(&styleWidth, "w", 0, UsageString["styleWidth"])
	styleCmd.IntVar(&styleBar, "bar", 0, UsageString["styleBar"])
	styleCmd.IntVar(&styleBar, "b", 0, UsageString["styleBar"])
	styleCmd.IntVar(&styleIcon, "icon", 0, UsageString["styleIcon"])
	styleCmd.IntVar(&styleIcon, "i", 0, UsageString["styleIcon"])

	styleCmd.Usage = func() {
		writer := flag.CommandLine.Output()
		fmt.Fprintf(writer, "%s\n\r", UsageString["styleCmd"])
		order := []string{"width", "bar", "icon"}
		for _, name := range order {
			f := styleCmd.Lookup(name)
			fmt.Printf("  -%v, -%v\n", Shorthand[f.Name], f.Name)
			fmt.Printf("\t%s\n", f.Usage)
		}
		fmt.Printf("\n")
	}

	toggleCmd := flag.NewFlagSet(programName+" toggle", flag.ExitOnError)
	toggleCmd.BoolVar(&toggleProgress, "progress", false, UsageString["toggleProgress"])
	toggleCmd.BoolVar(&toggleProgress, "p", false, UsageString["toggleProgress"])
	toggleCmd.BoolVar(&toggleBell, "bell", false, UsageString["toggleBell"])
	toggleCmd.BoolVar(&toggleBell, "l", false, UsageString["toggleBell"])
	toggleCmd.BoolVar(&toggleClock, "clock", false, UsageString["toggleClock"])
	toggleCmd.BoolVar(&toggleClock, "c", false, UsageString["toggleClock"])
	toggleCmd.BoolVar(&toggleSymbol, "symbol", false, UsageString["toggleSymbol"])
	toggleCmd.BoolVar(&toggleSymbol, "s", false, UsageString["toggleSymbol"])
	toggleCmd.BoolVar(&togglePercent, "percent", false, UsageString["togglePercent"])
	toggleCmd.BoolVar(&togglePercent, "P", false, UsageString["togglePercent"])
	toggleCmd.BoolVar(&toggleNotify, "notify", false, UsageString["toggleNotify"])
	toggleCmd.BoolVar(&toggleNotify, "n", false, UsageString["toggleNotify"])
	toggleCmd.BoolVar(&toggleTmux, "tmux", false, UsageString["toggleTmux"])
	toggleCmd.BoolVar(&toggleTmux, "t", false, UsageString["toggleTmux"])
	toggleCmd.BoolVar(&toggleRestart, "restart", false, UsageString["toggleRestart"])
	toggleCmd.BoolVar(&toggleRestart, "r", false, UsageString["toggleRestart"])
	toggleCmd.BoolVar(&toggleReverse, "reverse", false, UsageString["toggleReverse"])
	toggleCmd.BoolVar(&toggleReverse, "v", false, UsageString["toggleReverse"])

	toggleCmd.Usage = func() {
		writer := flag.CommandLine.Output()
		fmt.Fprintf(writer, "%s\n\r", UsageString["toggleCmd"])
		order := []string{
			"progress", "bell", "clock", "symbol", "percent", "restart", "reverse",
			"notify", "tmux"}

		for _, name := range order {
			f := toggleCmd.Lookup(name)
			fmt.Printf("  -%v, -%v\n", Shorthand[f.Name], f.Name)
			fmt.Printf("\t%s\n", f.Usage)
		}
		fmt.Printf("\n")
	}

	// invoking with program name by itself will render static output and exit
	if args == 1 {
		t, _ := InitializeTimer()
		t.GetTime()
		t.Render()
		os.Exit(0)
	}

	// if program is invoked with one argument, and it is --help
	if args == 2 {
		programCmd.Parse(os.Args[1:])
		if programHelp {
			PrintBasicUsage()
			setCmd.Usage()
			styleCmd.Usage()
			toggleCmd.Usage()
			os.Exit(0)
		}
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
		// HandleProgramRun()
	case "info":
		HandleInfo()
	case "status":
		HandleStatus()
	case "clean":
		HandleClean()
	case "set":
		HandleSetCmd(setCmd, &setTimer, &setBreak, &setAlert)
	case "style":
		HandleStyleCmd(styleCmd, &styleWidth, &styleBar, &styleIcon)
	case "toggle":
		HandleToggleCmd(toggleCmd, &toggleProgress, &toggleBell, &toggleClock, &toggleSymbol,
			&toggleNotify, &togglePercent, &toggleRestart, &toggleReverse, &toggleTmux)
	case "help":
		PrintBasicUsage()
		setCmd.Usage()
		styleCmd.Usage()
		toggleCmd.Usage()
	default:
		PrintBasicUsage()
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
func ValidateSetCmd(setCmd *flag.FlagSet, timeInterval, breakInterval, alertInterval *time.Duration) {
	setCmd.Parse(os.Args[2:])
	if len(os.Args) < 3 {
		setCmd.Usage()
		os.Exit(0)
	}
	if *timeInterval < zeroDuration || *breakInterval < zeroDuration || *alertInterval < zeroDuration {
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
func HandleSetCmd(setCmd *flag.FlagSet, timeInterval, breakInterval, alertInterval *time.Duration) {
	setCmd.Parse(os.Args[2:])
	ValidateSetCmd(setCmd, timeInterval, breakInterval, alertInterval)

	t, _ := InitializeTimer()
	if *timeInterval > zeroDuration {
		cmd := t.SetDuration("timer")
		cmd(*timeInterval)
	}
	if *breakInterval > zeroDuration {
		cmd := t.SetDuration("break")
		cmd(*breakInterval)
	}
	if *alertInterval > zeroDuration {
		cmd := t.SetDuration("alert")
		cmd(*alertInterval)
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

func HandleToggleCmd(toggleCmd *flag.FlagSet, progress, bell, clock, symbol, notify, percent, restart, reverse, tmux *bool) {
	toggleCmd.Parse(os.Args[2:])
	var cmd []func(bool)
	var config []bool
	t, _ := InitializeTimer()
	if *progress {
		cmd = append(cmd, t.ToggleOption("progress"))
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
	if *symbol {
		cmd = append(cmd, t.ToggleOption("symbol"))
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
	"programCmd":     "Usage of " + programName,
	"setCmd":         "Usage of " + programName + " set (duration)",
	"styleCmd":       "Usage of " + programName + " style (int)",
	"toggleCmd":      "Usage of " + programName + " toggle",
	"help":           "display full help",
	"start":          "start timer",
	"stop":           "stop timer",
	"pause":          "pause timer",
	"resume":         "start timer if paused",
	"break":          "start break",
	"run":            "display timer inline inside terminal",
	"clean":          "delete timer log file",
	"clear":          "clear the string for current task",
	"status":         "return current timer status",
	"info":           "return current timer interval values",
	"setTimer":       "set timer interval",
	"setBreak":       "set break interval",
	"setAlert":       "set threshold to start alert",
	"styleWidth":     "style progress bar width",
	"styleBar":       "style progress bar appearance",
	"styleIcon":      "style icon appearance",
	"task":           "set the string for current task",
	"toggleProgress": "turn progress bar on/off",
	"toggleBell":     "turn terminal bell on/off",
	"toggleClock":    "display time on/off",
	"toggleSymbol":   "display symbol on/off",
	"toggleNotify":   "turn timer notifications on/off for notify-send",
	"togglePercent":  "display interval percentage on/off",
	"toggleRestart":  "turn automatic timer restart on/off",
	"toggleReverse":  "timer displays time descending/ascending",
	"toggleTmux":     "turn timer notifications on/off for tmux",
}

var Shorthand = map[string]string{
	"bar":      "b",
	"progress": "p",
	"bell":     "l",
	"timer":    "t",
	"icon":     "i",
	"symbol":   "s",
	"percent":  "P",
	"notify":   "n",
	"tmux":     "t",
	"restart":  "r",
	"reverse":  "v",
	"clock":    "c",
	"alert":    "a",
	"break":    "k",
	"width":    "w",
	"help":     "h",
}

func PrintBasicUsage() {
	fmt.Printf("Usage of %v:\n", programName)
	fmt.Printf("  start\n\t%v\n", UsageString["start"])
	fmt.Printf("  stop\n\t%v\n", UsageString["stop"])
	fmt.Printf("  pause\n\t%v\n", UsageString["pause"])
	fmt.Printf("  resume\n\t%v\n", UsageString["resume"])
	fmt.Printf("  break\n\t%v\n", UsageString["break"])
	fmt.Printf("  run\n\t%v\n", UsageString["run"])
	fmt.Printf("  task\n\t%v\n", UsageString["task"])
	fmt.Printf("  clear\n\t%v\n", UsageString["clear"])
	fmt.Printf("  status\n\t%v\n", UsageString["status"])
	fmt.Printf("  info\n\t%v\n", UsageString["info"])
	fmt.Printf("  clean\n\t%v\n", UsageString["clean"])
	fmt.Printf("  help\n\t%v\n", UsageString["help"])
	fmt.Printf("\n")
}
