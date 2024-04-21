package main

import (
	"time"
)

// return command based on string input, check if passed string is in map, return matching function if true
func (t *Task) ExecuteCommand(s string) func() error {
	var cmd func() error
	cmd, ok := t.Command[s]
	if ok {
		t.State.Debug.Print(t.State.Debug.Trace(), "ExecuteCommand:", s)
		return cmd
	}
	return nil
}

func (t *Task) ToggleOption(s string) func(bool) {
	opt, ok := t.Toggle[s]
	if ok {
		t.State.Debug.Print(t.State.Debug.Trace(), "ToggleOption:", s)
		return opt
	}
	return nil
}

func (t *Task) SetOption(s string) func(int) {
	opt, ok := t.Option[s]
	if ok {
		t.State.Debug.Print(t.State.Debug.Trace(), "SetOption:", s)
		return opt
	}
	return nil
}

func (t *Task) SetDuration(s string) func(time.Duration) {
	dur, ok := t.Duration[s]
	if ok {
		t.State.Debug.Print(t.State.Debug.Trace(), "SetDuration:", s)
		return dur
	}
	return nil
}

func (t *Task) SetString(s string) func(string) {
	cmd, ok := t.Display[s]
	if ok {
		t.State.Debug.Print(t.State.Debug.Trace(), "SetString:", s)
		return cmd
	}
	return nil
}

// map input string to timer function
func (t *Task) LoadInputMaps() error {
	t.Command = map[string]func() error{
		"start":  t.Start,
		"stop":   t.Stop,
		"pause":  t.Pause,
		"resume": t.Resume,
		"break":  t.Break,
		"run":    t.RunInline,
		"clear":  t.Clear,
	}
	t.Duration = map[string]func(time.Duration){
		"timer": t.State.SetInterval,
		"break": t.State.SetBreak,
		"alert": t.State.SetAlert,
	}
	t.Toggle = map[string]func(bool){
		"restart": t.Config.SetRestart,
		"reverse": t.Config.SetReverse,
		"bell":    t.Config.SetBell,
		"clock":   t.Config.SetHideTime,
		"seconds": t.Config.SetHideSeconds,
		"bar":     t.Config.SetHideBar,
		"icon":    t.Config.SetHideIcon,
		"percent": t.Config.SetPercent,
		"notify":  t.Config.SetNotify,
		"tmux":    t.Config.SetNotifyTmux,
	}
	t.Option = map[string]func(int){
		"size":   t.Config.SetBarSize,
		"style":  t.Config.SetBarStyle,
		"symbol": t.Config.SetIcon,
	}
	t.Status = map[string]func() string{
		"info":   t.Info,
		"status": t.GetState,
	}
	t.Display = map[string]func(string){
		"task": t.State.SetTask,
	}
	return nil
}
