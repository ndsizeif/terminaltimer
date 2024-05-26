package main

import (
	"fmt"
	Tmux "terminalTimer/tmuxmenu"
	"time"
)

// create timer object, load state, config, and adjust timer values based on config
func InitializeTimer() (Task, error) {
	var t Task
	var err error

	t.State, err = InitializeState()
	if err != nil {
		// this should never != nil, but we log if something went wrong
		t.State.Debug.Print(t.State.Debug.Trace(), err)
	}
	t.Config, err = InitializeConfig()
	if err != nil {
		t.State.Debug.Print(t.State.Debug.Trace(), err)
	}
	t.LoadSymbols()
	t.LoadInputMaps()
	t.Tmux = Tmux.Initialize()
	if t.Config.Log { // defaults to false, must be explicitly set to true
		t.Log, err = InitializeTimeLog()
		if err != nil { // disable logging on any error
			t.State.Debug.Print(t.State.Debug.Trace(), err)
			t.Log.Enable(false)
			return t, nil
		}
		t.Log.Enable(t.Config.Log) // enable log if config file has it enabled
	}
	return t, nil
}

type Task struct {
	State  *State
	Config *Config
	Tmux   *Tmux.Menu
	Log    *History

	Symbols  map[string]string              // icon symbols
	Progress map[string]string              // progress bar characters
	Command  map[string]func() error        // timer command map
	Duration map[string]func(time.Duration) // set duration map
	Toggle   map[string]func(state bool)    // set boolean settings map
	Option   map[string]func(int)           // set integer settings map
	Status   map[string]func() string       // timer status map
	Display  map[string]func(string)        // display task string
}

func (t *Task) LoadSymbols() {
	switch {
	case t.Config.Icon == 1:
		t.Symbols = icon_solid
	case t.Config.Icon == 2:
		t.Symbols = icon_trace
	case t.Config.Icon == 3:
		t.Symbols = icon_ascii
	default:
		t.Symbols = icon_solid
	}
	switch {
	case t.Config.BarStyle == 1:
		t.Progress = bar_solid
	case t.Config.BarStyle == 2:
		t.Progress = bar_solid_rev
	case t.Config.BarStyle == 3:
		t.Progress = bar_shade
	case t.Config.BarStyle == 4:
		t.Progress = bar_shade_rev
	case t.Config.BarStyle == 5:
		t.Progress = bar_ascii
	default:
		t.Progress = bar_solid
	}
}

// handle notification and restart events
func (t *Task) GetTime() {
	if t.State.TimerHasExpired() {
		if t.Config.Restart {
			t.Start() // restart
		}
	}
	t.NotificationUpdate() // check notification status
}

func (t *Task) Start() error {
	t.State.SetStart(time.Now())
	t.State.SetPause(time.Time{})
	err := t.State.Save()
	if err != nil {
		t.State.Debug.Print("Start()", err)
		return err
	}
	t.Message("start")
	return nil
}

func (t *Task) Stop() error {
	t.State.SetStart(time.Time{})
	t.State.SetPause(time.Time{})
	err := t.State.Save()
	if err != nil {
		t.State.Debug.Print("Stop()", err)
		return err
	}
	t.Message("stop")
	return nil
}

func (t *Task) Pause() error {
	if t.State.TimerIsStopped() || t.State.TimerHasExpired() {
		return nil
	}
	if t.State.TimerIsPaused() {
		t.State.Debug.Print("Pause -> resume")
		t.Resume()
		return nil
	}
	t.State.SetPause(time.Now())
	t.State.Debug.Print("Pause -> pause")
	err := t.State.Save()
	if err != nil {
		t.State.Debug.Print("Pause() t.State.Save() error", err)
		return err
	}
	t.Message("pause")
	return nil
}

func (t *Task) Resume() error {
	if !t.State.TimerIsPaused() {
		t.State.Debug.Print("Resume(): timer not paused")
		return nil
	}
	pauseDuration := time.Since(t.State.TimePause)
	adjustedStart := t.State.TimeStart.Add(pauseDuration)
	t.State.TimeStart = adjustedStart
	t.State.TimePause = time.Time{}
	err := t.State.Save()
	if err != nil {
		t.State.Debug.Print("Resume() t.State.Save() error", err)
		return err
	}
	t.Message("resume")
	return nil
}

func (t *Task) Break() error {
	oldtime := t.State.TimeStart
	t.State.SetPause(time.Time{})
	t.State.SetStart(time.Now().Add(-t.State.TimeInterval)) // set start time to (time now - interval)
	dur := oldtime.Sub(t.State.TimeStart)
	t.State.Debug.Print("BREAK remaining time: ", dur)
	err := t.State.Save()
	if err != nil {
		t.State.Debug.Print("Break() t.State.Save() error", err)
		return err
	}
	t.Message("break")
	return nil
}

func (t *Task) Clear() error {
	t.State.ClearTask() // returned error will always be nil
	err := t.State.Save()
	if err != nil {
		t.State.Debug.Print("Clear() t.State.Save() error", err)
		return err
	}
	return nil
}

func (t *Task) Info() string {
	return fmt.Sprintf("%v %v\n", t.State.TimeInterval, t.State.TimeBreak)
}

func (t *Task) GetState() string {
	var state, notify, restart, task string
	switch {
	case t.State.TimerHasExpired():
		state = stateString["expired"]
	case t.State.TimerIsStopped():
		state = stateString["stopped"]
	case t.State.TimerOnBreak():
		if t.State.TimerIsPaused() {
			state = stateString["breakp"]
		} else {
			state = stateString["break"]
		}
	case !t.State.TimerIsStopped():
		if t.State.TimerIsPaused() {
			state = stateString["paused"]
		} else {
			state = stateString["on"]
		}
	}
	if t.Config.Restart {
		restart = fmt.Sprintf("%v ",t.Symbols["restart"])
	}
	if t.Config.Notify {
		notify = fmt.Sprintf("%v ",t.Symbols["notify"])
	}
	if len(t.State.Task) > 0 {
		task = t.State.Task + " "
		state = t.State.Task + " " + state
	}
	return fmt.Sprintf("%v%v %v%v", task, state, restart, notify)
}

func (t *Task) NotificationUpdate() {
	const notifyThreshold = 1000 * time.Millisecond
	if t.State.TimerIsPaused() { // don't send notifications if paused
		return
	}
	var message string
	switch {
	case t.State.TimerOnBreak():
		if (t.State.GetTotal()-t.State.GetElapsed()) < time.Duration(notifyThreshold) && (t.State.GetTotal()-t.State.GetElapsed()) > time.Duration(0*time.Second) {
			message = messageDone
			if t.Config.Restart {
				message = messageWork
			}
			t.SendNotification(message)
			t.SendTmuxNotification(message)
			t.Message("break over")
			return
		}
	case !t.State.TimerOnBreak():
		if t.State.GetIntervalRemaining() < time.Duration(notifyThreshold) && t.State.GetIntervalRemaining() > time.Duration(0*time.Second) {
			message = messageBreak
			t.SendNotification(message)
			t.SendTmuxNotification(message)
			t.Message("completed")
			return
		}
	}
}

func (t *Task) SendNotification(message string) {
	if !t.Config.Notify {
		return
	}
	if !t.State.TimerOnBreak() {
		message = messageBreak
	} else {
		message = messageDone
		if t.Config.Restart {
			message = messageWork
		}
	}
	err := NotifySend(message)
	if err != nil {
		message = err.Error() // show err in debug log
	}
	t.State.Debug.Print("NOTIFY-SEND:", message)
}

func (t *Task) SendTmuxNotification(message string) {
	var title, symbol string
	if !t.Config.NotifyTmux {
		return
	}
	if !t.State.TimerOnBreak() {
		message = messageBreak
		symbol = t.Symbols["break"]
	} else {
		message = messageDone
		symbol = t.Symbols["on"]
		if t.Config.Restart { // restarting timer displays a different message
			message = messageWork
		}
	}
	title = fmt.Sprintf(" %v %v %v ", symbol, "Notification", symbol)
	err := t.Tmux.Close() // close any existing tmux popup before spawning a new one
	if err != nil {
		t.State.Debug.Print("TMUX CLOSE POPUP:", err)
	}
	err = t.Tmux.Open(title, message)
	if err != nil {
		t.State.Debug.Print("TMUX OPEN POPUP:", err)
	}
	t.State.Debug.Print("TMUX NOTIFICATION:", message)
}

func (t *Task) FormatTime(remaining time.Duration) (result string) {
	days := int(remaining.Hours() / 24)
	hours := int(remaining.Hours()) % 24
	minutes := int(remaining.Minutes()) % 60
	seconds := int(remaining.Seconds()) % 60
	switch {
	case remaining > time.Duration(24*time.Hour):
		if t.Config.HideSeconds {
			result = fmt.Sprintf("%d:%02d:%02d", days, hours, minutes)
			return
		}
		result = fmt.Sprintf("%d:%02d:%02d:%02d", days, hours, minutes, seconds)
	case remaining > time.Duration(60*time.Minute):
		if t.Config.HideSeconds {
			result = fmt.Sprintf("%d:%02d", hours, minutes)
			return
		}
		result = fmt.Sprintf("%d:%02d:%02d", hours, minutes, seconds)
	case remaining > time.Duration(60*time.Second):
		if t.Config.HideSeconds {
			result = fmt.Sprintf("%02d", minutes)
			return
		}
		result = fmt.Sprintf("%02d:%02d", minutes, seconds)
	case remaining < time.Duration(0*time.Second): // Remove case after testing, it should never occur
		result = fmt.Sprintf("-%02d:%02d", minutes, seconds)
	default:
		if t.Config.HideSeconds {
			result = fmt.Sprintf("%02d", minutes)
			return
		}
		result = fmt.Sprintf("%02d:%02d", minutes, seconds)
	}
	return result
}

func (t *Task) Message(message string) {
	if !t.Config.Log { // do not attempt print if log == false, there will be no valid object to print
		return
	}
	format := fmt.Sprintf("%-10v\t%-20v\t%v", TimeStamp, message, t.State.Task)
	t.Log.Print(format)
}
